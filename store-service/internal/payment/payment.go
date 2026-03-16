package payment

import (
	"context"
	"database/sql"
	"log/slog"
	"store-service/internal/order"
	"store-service/internal/product"
	"store-service/internal/shipping"
	"time"
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
			slog.ErrorContext(ctx, "OrderRepository.GetOrderByOrderNumber not found", "orderNumber", orderNumber, "error", err)
			return SubmitedPaymentResponse{}, order.ErrOrderNotFound
		}
		slog.ErrorContext(ctx, "OrderRepository.GetOrderByOrderNumber internal error", "error", err)
		return SubmitedPaymentResponse{}, err
	}

	cardDetail, err := service.BankGateway.GetCardDetail(ctx, orgID, uid)
	if err != nil {
		slog.ErrorContext(ctx, "BankGateway.GetCardDetail internal error", "error", err)
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
	transactionId, err := service.BankGateway.Payment(ctx, paymentdetail)
	if err != nil {
		slog.ErrorContext(ctx, "BankGateway.Payment internal error", "error", err)
		return SubmitedPaymentResponse{}, err
	}

	orderProductList, err := service.OrderRepository.GetOrderProduct(ctx, orderDetail.ID)
	if err != nil {
		slog.ErrorContext(ctx, "OrderRepository.GetOrderProduct internal error", "error", err)
		return SubmitedPaymentResponse{}, err
	}
	for _, orderProduct := range orderProductList {
		err = service.ProductRepository.UpdateStock(ctx, orderProduct.ProductID, orderProduct.Quantity)
		if err != nil {
			slog.ErrorContext(ctx, "ProductRepository.UpdateStock internal error", "error", err)
			return SubmitedPaymentResponse{}, err
		}
	}

	err = service.OrderRepository.UpdateOrderTransaction(ctx, orderDetail.ID, transactionId)
	if err != nil {
		slog.ErrorContext(ctx, "OrderRepository.UpdateOrderTransaction internal error", "error", err)
		return SubmitedPaymentResponse{}, err
	}

	shippingGatewaySubmit := shipping.ShippingGatewaySubmit{
		ShippingMethodID: orderDetail.ShippingMethodID,
	}
	trackingNumber, err := service.ShippingGateway.GetTrackingNumber(ctx, shippingGatewaySubmit)
	if err != nil {
		slog.ErrorContext(ctx, "ShippingGateway.GetTrackingNumber internal error", "error", err)
		return SubmitedPaymentResponse{}, err
	}

	err = service.OrderRepository.UpdateOrderTrackingNumber(ctx, orderDetail.ID, trackingNumber)
	if err != nil {
		slog.ErrorContext(ctx, "OrderRepository.UpdateOrderTrackingNumber internal error", "error", err)
		return SubmitedPaymentResponse{}, err
	}

	return SubmitedPaymentResponse{
		OrderNumber:      orderDetail.OrderNumber,
		PaymentDate:      now,
		ShippingMethodID: orderDetail.ShippingMethodID,
		TrackingNumber:   trackingNumber,
	}, nil
}
