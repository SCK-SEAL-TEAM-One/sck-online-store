package payment

import (
	"context"
	"database/sql"
	"log/slog"
	"store-service/internal/metrics"
	"store-service/internal/order"
	"store-service/internal/product"
	"store-service/internal/shipping"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type PaymentInterface interface {
	ConfirmPayment(ctx context.Context, uid int, submitedPayment SubmitedPayment) (SubmitedPaymentResponse, error)
}

type PaymentService struct {
	BankGateway       BankGatewayInterface
	ShippingGateway   ShippingGatewayInterface
	OrderRepository   order.OrderRepository
	ProductRepository product.ProductRepository
}

type BankGatewayInterface interface {
	Payment(ctx context.Context, paymentDetail PaymentDetail) (string, error)
	GetCardDetail(ctx context.Context, orgID int, userID int) (CardDetail, error)
}

type ShippingGatewayInterface interface {
	GetTrackingNumber(ctx context.Context, shippingGatewaySubmit shipping.ShippingGatewaySubmit) (string, error)
}

func (service PaymentService) ConfirmPayment(ctx context.Context, uid int, submitedPayment SubmitedPayment) (SubmitedPaymentResponse, error) {
	orgID := 1
	orderNumber := submitedPayment.OrderNumber
	currency := "USD"
	now := time.Now()

	orderDetail, err := service.OrderRepository.GetOrderByOrderNumber(ctx, orderNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.ErrorContext(ctx, "Order not found for payment",
				"log_type", "error", "error_code", "ORDER_NOT_FOUND", "error_message", err.Error(),
				"user_id", uid, "order_number", orderNumber)
			return SubmitedPaymentResponse{}, order.ErrOrderNotFound
		}
		slog.ErrorContext(ctx, "OrderRepository.GetOrderByOrderNumber failed",
			"log_type", "error", "error_code", "ORDER_QUERY_FAILED", "error_message", err.Error(),
			"user_id", uid, "order_number", orderNumber)
		return SubmitedPaymentResponse{}, err
	}

	cardDetail, err := service.BankGateway.GetCardDetail(ctx, orgID, uid)
	if err != nil {
		slog.ErrorContext(ctx, "BankGateway.GetCardDetail failed",
			"log_type", "error", "error_code", "CARD_DETAIL_FAILED", "error_message", err.Error(),
			"user_id", uid, "order_number", orderNumber)
		return SubmitedPaymentResponse{}, err
	}

	paymentdetail := PaymentDetail{
		CardNumber:   cardDetail.CardNumber,
		CVV:          cardDetail.CVV,
		ExpiredMonth: cardDetail.ExpiredMonth,
		ExpiredYear:  cardDetail.ExpiredYear,
		CardName:     cardDetail.CardName,
		Amount:       orderDetail.TotalPrice,
		Currency:     currency,
		MerchantID:   orgID,
	}
	paymentStart := time.Now()
	transactionId, err := service.BankGateway.Payment(ctx, paymentdetail)
	if metrics.PaymentDuration != nil {
		metrics.PaymentDuration.Record(ctx, time.Since(paymentStart).Seconds())
	}
	if err != nil {
		slog.ErrorContext(ctx, "BankGateway.Payment failed",
			"log_type", "error", "error_code", "BANK_PAYMENT_FAILED", "error_message", err.Error(),
			"user_id", uid, "order_number", orderNumber, "amount", orderDetail.TotalPrice)
		if metrics.PaymentAttempts != nil {
			metrics.PaymentAttempts.Add(ctx, 1,
				metric.WithAttributes(
					attribute.String("status", "failed"),
					attribute.String("error_type", "bank_gateway_error"),
				),
			)
		}
		return SubmitedPaymentResponse{}, err
	}

	orderProductList, err := service.OrderRepository.GetOrderProduct(ctx, orderDetail.ID)
	if err != nil {
		slog.ErrorContext(ctx, "OrderRepository.GetOrderProduct failed",
			"log_type", "error", "error_code", "ORDER_PRODUCT_QUERY_FAILED", "error_message", err.Error(),
			"user_id", uid, "order_number", orderNumber)
		return SubmitedPaymentResponse{}, err
	}
	for _, orderProduct := range orderProductList {
		err = service.ProductRepository.UpdateStock(ctx, orderProduct.ProductID, orderProduct.Quantity)
		if err != nil {
			slog.ErrorContext(ctx, "ProductRepository.UpdateStock failed",
				"log_type", "error", "error_code", "STOCK_UPDATE_FAILED", "error_message", err.Error(),
				"user_id", uid, "product_id", orderProduct.ProductID)
			return SubmitedPaymentResponse{}, err
		}
		slog.InfoContext(ctx, "Stock updated",
			"log_type", "state_change",
			"entity_type", "product_stock",
			"entity_id", orderProduct.ProductID,
			"changed_by", uid,
			slog.Any("after", map[string]any{"quantity_deducted": orderProduct.Quantity}),
			slog.Any("changed_fields", []string{"stock"}),
		)
	}

	err = service.OrderRepository.UpdateOrderTransaction(ctx, orderDetail.ID, transactionId)
	if err != nil {
		slog.ErrorContext(ctx, "OrderRepository.UpdateOrderTransaction failed",
			"log_type", "error", "error_code", "ORDER_TRANSACTION_FAILED", "error_message", err.Error(),
			"user_id", uid, "order_number", orderNumber)
		return SubmitedPaymentResponse{}, err
	}

	slog.InfoContext(ctx, "Order payment recorded",
		"log_type", "state_change",
		"entity_type", "order",
		"entity_id", orderNumber,
		"changed_by", uid,
		slog.Any("after", map[string]any{"transaction_id": transactionId, "status": "paid"}),
		slog.Any("changed_fields", []string{"transaction_id", "status"}),
	)

	shippingGatewaySubmit := shipping.ShippingGatewaySubmit{
		ShippingMethodID: orderDetail.ShippingMethodID,
	}
	trackingNumber, err := service.ShippingGateway.GetTrackingNumber(ctx, shippingGatewaySubmit)
	if err != nil {
		slog.ErrorContext(ctx, "ShippingGateway.GetTrackingNumber failed",
			"log_type", "error", "error_code", "TRACKING_NUMBER_FAILED", "error_message", err.Error(),
			"user_id", uid, "order_number", orderNumber)
		return SubmitedPaymentResponse{}, err
	}

	err = service.OrderRepository.UpdateOrderTrackingNumber(ctx, orderDetail.ID, trackingNumber)
	if err != nil {
		slog.ErrorContext(ctx, "OrderRepository.UpdateOrderTrackingNumber failed",
			"log_type", "error", "error_code", "TRACKING_UPDATE_FAILED", "error_message", err.Error(),
			"user_id", uid, "order_number", orderNumber)
		return SubmitedPaymentResponse{}, err
	}

	slog.InfoContext(ctx, "Tracking number assigned",
		"log_type", "state_change",
		"entity_type", "order",
		"entity_id", orderNumber,
		"changed_by", uid,
		slog.Any("after", map[string]any{"tracking_number": trackingNumber}),
		slog.Any("changed_fields", []string{"tracking_number"}),
	)

	if metrics.PaymentAttempts != nil {
		metrics.PaymentAttempts.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("status", "success"),
				attribute.String("error_type", ""),
			),
		)
	}

	return SubmitedPaymentResponse{
		OrderNumber:      orderDetail.OrderNumber,
		PaymentDate:      now,
		ShippingMethodID: orderDetail.ShippingMethodID,
		TrackingNumber:   trackingNumber,
	}, nil
}
