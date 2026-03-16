package payment_test

import (
	"context"
	"errors"
	"store-service/internal/order"
	"store-service/internal/payment"
	"store-service/internal/shipping"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_ConfirmPayment_Input_OrderNumber_2603159522001_Should_Be_Return_TrackingNumber_KR_307676366_No_Error(t *testing.T) {
	uid := 1
	orderNumber := "2603159522001"
	oid := 8004359103
	orgID := 1
	shippingMethodID := 1
	paymentMethodID := 1
	trackingNumber := "KR-307676366"

	expected := payment.SubmitedPaymentResponse{
		OrderNumber:      "2603159522001",
		PaymentDate:      time.Now(),
		ShippingMethodID: 1,
		TrackingNumber:   trackingNumber,
	}

	mockOrderRepository := new(mockOrderRepository)
	mockOrderRepository.On("GetOrderByOrderNumber", mock.Anything, orderNumber).Return(order.OrderDetail{
		ID:               oid,
		OrderNumber:      orderNumber,
		UserID:           uid,
		ShippingMethodID: shippingMethodID,
		PaymentMethodID:  paymentMethodID,
		BurnPoint:        0,
		SubTotalPrice:    100.00,
		DiscountPrice:    10.00,
		TotalPrice:       90.00,
		TransactionID:    "",
		Status:           "created",
	}, nil)

	mockBankGateway := new(mockBankGateway)
	mockBankGateway.On("GetCardDetail", mock.Anything, orgID, uid).Return(payment.CardDetail{
		CardNumber:   "4719700591590995",
		CVV:          752,
		ExpiredMonth: 12,
		ExpiredYear:  27,
		CardName:     "SCK ShuHaRi",
	}, nil)

	paymentDetail := payment.PaymentDetail{
		CardNumber:   "4719700591590995",
		CVV:          752,
		ExpiredMonth: 12,
		ExpiredYear:  27,
		CardName:     "SCK ShuHaRi",
		Amount:       90.00,
		Currency:     "USD",
		MerchantID:   1,
	}
	mockBankGateway.On("Payment", mock.Anything, paymentDetail).Return("TRANSACTION_ID", nil)

	mockOrderRepository.On("GetOrderProduct", mock.Anything, oid).Return([]order.OrderProduct{
		{
			ProductID: 2,
			Quantity:  2,
		},
	}, nil)

	mockProductRepository := new(mockProductRepository)
	mockProductRepository.On("UpdateStock", mock.Anything, 2, 2).Return(nil)

	mockOrderRepository.On("UpdateOrderTransaction", mock.Anything, oid, "TRANSACTION_ID").Return(nil)

	mockShippingGateway := new(mockShippingGateway)
	mockShippingGateway.On("GetTrackingNumber", mock.Anything, shipping.ShippingGatewaySubmit{
		ShippingMethodID: shippingMethodID,
	}).Return(trackingNumber, nil)
	mockOrderRepository.On("UpdateOrderTrackingNumber", mock.Anything, oid, trackingNumber).Return(nil)

	paymentService := payment.PaymentService{
		BankGateway:       mockBankGateway,
		ShippingGateway:   mockShippingGateway,
		OrderRepository:   mockOrderRepository,
		ProductRepository: mockProductRepository,
	}

	submitedPayment := payment.SubmitedPayment{
		OrderNumber: orderNumber,
		OTP:         123456,
		RefOTP:      "REF_OTP",
	}

	actual, err := paymentService.ConfirmPayment(context.Background(), uid, submitedPayment)
	assert.Equal(t, expected.OrderNumber, actual.OrderNumber)
	assert.Equal(t, expected.ShippingMethodID, actual.ShippingMethodID)
	assert.Equal(t, expected.TrackingNumber, actual.TrackingNumber)
	assert.Equal(t, nil, err)
}

func Test_ConfirmPayment_Input_OrderNumber_2603159533002_Should_Be_Return_OrderRepository_GetOrderByOrderNumber_Error(t *testing.T) {
	expected := payment.SubmitedPaymentResponse{}

	uid := 1
	orderNumber := "2603159533002"

	mockOrderRepository := new(mockOrderRepository)
	mockOrderRepository.On("GetOrderByOrderNumber", mock.Anything, orderNumber).Return(order.OrderDetail{}, errors.New("GetOrderByOrderNumber Error"))

	paymentService := payment.PaymentService{
		OrderRepository: mockOrderRepository,
	}

	submitedPayment := payment.SubmitedPayment{
		OrderNumber: orderNumber,
		OTP:         123456,
		RefOTP:      "REF_OTP",
	}

	actual, err := paymentService.ConfirmPayment(context.Background(), uid, submitedPayment)

	assert.Equal(t, expected, actual)
	assert.NotNil(t, err)
}

func Test_ConfirmPayment_Input_OrderNumber_2603159544003_Should_Be_Return_BankGateway_GetCardDetail_Error(t *testing.T) {
	expected := payment.SubmitedPaymentResponse{}

	uid := 1
	oid := 8004359103
	orgID := 1
	shippingMethodID := 1
	paymentMethodID := 1
	orderNumber := "2603159544003"

	mockOrderRepository := new(mockOrderRepository)
	mockOrderRepository.On("GetOrderByOrderNumber", mock.Anything, orderNumber).Return(order.OrderDetail{
		ID:               oid,
		OrderNumber:      orderNumber,
		UserID:           uid,
		ShippingMethodID: shippingMethodID,
		PaymentMethodID:  paymentMethodID,
		BurnPoint:        0,
		SubTotalPrice:    100.00,
		DiscountPrice:    10.00,
		TotalPrice:       90.00,
		TransactionID:    "",
		Status:           "created",
	}, nil)

	mockBankGateway := new(mockBankGateway)
	mockBankGateway.On("GetCardDetail", mock.Anything, orgID, uid).Return(payment.CardDetail{}, errors.New("GetCardDetail Error"))

	paymentService := payment.PaymentService{
		BankGateway:     mockBankGateway,
		OrderRepository: mockOrderRepository,
	}

	submitedPayment := payment.SubmitedPayment{
		OrderNumber: orderNumber,
		OTP:         123456,
		RefOTP:      "REF_OTP",
	}

	actual, err := paymentService.ConfirmPayment(context.Background(), uid, submitedPayment)

	assert.Equal(t, expected, actual)
	assert.NotNil(t, err)
}

func Test_ConfirmPayment_Input_OrderNumber_2603159822004_Should_Be_Return_BankGateway_Payment_Error(t *testing.T) {
	expected := payment.SubmitedPaymentResponse{}

	uid := 1
	oid := 8004359103
	orgID := 1
	shippingMethodID := 1
	paymentMethodID := 1
	orderNumber := "2603159822004"

	mockOrderRepository := new(mockOrderRepository)
	mockOrderRepository.On("GetOrderByOrderNumber", mock.Anything, orderNumber).Return(order.OrderDetail{
		ID:               oid,
		OrderNumber:      orderNumber,
		UserID:           uid,
		ShippingMethodID: shippingMethodID,
		PaymentMethodID:  paymentMethodID,
		BurnPoint:        0,
		SubTotalPrice:    100.00,
		DiscountPrice:    10.00,
		TotalPrice:       90.00,
		TransactionID:    "",
		Status:           "created",
	}, nil)

	mockBankGateway := new(mockBankGateway)
	mockBankGateway.On("GetCardDetail", mock.Anything, orgID, uid).Return(payment.CardDetail{
		CardNumber:   "4719700591590995",
		CVV:          752,
		ExpiredMonth: 12,
		ExpiredYear:  27,
		CardName:     "SCK ShuHaRi",
	}, nil)

	paymentDetail := payment.PaymentDetail{
		CardNumber:   "4719700591590995",
		CVV:          752,
		ExpiredMonth: 12,
		ExpiredYear:  27,
		CardName:     "SCK ShuHaRi",
		Amount:       90.00,
		Currency:     "USD",
		MerchantID:   1,
	}
	mockBankGateway.On("Payment", mock.Anything, paymentDetail).Return("", errors.New("Payment Error"))

	paymentService := payment.PaymentService{
		BankGateway:     mockBankGateway,
		OrderRepository: mockOrderRepository,
	}

	submitedPayment := payment.SubmitedPayment{
		OrderNumber: orderNumber,
		OTP:         123456,
		RefOTP:      "REF_OTP",
	}

	actual, err := paymentService.ConfirmPayment(context.Background(), uid, submitedPayment)

	assert.Equal(t, expected, actual)
	assert.NotNil(t, err)
}

func Test_ConfirmPayment_Input_OrderNumber_2603159833005_Should_Be_Return_OrderRepository_GetOrderProduct_Error(t *testing.T) {
	expected := payment.SubmitedPaymentResponse{}

	uid := 1
	oid := 8004359103
	orgID := 1
	shippingMethodID := 1
	paymentMethodID := 1
	orderNumber := "2603159833005"

	mockOrderRepository := new(mockOrderRepository)
	mockOrderRepository.On("GetOrderByOrderNumber", mock.Anything, orderNumber).Return(order.OrderDetail{
		ID:               oid,
		OrderNumber:      orderNumber,
		UserID:           uid,
		ShippingMethodID: shippingMethodID,
		PaymentMethodID:  paymentMethodID,
		BurnPoint:        0,
		SubTotalPrice:    100.00,
		DiscountPrice:    10.00,
		TotalPrice:       90.00,
		TransactionID:    "",
		Status:           "created",
	}, nil)

	mockBankGateway := new(mockBankGateway)
	mockBankGateway.On("GetCardDetail", mock.Anything, orgID, uid).Return(payment.CardDetail{
		CardNumber:   "4719700591590995",
		CVV:          752,
		ExpiredMonth: 12,
		ExpiredYear:  27,
		CardName:     "SCK ShuHaRi",
	}, nil)

	paymentDetail := payment.PaymentDetail{
		CardNumber:   "4719700591590995",
		CVV:          752,
		ExpiredMonth: 12,
		ExpiredYear:  27,
		CardName:     "SCK ShuHaRi",
		Amount:       90.00,
		Currency:     "USD",
		MerchantID:   1,
	}
	mockBankGateway.On("Payment", mock.Anything, paymentDetail).Return("TRANSACTION_ID", nil)

	mockOrderRepository.On("GetOrderProduct", mock.Anything, oid).Return([]order.OrderProduct{}, errors.New("GetOrderProduct Error"))

	paymentService := payment.PaymentService{
		BankGateway:     mockBankGateway,
		OrderRepository: mockOrderRepository,
	}

	submitedPayment := payment.SubmitedPayment{
		OrderNumber: orderNumber,
		OTP:         123456,
		RefOTP:      "REF_OTP",
	}

	actual, err := paymentService.ConfirmPayment(context.Background(), uid, submitedPayment)

	assert.Equal(t, expected, actual)
	assert.NotNil(t, err)
}

func Test_ConfirmPayment_Input_OrderNumber_2603159844006_Should_Be_Return_ProductRepository_UpdateStock_Error(t *testing.T) {
	expected := payment.SubmitedPaymentResponse{}

	uid := 1
	oid := 8004359103
	orgID := 1
	shippingMethodID := 1
	paymentMethodID := 1
	orderNumber := "2603159844006"

	mockOrderRepository := new(mockOrderRepository)
	mockOrderRepository.On("GetOrderByOrderNumber", mock.Anything, orderNumber).Return(order.OrderDetail{
		ID:               oid,
		OrderNumber:      orderNumber,
		UserID:           uid,
		ShippingMethodID: shippingMethodID,
		PaymentMethodID:  paymentMethodID,
		BurnPoint:        0,
		SubTotalPrice:    100.00,
		DiscountPrice:    10.00,
		TotalPrice:       90.00,
		TransactionID:    "",
		Status:           "created",
	}, nil)

	mockBankGateway := new(mockBankGateway)
	mockBankGateway.On("GetCardDetail", mock.Anything, orgID, uid).Return(payment.CardDetail{
		CardNumber:   "4719700591590995",
		CVV:          752,
		ExpiredMonth: 12,
		ExpiredYear:  27,
		CardName:     "SCK ShuHaRi",
	}, nil)

	paymentDetail := payment.PaymentDetail{
		CardNumber:   "4719700591590995",
		CVV:          752,
		ExpiredMonth: 12,
		ExpiredYear:  27,
		CardName:     "SCK ShuHaRi",
		Amount:       90.00,
		Currency:     "USD",
		MerchantID:   1,
	}
	mockBankGateway.On("Payment", mock.Anything, paymentDetail).Return("TRANSACTION_ID", nil)

	mockOrderRepository.On("GetOrderProduct", mock.Anything, oid).Return([]order.OrderProduct{
		{
			ProductID: 2,
			Quantity:  2,
		},
	}, nil)

	mockProductRepository := new(mockProductRepository)
	mockProductRepository.On("UpdateStock", mock.Anything, 2, 2).Return(errors.New("UpdateStock Error"))

	paymentService := payment.PaymentService{
		BankGateway:       mockBankGateway,
		OrderRepository:   mockOrderRepository,
		ProductRepository: mockProductRepository,
	}

	submitedPayment := payment.SubmitedPayment{
		OrderNumber: orderNumber,
		OTP:         123456,
		RefOTP:      "REF_OTP",
	}

	actual, err := paymentService.ConfirmPayment(context.Background(), uid, submitedPayment)

	assert.Equal(t, expected, actual)
	assert.NotNil(t, err)
}

func Test_ConfirmPayment_Input_OrderNumber_2603159522179_Should_Be_Return_OrderRepository_UpdateOrderTransaction_Error(t *testing.T) {
	expected := payment.SubmitedPaymentResponse{}

	uid := 1
	oid := 8004359103
	orgID := 1
	shippingMethodID := 1
	paymentMethodID := 1
	orderNumber := "2603159522179"

	mockOrderRepository := new(mockOrderRepository)
	mockOrderRepository.On("GetOrderByOrderNumber", mock.Anything, orderNumber).Return(order.OrderDetail{
		ID:               oid,
		UserID:           uid,
		ShippingMethodID: shippingMethodID,
		PaymentMethodID:  paymentMethodID,
		BurnPoint:        0,
		SubTotalPrice:    100.00,
		DiscountPrice:    10.00,
		TotalPrice:       90.00,
		TransactionID:    "",
		Status:           "created",
	}, nil)

	mockBankGateway := new(mockBankGateway)
	mockBankGateway.On("GetCardDetail", mock.Anything, orgID, uid).Return(payment.CardDetail{
		CardNumber:   "4719700591590995",
		CVV:          752,
		ExpiredMonth: 12,
		ExpiredYear:  27,
		CardName:     "SCK ShuHaRi",
	}, nil)

	paymentDetail := payment.PaymentDetail{
		CardNumber:   "4719700591590995",
		CVV:          752,
		ExpiredMonth: 12,
		ExpiredYear:  27,
		CardName:     "SCK ShuHaRi",
		Amount:       90.00,
		Currency:     "USD",
		MerchantID:   1,
	}
	mockBankGateway.On("Payment", mock.Anything, paymentDetail).Return("TRANSACTION_ID", nil)

	mockOrderRepository.On("GetOrderProduct", mock.Anything, oid).Return([]order.OrderProduct{
		{
			ProductID: 2,
			Quantity:  2,
		},
	}, nil)

	mockProductRepository := new(mockProductRepository)
	mockProductRepository.On("UpdateStock", mock.Anything, 2, 2).Return(nil)

	mockOrderRepository.On("UpdateOrderTransaction", mock.Anything, oid, "TRANSACTION_ID").Return(errors.New("UpdateOrderTransaction Error"))

	paymentService := payment.PaymentService{
		BankGateway:       mockBankGateway,
		OrderRepository:   mockOrderRepository,
		ProductRepository: mockProductRepository,
	}

	submitedPayment := payment.SubmitedPayment{
		OrderNumber: orderNumber,
		OTP:         123456,
		RefOTP:      "REF_OTP",
	}

	actual, err := paymentService.ConfirmPayment(context.Background(), uid, submitedPayment)

	assert.Equal(t, expected, actual)
	assert.NotNil(t, err)
}

func Test_ConfirmPayment_Input_OrderNumber_2603159533899_Should_Be_Return_ShippingGateway_GetTrackingNumber_Error(t *testing.T) {
	expected := payment.SubmitedPaymentResponse{}

	uid := 1
	oid := 8004359103
	orgID := 1
	shippingMethodID := 1
	paymentMethodID := 1
	orderNumber := "2603159533899"

	mockOrderRepository := new(mockOrderRepository)
	mockOrderRepository.On("GetOrderByOrderNumber", mock.Anything, orderNumber).Return(order.OrderDetail{
		ID:               oid,
		OrderNumber:      orderNumber,
		UserID:           uid,
		ShippingMethodID: shippingMethodID,
		PaymentMethodID:  paymentMethodID,
		BurnPoint:        0,
		SubTotalPrice:    100.00,
		DiscountPrice:    10.00,
		TotalPrice:       90.00,
		TransactionID:    "",
		Status:           "created",
	}, nil)

	mockBankGateway := new(mockBankGateway)
	mockBankGateway.On("GetCardDetail", mock.Anything, orgID, uid).Return(payment.CardDetail{
		CardNumber:   "4719700591590995",
		CVV:          752,
		ExpiredMonth: 12,
		ExpiredYear:  27,
		CardName:     "SCK ShuHaRi",
	}, nil)

	paymentDetail := payment.PaymentDetail{
		CardNumber:   "4719700591590995",
		CVV:          752,
		ExpiredMonth: 12,
		ExpiredYear:  27,
		CardName:     "SCK ShuHaRi",
		Amount:       90.00,
		Currency:     "USD",
		MerchantID:   1,
	}
	mockBankGateway.On("Payment", mock.Anything, paymentDetail).Return("TRANSACTION_ID", nil)

	mockOrderRepository.On("GetOrderProduct", mock.Anything, oid).Return([]order.OrderProduct{
		{
			ProductID: 2,
			Quantity:  2,
		},
	}, nil)

	mockProductRepository := new(mockProductRepository)
	mockProductRepository.On("UpdateStock", mock.Anything, 2, 2).Return(nil)

	mockOrderRepository.On("UpdateOrderTransaction", mock.Anything, oid, "TRANSACTION_ID").Return(nil)

	mockShippingGateway := new(mockShippingGateway)
	mockShippingGateway.On("GetTrackingNumber", mock.Anything, shipping.ShippingGatewaySubmit{
		ShippingMethodID: shippingMethodID,
	}).Return("", errors.New("GetTrackingNumber Error"))

	paymentService := payment.PaymentService{
		BankGateway:       mockBankGateway,
		ShippingGateway:   mockShippingGateway,
		OrderRepository:   mockOrderRepository,
		ProductRepository: mockProductRepository,
	}

	submitedPayment := payment.SubmitedPayment{
		OrderNumber: orderNumber,
		OTP:         123456,
		RefOTP:      "REF_OTP",
	}

	actual, err := paymentService.ConfirmPayment(context.Background(), uid, submitedPayment)

	assert.Equal(t, expected, actual)
	assert.NotNil(t, err)
}
