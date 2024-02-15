package order

import (
	"fmt"
	"log"
	"store-service/internal/cart"
	"store-service/internal/common"
	"store-service/internal/point"
	"store-service/internal/product"
	"store-service/internal/shipping"
)

type OrderInterface interface {
	CreateOrder(uid int, submitedOrder SubmitedOrder) (Order, error)
	OrderBurnPoint(uid int, burn int) (point.TotalPoint, error)
}

type OrderService struct {
	CartRepository     cart.CartRepository
	OrderRepository    OrderRepository
	PointService       point.PointInterface
	ProductRepository  product.ProductRepository
	ShippingRepository shipping.ShippingRepository
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
