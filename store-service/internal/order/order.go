package order

import (
	"fmt"
	"log"
	"store-service/internal/cart"
	"store-service/internal/point"
	"store-service/internal/product"
	"time"
)

type OrderService struct {
	CartRepository    cart.CartRepository
	OrderRepository   OrderRepository
	PointService      point.PointServiceInterface
	ProductRepository product.ProductRepository
}

type CartRepository interface {
	DeleteCart(userID int, productID int)
}

type OrderInterface interface {
	CreateOrder(uid int, submitedOrder SubmitedOrder) (Order, error)
}

type PointService interface {
	DeductPoint(uid int, submitedPoint point.SubmitedPoint) (point.TotalPoint, error)
}

type ProductRepository interface {
	GetProductByID(id int) product.ProductDetail
}

func (orderService OrderService) CreateOrder(uid int, submitedOrder SubmitedOrder) (Order, error) {
	_, err := orderService.PointService.CheckBurnPoint(uid, -(submitedOrder.BurnPoint))
	if err != nil {
		return Order{}, err
	}

	orderID, err := orderService.OrderRepository.CreateOrder(uid, submitedOrder)
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

func SendNotification(orderID int, trackingNumber string, dateTime time.Time) string {
	return fmt.Sprintf("วันเวลาที่ชำระเงิน %s หมายเลขคำสั่งซื้อ %d คุณสามารถติดตามสินค้าผ่านช่องทาง xx หมายเลข %s", dateTime.Format("2/1/2006 15:04:05"), orderID, trackingNumber)
}
