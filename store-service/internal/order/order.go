package order

import (
	"fmt"
	"log"
	"math"
	"store-service/internal/auth"
	"store-service/internal/cart"
	"store-service/internal/common"
	"store-service/internal/point"
	"store-service/internal/product"
	"store-service/internal/shipping"
)

type OrderInterface interface {
	CreateOrder(uid int, submitedOrder SubmitedOrder) (Order, error)
	OrderBurnPoint(uid int, burn int) (point.TotalPoint, error)
	GetOrderSummary(orderID int) (OrderSummary, error)
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

func (orderService OrderService) CreateOrder(uid int, submitedOrder SubmitedOrder) (Order, error) {
	_, err := orderService.PointService.CheckBurnPoint(uid, -(submitedOrder.BurnPoint))
	if err != nil {
		return Order{}, err
	}

	if len(submitedOrder.Cart) == 0 {
		return Order{}, fmt.Errorf("There is no product in order, please try again")
	}

	subtotalPrice := 0.0
	for _, productSelected := range submitedOrder.Cart {
		product, _ := orderService.ProductRepository.GetProductByID(productSelected.ProductID)
		subtotalPrice = subtotalPrice + (product.Price * float64(productSelected.Quantity))
	}

	subtotalPriceTHB := common.ConvertToThb(subtotalPrice).LongDecimal
	discountPriceTHB := common.ConvertToThb(submitedOrder.DiscountPrice).LongDecimal
	totalPriceTHB := subtotalPriceTHB - discountPriceTHB

	shippingDetail, _ := orderService.ShippingRepository.GetShippingMethodByID(submitedOrder.ShippingMethodID)
	shippingFeeTHB := shippingDetail.Fee

	orderDetail := OrderDetail{
		ShippingMethodID: submitedOrder.ShippingMethodID,
		PaymentMethodID:  submitedOrder.PaymentMethodID,
		SubTotalPrice:    subtotalPriceTHB,
		DiscountPrice:    discountPriceTHB,
		TotalPrice:       totalPriceTHB + shippingFeeTHB,
		ShippingFee:      shippingFeeTHB,
		BurnPoint:        submitedOrder.BurnPoint,
		EarnPoint:        common.CalculatePoint(totalPriceTHB),
	}
	orderID, err := orderService.OrderRepository.CreateOrder(uid, orderDetail)
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
	_, err = orderService.OrderRepository.CreateShipping(uid, orderID, shippingInfo)
	if err != nil {
		log.Printf("OrderRepository.CreateShipping internal error %s", err.Error())
		return Order{}, err
	}

	for _, selectedProduct := range submitedOrder.Cart {
		product, err := orderService.ProductRepository.GetProductByID(selectedProduct.ProductID)
		err = orderService.OrderRepository.CreateOrderProduct(orderID, selectedProduct.ProductID, selectedProduct.Quantity, product.Price)
		if err != nil {
			log.Printf("OrderRepository.CreateOrderProduct internal error %s", err.Error())
			return Order{}, err
		}

		orderService.CartRepository.DeleteCart(uid, selectedProduct.ProductID)
	}

	if submitedOrder.BurnPoint > 0 {
		orderService.OrderBurnPoint(uid, submitedOrder.BurnPoint)
	}

	return Order{
		OrderID: orderID,
	}, nil
}

func (orderService OrderService) OrderBurnPoint(uid int, burn int) (point.TotalPoint, error) {
	submit := point.SubmitedPoint{
		Amount: -(burn),
	}

	totalPoint, err := orderService.PointService.DeductPoint(uid, submit)
	if err != nil {
		return point.TotalPoint{}, err
	}
	return totalPoint, nil
}

func (orderService OrderService) GetOrderSummary(orderID int) (OrderSummary, error) {
	orderDetail, err := orderService.OrderRepository.GetOrderWithTrackingNumberByID(orderID)
	if err != nil {
		log.Printf("OrderRepository.GetOrderByID internal error for orderID %d: %s", orderID, err.Error())
		return OrderSummary{}, err
	}

	orderedProducts, err := orderService.OrderRepository.GetOrderProductWithPrice(orderID)
	if err != nil {
		log.Printf("OrderRepository.GetOrderProduct internal error %s", err.Error())
		return OrderSummary{}, err
	}

	var productList []OrderSummaryProduct
	for _, orderProduct := range orderedProducts {
		priceTHB := common.ConvertToThb(orderProduct.Price)
		product := OrderSummaryProduct{
			ProductBrand: orderProduct.ProductBrand,
			ProductName:  orderProduct.ProductName,
			Quantity:     orderProduct.Quantity,
			PriceTHB:     priceTHB.ShortDecimal,
		}
		productList = append(productList, product)
	}

	paymentMethod := PaymentMethod[orderDetail.PaymentMethodID]
	shippingMethod := ShippingMethod[orderDetail.ShippingMethodID]

	userDetail, err := orderService.UserRepository.FindByID(orderDetail.UserID)
	if err != nil {
		log.Printf("UserRepository.FindByID internal error %s", err.Error())
		return OrderSummary{}, err
	}

	factor2 := math.Pow(10, 2)
	subTotal := math.Round(orderDetail.SubTotalPrice*factor2) / factor2
	totalPrice := math.Round(orderDetail.TotalPrice*factor2) / factor2
	shippingFee := math.Round(orderDetail.ShippingFee*factor2) / factor2

	orderSummary := OrderSummary{
		OrderID:          orderID,
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
