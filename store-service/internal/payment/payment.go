package payment

import (
	"database/sql"
	"log"
	"store-service/internal/order"
	"store-service/internal/product"
	"store-service/internal/shipping"
	"time"
)

type PaymentInterface interface {
	ConfirmPayment(uid int, submitedPayment SubmitedPayment) (SubmitedPaymentResponse, error)
}

type PaymentService struct {
	BankGateway       BankGatewayInterface
	ShippingGateway   ShippingGatewayInterface
	OrderRepository   order.OrderRepository
	ProductRepository product.ProductRepository
}

type BankGatewayInterface interface {
	Payment(paymentDetail PaymentDetail) (string, error)
	GetCardDetail(orgID int, userID int) (CardDetail, error)
}

type ShippingGatewayInterface interface {
	GetTrackingNumber(shippingGatewaySubmit shipping.ShippingGatewaySubmit) (string, error)
}

func (service PaymentService) ConfirmPayment(uid int, submitedPayment SubmitedPayment) (SubmitedPaymentResponse, error) {
	orgID := 1
	orderNumber := submitedPayment.OrderNumber
	currency := "USD"
	now := time.Now()

	orderDetail, err := service.OrderRepository.GetOrderByOrderNumber(orderNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("OrderRepository.GetOrderByOrderNumber not found for Order Number %s: %s", orderNumber, err.Error())
			return SubmitedPaymentResponse{}, order.ErrOrderNotFound
		}
		log.Printf("OrderRepository.GetOrderByOrderNumber internal error %s", err.Error())
		return SubmitedPaymentResponse{}, err
	}

	cardDetail, err := service.BankGateway.GetCardDetail(orgID, uid)
	if err != nil {
		log.Printf("BankGateway.GetCardDetail internal error %s", err.Error())
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
	transactionId, err := service.BankGateway.Payment(paymentdetail)
	if err != nil {
		log.Printf("BankGateway.Payment internal error %s", err.Error())
		return SubmitedPaymentResponse{}, err
	}

	orderProductList, err := service.OrderRepository.GetOrderProduct(orderDetail.ID)
	if err != nil {
		log.Printf("OrderRepository.GetOrderProduct internal error %s", err.Error())
		return SubmitedPaymentResponse{}, err
	}
	for _, orderProduct := range orderProductList {
		err = service.ProductRepository.UpdateStock(orderProduct.ProductID, orderProduct.Quantity)
		if err != nil {
			log.Printf("ProductRepository.UpdateStock internal error %s", err.Error())
			return SubmitedPaymentResponse{}, err
		}
	}

	err = service.OrderRepository.UpdateOrderTransaction(orderDetail.ID, transactionId)
	if err != nil {
		log.Printf("OrderRepository.UpdateOrderTransaction internal error %s", err.Error())
		return SubmitedPaymentResponse{}, err
	}

	shippingGatewaySubmit := shipping.ShippingGatewaySubmit{
		ShippingMethodID: orderDetail.ShippingMethodID,
	}
	trackingNumber, err := service.ShippingGateway.GetTrackingNumber(shippingGatewaySubmit)
	if err != nil {
		log.Printf("ShippingGateway.GetTrackingNumber internal error %s", err.Error())
		return SubmitedPaymentResponse{}, err
	}

	err = service.OrderRepository.UpdateOrderTrackingNumber(orderDetail.ID, trackingNumber)
	if err != nil {
		log.Printf("OrderRepository.UpdateOrderTrackingNumber internal error %s", err.Error())
		return SubmitedPaymentResponse{}, err
	}

	return SubmitedPaymentResponse{
		OrderNumber:      orderDetail.OrderNumber,
		PaymentDate:      now,
		ShippingMethodID: orderDetail.ShippingMethodID,
		TrackingNumber:   trackingNumber,
	}, nil
}
