package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"store-service/internal/auth"
	"store-service/internal/cart"
	"store-service/internal/common"
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

	seq := 1 // Defalt SEQ number for beginning new year
	now := orderService.Clock()
	yearPrefix := now.Format("06") // Format: YY
	lastOrderNumber, err := orderService.OrderRepository.GetLastOrderNumber(ctx, yearPrefix)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("OrderRepository.GetLastOrderNumber internal error %s", err.Error())
		return Order{}, err
	}

	if err == nil {
		seq, err = orderService.OrderHelper.GetNextSequence(lastOrderNumber)
		if err != nil {
			log.Printf("OrderHelper.GetNextSequence internal error %s", err.Error())
			return Order{}, err
		}
	}

	orderNumber, err := orderService.OrderHelper.GenerateOrderNumber(submitedOrder.PaymentMethodID, submitedOrder.ShippingMethodID, seq, now)
	if err != nil {
		log.Printf("OrderHelper.GenerateOrderNumber internal error %s", err.Error())
		return Order{}, err
	}
	orderDetail := OrderDetail{
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

	orderID, err := orderService.OrderRepository.CreateOrder(ctx, uid, orderDetail)
	if err != nil {
		log.Printf("OrderRepository.CreateOrder internal error %s", err.Error())
		return Order{}, err
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
		log.Printf("OrderRepository.CreateShipping internal error %s", err.Error())
		return Order{}, err
	}

	for _, selectedProduct := range submitedOrder.Cart {
		product, err := orderService.ProductRepository.GetProductByID(ctx, selectedProduct.ProductID)
		err = orderService.OrderRepository.CreateOrderProduct(ctx, orderID, selectedProduct.ProductID, selectedProduct.Quantity, product.Price)
		if err != nil {
			log.Printf("OrderRepository.CreateOrderProduct internal error %s", err.Error())
			return Order{}, err
		}

		orderService.CartRepository.DeleteCart(ctx, uid, selectedProduct.ProductID)
	}

	if submitedOrder.BurnPoint > 0 {
		orderService.OrderBurnPoint(ctx, uid, submitedOrder.BurnPoint)
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
			log.Printf("OrderRepository.GetOrderWithTrackingNumberByOrderNumber not found for Order Number %s: %s", orderNumber, err.Error())
			return OrderSummary{}, ErrOrderNotFound
		}
		log.Printf("OrderRepository.GetOrderWithTrackingNumberByOrderNumber internal error for %s", err.Error())
		return OrderSummary{}, err
	}

	orderedProducts, err := orderService.OrderRepository.GetOrderProductWithPrice(ctx, orderDetail.ID)
	if err != nil {
		log.Printf("OrderRepository.GetOrderProduct internal error %s", err.Error())
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
		log.Printf("UserRepository.FindByID internal error %s", err.Error())
		return OrderSummary{}, err
	}

	factor2 := math.Pow(10, 2)
	subTotal := math.Round(orderDetail.SubTotalPrice*factor2) / factor2
	totalPrice := math.Round(orderDetail.TotalPrice*factor2) / factor2
	shippingFee := math.Round(orderDetail.ShippingFee*factor2) / factor2

	bangkok, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Fatal("Could not load timezone:", err)
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
		log.Printf("PDFHelper.GenerateOrderSummaryPDF internal error %s", err.Error())
		return []byte(""), err
	}

	return pdfBytes, nil
}
