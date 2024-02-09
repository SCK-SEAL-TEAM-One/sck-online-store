package order

import (
	"fmt"
	"log"
	"store-service/internal/product"
	"time"
)

type OrderService struct {
	ProductRepository product.ProductRepository
	OrderRepository   OrderRepository
}

type OrderInterface interface {
	CreateOrder(submitedOrder SubmitedOrder) Order
}

type ProductRepository interface {
	GetProductByID(id int) product.ProductDetail
}

func (orderService OrderService) CreateOrder(submitedOrder SubmitedOrder) Order {
	orderID, err := orderService.OrderRepository.CreateOrder(submitedOrder)
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
	_, err = orderService.OrderRepository.CreateShipping(orderID, shippingInfo)
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
	}
	return Order{
		OrderID: orderID,
	}
	TODO
	update point & clear cart
}

// func (orderService OrderService) GetTotalProductPrice(submitedOrder SubmitedOrder) float64 {
// 	totalProductPrice := 0.00
// 	for _, cartItem := range submitedOrder.Cart {
// 		product, _ := orderService.ProductRepository.GetProductByID(cartItem.ProductID)
// 		totalProductPrice += product.Price * float64(cartItem.Quantity)
// 	}
// 	return totalProductPrice
// }

// func (orderService OrderService) GetTotalAmount(order SubmitedOrder) float64 {
// 	return orderService.GetTotalProductPrice(order) + order.GetShippingFee()
// }

func SendNotification(orderID int, trackingNumber string, dateTime time.Time) string {
	return fmt.Sprintf("วันเวลาที่ชำระเงิน %s หมายเลขคำสั่งซื้อ %d คุณสามารถติดตามสินค้าผ่านช่องทาง xx หมายเลข %s", dateTime.Format("2/1/2006 15:04:05"), orderID, trackingNumber)
}
