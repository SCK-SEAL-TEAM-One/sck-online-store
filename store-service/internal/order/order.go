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
	PointGateway      point.PointGateway
	ProductRepository product.ProductRepository
}

type CartRepository interface {
	DeleteCart(userID int, productID int)
}

type OrderInterface interface {
	CreateOrder(submitedOrder SubmitedOrder) Order
}

type PointGateway interface {
	CreatePoint(uid int, body point.Point) point.Point
}

type ProductRepository interface {
	GetProductByID(id int) product.ProductDetail
}

func (orderService OrderService) CreateOrder(submitedOrder SubmitedOrder) Order {
	uid := 1
	orderID, err := orderService.OrderRepository.CreateOrder(uid, submitedOrder)
	if err != nil {
		log.Printf("OrderRepository.CreateOrder internal error %s", err.Error())
		return Order{}
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
		return Order{}
	}

	for _, selectedProduct := range submitedOrder.Cart {
		product, err := orderService.ProductRepository.GetProductByID(selectedProduct.ProductID)
		err = orderService.OrderRepository.CreateOrderProduct(orderID, selectedProduct.ProductID, selectedProduct.Quantity, product.Price)
		if err != nil {
			log.Printf("OrderRepository.CreateOrderProduct internal error %s", err.Error())
			return Order{}
		}

		orderService.CartRepository.DeleteCart(uid, selectedProduct.ProductID)
	}

	if submitedOrder.BurnPoint > 0 {
		orderService.OrderBurnPoint(uid, submitedOrder.BurnPoint)
	}

	return Order{
		OrderID: orderID,
	}
}

func (orderService OrderService) OrderBurnPoint(uid int, burn int) point.Point {
	submit := point.Point{
		OrgID:  1,
		UserID: uid,
		Amount: -(burn),
	}

	totalPoint, err := orderService.PointGateway.CreatePoint(uid, submit)
	if err != nil {
		log.Printf("orderService.PointService.DeductPoint internal error %s", err.Error())
		return point.Point{}
	}
	return totalPoint
}

func SendNotification(orderID int, trackingNumber string, dateTime time.Time) string {
	return fmt.Sprintf("วันเวลาที่ชำระเงิน %s หมายเลขคำสั่งซื้อ %d คุณสามารถติดตามสินค้าผ่านช่องทาง xx หมายเลข %s", dateTime.Format("2/1/2006 15:04:05"), orderID, trackingNumber)
}
