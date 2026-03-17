package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"store-service/internal/auth"
	"store-service/internal/cart"
	"store-service/internal/common"
	"store-service/internal/metrics"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"store-service/internal/point"
	"store-service/internal/product"
	"store-service/internal/shipping"
	"time"
)

type OrderInterface interface {
	CreateOrder(ctx context.Context, uid int, submitedOrder SubmitedOrder) (Order, error)
	OrderBurnPoint(ctx context.Context, uid int, burn int) (point.TotalPoint, error)
	GetOrderSummary(ctx context.Context, orderNumber string) (OrderSummary, error)
	GeneratePDFFromData(orderDetail OrderSummary) ([]byte, error)
}

type OrderService struct {
	CartRepository     cart.CartRepository
	OrderRepository    OrderRepository
	PointService       point.PointInterface
	ProductRepository  product.ProductRepository
	ShippingRepository shipping.ShippingRepository
	UserRepository     auth.UserRepository
	PDFHelper          PDFHelper
	OrderHelper        OrderHelperInterface
	Clock              func() time.Time
}

type CartRepository interface {
	DeleteCart(userID int, productID int)
}
type PointService interface {
	DeductPoint(uid int, submitedPoint point.SubmitedPoint) (point.TotalPoint, error)
}

type ProductRepository interface {
	GetProductByID(id int) (product.ProductDetail, error)
}

type ShippingRepository interface {
	GetShippingMethodByID(id int) (shipping.ShippingMethodDetail, error)
}

var PaymentMethod = map[int]string{
	1: "Credit Card / Debit Card",
	2: "Line Pay",
}

var ShippingMethod = map[int]string{
	1: "Kerry",
	2: "Thai Post",
	3: "Lineman",
}

var ErrOrderNotFound = errors.New("Order not found")

func (orderService OrderService) CreateOrder(ctx context.Context, uid int, submitedOrder SubmitedOrder) (Order, error) {
	_, err := orderService.PointService.CheckBurnPoint(ctx, uid, -(submitedOrder.BurnPoint))
	if err != nil {
		return Order{}, err
	}

	if len(submitedOrder.Cart) == 0 {
		return Order{}, fmt.Errorf("There is no product in order, please try again")
	}

	subtotalPrice := 0.0
	for _, productSelected := range submitedOrder.Cart {
		product, _ := orderService.ProductRepository.GetProductByID(ctx, productSelected.ProductID)
		subtotalPrice = subtotalPrice + (product.Price * float64(productSelected.Quantity))
	}

	subtotalPriceTHB := common.ConvertToThb(subtotalPrice).LongDecimal
	discountPriceTHB := common.ConvertToThb(submitedOrder.DiscountPrice).LongDecimal
	totalPriceTHB := subtotalPriceTHB - discountPriceTHB

	shippingDetail, _ := orderService.ShippingRepository.GetShippingMethodByID(ctx, submitedOrder.ShippingMethodID)
	shippingFeeTHB := shippingDetail.Fee

	now := orderService.Clock()
	yearPrefix := now.Format("06") // Format: YY

	const maxRetries = 3
	var orderID int
	var orderNumber string
	var orderDetail OrderDetail

	for attempt := 0; attempt < maxRetries; attempt++ {
		seq := 1 // Default SEQ number for beginning new year
		lastOrderNumber, err := orderService.OrderRepository.GetLastOrderNumber(ctx, yearPrefix)
		if err != nil && err != sql.ErrNoRows {
			slog.ErrorContext(ctx, "OrderRepository.GetLastOrderNumber failed",
				"log_type", "error", "error_code", "ORDER_SEQ_FAILED", "error_message", err.Error(), "user_id", uid)
			return Order{}, err
		}

		if err == nil {
			seq, err = orderService.OrderHelper.GetNextSequence(lastOrderNumber)
			if err != nil {
				slog.ErrorContext(ctx, "OrderHelper.GetNextSequence failed",
					"log_type", "error", "error_code", "ORDER_SEQ_FAILED", "error_message", err.Error(), "user_id", uid)
				return Order{}, err
			}
		}

		orderNumber, err = orderService.OrderHelper.GenerateOrderNumber(submitedOrder.PaymentMethodID, submitedOrder.ShippingMethodID, seq, now)
		if err != nil {
			slog.ErrorContext(ctx, "OrderHelper.GenerateOrderNumber failed",
				"log_type", "error", "error_code", "ORDER_NUMBER_FAILED", "error_message", err.Error(), "user_id", uid)
			return Order{}, err
		}
		orderDetail = OrderDetail{
			OrderNumber:      orderNumber,
			ShippingMethodID: submitedOrder.ShippingMethodID,
			PaymentMethodID:  submitedOrder.PaymentMethodID,
			SubTotalPrice:    subtotalPriceTHB,
			DiscountPrice:    discountPriceTHB,
			TotalPrice:       totalPriceTHB + shippingFeeTHB,
			ShippingFee:      shippingFeeTHB,
			BurnPoint:        submitedOrder.BurnPoint,
			EarnPoint:        common.CalculatePoint(totalPriceTHB),
		}

		orderID, err = orderService.OrderRepository.CreateOrder(ctx, uid, orderDetail)
		if err != nil {
			if isDuplicateKeyError(err) && attempt < maxRetries-1 {
				slog.WarnContext(ctx, "Duplicate order number, retrying",
					"log_type", "state_change", "order_number", orderNumber, "attempt", attempt+1, "user_id", uid)
				continue
			}
			slog.ErrorContext(ctx, "OrderRepository.CreateOrder failed",
				"log_type", "error", "error_code", "ORDER_INSERT_FAILED", "error_message", err.Error(), "user_id", uid)
			return Order{}, err
		}
		break
	}

	shippingInfo := ShippingInfo{
		ShippingMethodID:     submitedOrder.ShippingMethodID,
		ShippingAddress:      submitedOrder.ShippingAddress,
		ShippingSubDistrict:  submitedOrder.ShippingSubDistrict,
		ShippingDistrict:     submitedOrder.ShippingDistrict,
		ShippingProvince:     submitedOrder.ShippingProvince,
		ShippingZipCode:      submitedOrder.ShippingZipCode,
		RecipientFirstName:   submitedOrder.RecipientFirstName,
		RecipientLastName:    submitedOrder.RecipientLastName,
		RecipientPhoneNumber: submitedOrder.RecipientPhoneNumber,
	}
	_, err = orderService.OrderRepository.CreateShipping(ctx, uid, orderID, shippingInfo)
	if err != nil {
		slog.ErrorContext(ctx, "OrderRepository.CreateShipping failed",
			"log_type", "error", "error_code", "SHIPPING_INSERT_FAILED", "error_message", err.Error(), "user_id", uid)
		return Order{}, err
	}

	for _, selectedProduct := range submitedOrder.Cart {
		product, err := orderService.ProductRepository.GetProductByID(ctx, selectedProduct.ProductID)
		err = orderService.OrderRepository.CreateOrderProduct(ctx, orderID, selectedProduct.ProductID, selectedProduct.Quantity, product.Price)
		if err != nil {
			slog.ErrorContext(ctx, "OrderRepository.CreateOrderProduct failed",
				"log_type", "error", "error_code", "ORDER_PRODUCT_FAILED", "error_message", err.Error(), "user_id", uid,
				"product_id", selectedProduct.ProductID)
			return Order{}, err
		}

		orderService.CartRepository.DeleteCart(ctx, uid, selectedProduct.ProductID)
	}

	if submitedOrder.BurnPoint > 0 {
		orderService.OrderBurnPoint(ctx, uid, submitedOrder.BurnPoint)
	}

	if metrics.OrdersCreated != nil {
		metrics.OrdersCreated.Add(ctx, 1,
			metric.WithAttributes(
				attribute.String("status", "success"),
				attribute.String("payment_method", PaymentMethod[submitedOrder.PaymentMethodID]),
				attribute.String("shipping_method", ShippingMethod[submitedOrder.ShippingMethodID]),
			),
		)
		metrics.OrderRevenue.Add(ctx, orderDetail.TotalPrice,
			metric.WithAttributes(
				attribute.String("payment_method", PaymentMethod[submitedOrder.PaymentMethodID]),
			),
		)
		metrics.OrderItemsCount.Record(ctx, int64(len(submitedOrder.Cart)))
	}

	return Order{
		OrderNumber: orderNumber,
	}, nil
}

func (orderService OrderService) OrderBurnPoint(ctx context.Context, uid int, burn int) (point.TotalPoint, error) {
	submit := point.SubmitedPoint{
		Amount: -(burn),
	}

	totalPoint, err := orderService.PointService.DeductPoint(ctx, uid, submit)
	if err != nil {
		return point.TotalPoint{}, err
	}
	return totalPoint, nil
}

func (orderService OrderService) GetOrderSummary(ctx context.Context, orderNumber string) (OrderSummary, error) {
	orderDetail, err := orderService.OrderRepository.GetOrderWithTrackingNumberByOrderNumber(ctx, orderNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.ErrorContext(ctx, "Order not found",
				"log_type", "error", "error_code", "ORDER_NOT_FOUND", "error_message", err.Error(),
				"user_id", 0, "order_number", orderNumber)
			return OrderSummary{}, ErrOrderNotFound
		}
		slog.ErrorContext(ctx, "OrderRepository.GetOrderWithTrackingNumberByOrderNumber failed",
			"log_type", "error", "error_code", "ORDER_QUERY_FAILED", "error_message", err.Error(),
			"user_id", 0, "order_number", orderNumber)
		return OrderSummary{}, err
	}

	orderedProducts, err := orderService.OrderRepository.GetOrderProductWithPrice(ctx, orderDetail.ID)
	if err != nil {
		slog.ErrorContext(ctx, "OrderRepository.GetOrderProductWithPrice failed",
			"log_type", "error", "error_code", "ORDER_PRODUCT_QUERY_FAILED", "error_message", err.Error(),
			"user_id", 0, "order_number", orderNumber)
		return OrderSummary{}, err
	}

	var productList []OrderSummaryProduct
	for _, orderProduct := range orderedProducts {
		totalPrice := orderProduct.Price * float64(orderProduct.Quantity)

		totalPriceTHB := common.ConvertToThb(totalPrice)
		priceTHB := common.ConvertToThb(orderProduct.Price)
		product := OrderSummaryProduct{
			ProductBrand:  orderProduct.ProductBrand,
			ProductName:   orderProduct.ProductName,
			Quantity:      orderProduct.Quantity,
			PriceTHB:      priceTHB.ShortDecimal,
			TotalPriceTHB: totalPriceTHB.ShortDecimal,
		}
		productList = append(productList, product)
	}

	paymentMethod := PaymentMethod[orderDetail.PaymentMethodID]
	shippingMethod := ShippingMethod[orderDetail.ShippingMethodID]

	userDetail, err := orderService.UserRepository.FindByID(ctx, orderDetail.UserID)
	if err != nil {
		slog.ErrorContext(ctx, "UserRepository.FindByID failed",
			"log_type", "error", "error_code", "USER_QUERY_FAILED", "error_message", err.Error(),
			"user_id", orderDetail.UserID)
		return OrderSummary{}, err
	}

	factor2 := math.Pow(10, 2)
	subTotal := math.Round(orderDetail.SubTotalPrice*factor2) / factor2
	totalPrice := math.Round(orderDetail.TotalPrice*factor2) / factor2
	shippingFee := math.Round(orderDetail.ShippingFee*factor2) / factor2

	bangkok, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		slog.ErrorContext(ctx, "Could not load timezone",
			"log_type", "error", "error_code", "TIMEZONE_LOAD_FAILED", "error_message", err.Error(), "user_id", 0)
		return OrderSummary{}, err
	}

	issuedDate := orderDetail.Updated.In(bangkok).Format("02-01-2006 15:04:05")

	orderSummary := OrderSummary{
		OrderNumber:      orderDetail.OrderNumber,
		FirstName:        userDetail.FirstName,
		LastName:         userDetail.LastName,
		TrackingNumber:   orderDetail.TrackingNumber,
		ShippingMethod:   shippingMethod,
		PaymentMethod:    paymentMethod,
		OrderProductList: productList,
		SubTotalPrice:    subTotal,
		TotalPrice:       totalPrice,
		ShippingFee:      shippingFee,
		ReceivingPoint:   orderDetail.EarnPoint,
		IssuedDate:       issuedDate,
	}

	return orderSummary, nil
}

func (orderService OrderService) GeneratePDFFromData(orderSummary OrderSummary) ([]byte, error) {
	pdfBytes, err := orderService.PDFHelper.GenerateOrderSummaryPDF(orderSummary)
	if err != nil {
		slog.Error("PDFHelper.GenerateOrderSummaryPDF failed",
			"log_type", "error", "error_code", "PDF_GENERATION_FAILED", "error_message", err.Error(), "user_id", 0)
		return []byte(""), err
	}

	return pdfBytes, nil
}

func isDuplicateKeyError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Duplicate entry")
}
