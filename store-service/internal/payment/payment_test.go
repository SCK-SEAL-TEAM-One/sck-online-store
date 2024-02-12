package payment_test

import (
	"errors"
	"store-service/internal/order"
	"store-service/internal/payment"
	"store-service/internal/shipping"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_ConfirmPayment_Input_OrderID_8004359103_Should_Be_Return_TrackingNumber_KR_307676366_No_Error(t *testing.T) {
	expected := payment.SubmitedPaymentResponse{
		OrderID:          8004359103,
		PaymentDate:      time.Now(),
		ShippingMethodID: 1,
		TrackingNumber:   "KR-307676366",
	}

	uid := 1
	oid := 8004359103
	orgID := 1
	shippingMethodID := 1
	paymentMethodID := 1

	mockOrderRepository := new(mockOrderRepository)
	mockOrderRepository.On("GetOrderByID", oid).Return(order.OrderDetail{
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
	mockBankGateway.On("GetCardDetail", orgID, uid).Return(payment.CardDetail{
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
	mockBankGateway.On("Payment", paymentDetail).Return("TRANSACTION_ID", nil)

	mockOrderRepository.On("GetOrderProduct", oid).Return([]order.OrderProduct{
		{
			ProductID: 2,
			Quantity:  2,
		},
	}, nil)

	mockProductRepository := new(mockProductRepository)
	mockProductRepository.On("UpdateStock", 2, 2).Return(nil)

	mockOrderRepository.On("UpdateOrderTransaction", oid, "TRANSACTION_ID").Return(nil)

	mockShippingGateway := new(mockShippingGateway)
	mockShippingGateway.On("GetTrackingNumber", shipping.ShippingGatewaySubmit{
		ShippingMethodID: shippingMethodID,
	}).Return("KR-307676366", nil)

	paymentService := payment.PaymentService{
		BankGateway:       mockBankGateway,
		ShippingGateway:   mockShippingGateway,
		OrderRepository:   mockOrderRepository,
		ProductRepository: mockProductRepository,
	}

	submitedPayment := payment.SubmitedPayment{
		OrderID: oid,
		OTP:     123456,
		RefOTP:  "REF_OTP",
	}

	actual, err := paymentService.ConfirmPayment(uid, submitedPayment)
	assert.Equal(t, expected.OrderID, actual.OrderID)
	assert.Equal(t, expected.ShippingMethodID, actual.ShippingMethodID)
	assert.Equal(t, expected.TrackingNumber, actual.TrackingNumber)
	assert.Equal(t, nil, err)
}

func Test_ConfirmPayment_Input_OrderID_8004359103_Should_Be_Return_OrderRepository_GetOrderByID_Error(t *testing.T) {
	expected := payment.SubmitedPaymentResponse{}

	uid := 1
	oid := 8004359103

	mockOrderRepository := new(mockOrderRepository)
	mockOrderRepository.On("GetOrderByID", oid).Return(order.OrderDetail{}, errors.New("GetOrderByID Error"))

	paymentService := payment.PaymentService{
		OrderRepository: mockOrderRepository,
	}

	submitedPayment := payment.SubmitedPayment{
		OrderID: oid,
		OTP:     123456,
		RefOTP:  "REF_OTP",
	}

	actual, err := paymentService.ConfirmPayment(uid, submitedPayment)

	assert.Equal(t, expected, actual)
	assert.NotNil(t, err)
}

func Test_ConfirmPayment_Input_OrderID_8004359103_Should_Be_Return_BankGateway_GetCardDetail_Error(t *testing.T) {
	expected := payment.SubmitedPaymentResponse{}

	uid := 1
	oid := 8004359103
	orgID := 1
	shippingMethodID := 1
	paymentMethodID := 1

	mockOrderRepository := new(mockOrderRepository)
	mockOrderRepository.On("GetOrderByID", oid).Return(order.OrderDetail{
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
	mockBankGateway.On("GetCardDetail", orgID, uid).Return(payment.CardDetail{}, errors.New("GetCardDetail Error"))

	paymentService := payment.PaymentService{
		BankGateway:     mockBankGateway,
		OrderRepository: mockOrderRepository,
	}

	submitedPayment := payment.SubmitedPayment{
		OrderID: oid,
		OTP:     123456,
		RefOTP:  "REF_OTP",
	}

	actual, err := paymentService.ConfirmPayment(uid, submitedPayment)

	assert.Equal(t, expected, actual)
	assert.NotNil(t, err)
}

func Test_ConfirmPayment_Input_OrderID_8004359103_Should_Be_Return_BankGateway_Payment_Error(t *testing.T) {
	expected := payment.SubmitedPaymentResponse{}

	uid := 1
	oid := 8004359103
	orgID := 1
	shippingMethodID := 1
	paymentMethodID := 1

	mockOrderRepository := new(mockOrderRepository)
	mockOrderRepository.On("GetOrderByID", oid).Return(order.OrderDetail{
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
	mockBankGateway.On("GetCardDetail", orgID, uid).Return(payment.CardDetail{
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
	mockBankGateway.On("Payment", paymentDetail).Return("", errors.New("Payment Error"))

	paymentService := payment.PaymentService{
		BankGateway:     mockBankGateway,
		OrderRepository: mockOrderRepository,
	}

	submitedPayment := payment.SubmitedPayment{
		OrderID: oid,
		OTP:     123456,
		RefOTP:  "REF_OTP",
	}

	actual, err := paymentService.ConfirmPayment(uid, submitedPayment)

	assert.Equal(t, expected, actual)
	assert.NotNil(t, err)
}

func Test_ConfirmPayment_Input_OrderID_8004359103_Should_Be_Return_OrderRepository_GetOrderProduct_Error(t *testing.T) {
	expected := payment.SubmitedPaymentResponse{}

	uid := 1
	oid := 8004359103
	orgID := 1
	shippingMethodID := 1
	paymentMethodID := 1

	mockOrderRepository := new(mockOrderRepository)
	mockOrderRepository.On("GetOrderByID", oid).Return(order.OrderDetail{
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
	mockBankGateway.On("GetCardDetail", orgID, uid).Return(payment.CardDetail{
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
	mockBankGateway.On("Payment", paymentDetail).Return("TRANSACTION_ID", nil)

	mockOrderRepository.On("GetOrderProduct", oid).Return([]order.OrderProduct{}, errors.New("GetOrderProduct Error"))

	paymentService := payment.PaymentService{
		BankGateway:     mockBankGateway,
		OrderRepository: mockOrderRepository,
	}

	submitedPayment := payment.SubmitedPayment{
		OrderID: oid,
		OTP:     123456,
		RefOTP:  "REF_OTP",
	}

	actual, err := paymentService.ConfirmPayment(uid, submitedPayment)

	assert.Equal(t, expected, actual)
	assert.NotNil(t, err)
}

func Test_ConfirmPayment_Input_OrderID_8004359103_Should_Be_Return_ProductRepository_UpdateStock_Error(t *testing.T) {
	expected := payment.SubmitedPaymentResponse{}

	uid := 1
	oid := 8004359103
	orgID := 1
	shippingMethodID := 1
	paymentMethodID := 1

	mockOrderRepository := new(mockOrderRepository)
	mockOrderRepository.On("GetOrderByID", oid).Return(order.OrderDetail{
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
	mockBankGateway.On("GetCardDetail", orgID, uid).Return(payment.CardDetail{
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
	mockBankGateway.On("Payment", paymentDetail).Return("TRANSACTION_ID", nil)

	mockOrderRepository.On("GetOrderProduct", oid).Return([]order.OrderProduct{
		{
			ProductID: 2,
			Quantity:  2,
		},
	}, nil)

	mockProductRepository := new(mockProductRepository)
	mockProductRepository.On("UpdateStock", 2, 2).Return(errors.New("UpdateStock Error"))

	paymentService := payment.PaymentService{
		BankGateway:       mockBankGateway,
		OrderRepository:   mockOrderRepository,
		ProductRepository: mockProductRepository,
	}

	submitedPayment := payment.SubmitedPayment{
		OrderID: oid,
		OTP:     123456,
		RefOTP:  "REF_OTP",
	}

	actual, err := paymentService.ConfirmPayment(uid, submitedPayment)

	assert.Equal(t, expected, actual)
	assert.NotNil(t, err)
}

func Test_ConfirmPayment_Input_OrderID_8004359103_Should_Be_Return_OrderRepository_UpdateOrderTransaction_Error(t *testing.T) {
	expected := payment.SubmitedPaymentResponse{}

	uid := 1
	oid := 8004359103
	orgID := 1
	shippingMethodID := 1
	paymentMethodID := 1

	mockOrderRepository := new(mockOrderRepository)
	mockOrderRepository.On("GetOrderByID", oid).Return(order.OrderDetail{
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
	mockBankGateway.On("GetCardDetail", orgID, uid).Return(payment.CardDetail{
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
	mockBankGateway.On("Payment", paymentDetail).Return("TRANSACTION_ID", nil)

	mockOrderRepository.On("GetOrderProduct", oid).Return([]order.OrderProduct{
		{
			ProductID: 2,
			Quantity:  2,
		},
	}, nil)

	mockProductRepository := new(mockProductRepository)
	mockProductRepository.On("UpdateStock", 2, 2).Return(nil)

	mockOrderRepository.On("UpdateOrderTransaction", oid, "TRANSACTION_ID").Return(errors.New("UpdateOrderTransaction Error"))

	paymentService := payment.PaymentService{
		BankGateway:       mockBankGateway,
		OrderRepository:   mockOrderRepository,
		ProductRepository: mockProductRepository,
	}

	submitedPayment := payment.SubmitedPayment{
		OrderID: oid,
		OTP:     123456,
		RefOTP:  "REF_OTP",
	}

	actual, err := paymentService.ConfirmPayment(uid, submitedPayment)

	assert.Equal(t, expected, actual)
	assert.NotNil(t, err)
}

func Test_ConfirmPayment_Input_OrderID_8004359103_Should_Be_Return_ShippingGateway_GetTrackingNumber_Error(t *testing.T) {
	expected := payment.SubmitedPaymentResponse{}

	uid := 1
	oid := 8004359103
	orgID := 1
	shippingMethodID := 1
	paymentMethodID := 1

	mockOrderRepository := new(mockOrderRepository)
	mockOrderRepository.On("GetOrderByID", oid).Return(order.OrderDetail{
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
	mockBankGateway.On("GetCardDetail", orgID, uid).Return(payment.CardDetail{
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
	mockBankGateway.On("Payment", paymentDetail).Return("TRANSACTION_ID", nil)

	mockOrderRepository.On("GetOrderProduct", oid).Return([]order.OrderProduct{
		{
			ProductID: 2,
			Quantity:  2,
		},
	}, nil)

	mockProductRepository := new(mockProductRepository)
	mockProductRepository.On("UpdateStock", 2, 2).Return(nil)

	mockOrderRepository.On("UpdateOrderTransaction", oid, "TRANSACTION_ID").Return(nil)

	mockShippingGateway := new(mockShippingGateway)
	mockShippingGateway.On("GetTrackingNumber", shipping.ShippingGatewaySubmit{
		ShippingMethodID: shippingMethodID,
	}).Return("", errors.New("GetTrackingNumber Error"))

	paymentService := payment.PaymentService{
		BankGateway:       mockBankGateway,
		ShippingGateway:   mockShippingGateway,
		OrderRepository:   mockOrderRepository,
		ProductRepository: mockProductRepository,
	}

	submitedPayment := payment.SubmitedPayment{
		OrderID: oid,
		OTP:     123456,
		RefOTP:  "REF_OTP",
	}

	actual, err := paymentService.ConfirmPayment(uid, submitedPayment)

	assert.Equal(t, expected, actual)
	assert.NotNil(t, err)
}
