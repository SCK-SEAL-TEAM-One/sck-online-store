package order_test

import (
	"database/sql"
	"errors"
	"fmt"
	"store-service/internal/auth"
	"store-service/internal/order"
	"store-service/internal/point"
	"store-service/internal/product"
	"store-service/internal/shipping"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_CreateOrder_Input_Submitted_Order_Should_be_OrderNumber_2601069522001(t *testing.T) {
	uid := 1
	oid := 8004359103
	orderNumber := "2601069522001"
	productPrice := 12.95
	nextSeq := 1
	fixedTime := time.Date(2026, 1, 6, 0, 0, 0, 0, time.UTC)
	yearPrefix := "26"

	expected := order.Order{
		OrderNumber: orderNumber,
	}

	submittedOrder := order.SubmitedOrder{
		Cart: []order.OrderProduct{
			{
				ProductID: 2,
				Quantity:  1,
			},
		},
		ShippingMethodID:     1,
		ShippingAddress:      "405/37 ถ.มหิดล",
		ShippingSubDistrict:  "ท่าศาลา",
		ShippingDistrict:     "เมือง",
		ShippingProvince:     "เชียงใหม่",
		ShippingZipCode:      "50000",
		RecipientFirstName:   "ณัฐญา",
		RecipientLastName:    "ชุติบุตร",
		RecipientPhoneNumber: "0970809292",
		PaymentMethodID:      1,
		BurnPoint:            0,
		SubTotalPrice:        100.00,
		DiscountPrice:        0,
		TotalPrice:           100.00,
	}

	mockPointInterface := new(mockPointInterface)
	mockPointInterface.On("CheckBurnPoint", uid, submittedOrder.BurnPoint).Return(true, nil)

	mockProductRepository := new(mockProductRepository)
	mockProductRepository.On("GetProductByID", submittedOrder.Cart[0].ProductID).Return(product.ProductDetail{
		ID:           submittedOrder.Cart[0].ProductID,
		Name:         "43 Piece dinner Set",
		Price:        productPrice,
		PriceTHB:     0,
		PriceFullTHB: 0,
		Stock:        1,
		Brand:        "Coolkidz",
		Image:        "43_Piece_Dinner_Set.jpg",
	}, nil)

	mockShippingRepository := new(mockShippingRepository)
	mockShippingRepository.On("GetShippingMethodByID", submittedOrder.ShippingMethodID).Return(shipping.ShippingMethodDetail{
		ID:          1,
		Name:        "Kerry",
		Description: "",
		Fee:         50,
	}, nil)

	mockOrderRepository := new(mockOrderRepository)
	mockOrderRepository.On("GetLastOrderNumber", yearPrefix).Return("", sql.ErrNoRows)

	mockOrderHelper := new(mockOrderHelper)
	mockOrderHelper.On("GetNextSequence", "").Return(nextSeq, nil)
	mockOrderHelper.On("GenerateOrderNumber", submittedOrder.PaymentMethodID, submittedOrder.ShippingMethodID, nextSeq, fixedTime).Return(orderNumber, nil)

	orderDetail := order.OrderDetail{
		OrderNumber:      orderNumber,
		ShippingMethodID: submittedOrder.ShippingMethodID,
		PaymentMethodID:  submittedOrder.PaymentMethodID,
		SubTotalPrice:    465.811034,
		DiscountPrice:    0,
		TotalPrice:       515.8110340000001,
		ShippingFee:      50,
		BurnPoint:        0,
		EarnPoint:        4,
	}

	mockOrderRepository.On("CreateOrder", uid, orderDetail).Return(oid, nil)

	shippingInfo := order.ShippingInfo{
		ShippingMethodID:     1,
		ShippingAddress:      "405/37 ถ.มหิดล",
		ShippingSubDistrict:  "ท่าศาลา",
		ShippingDistrict:     "เมือง",
		ShippingProvince:     "เชียงใหม่",
		ShippingZipCode:      "50000",
		RecipientFirstName:   "ณัฐญา",
		RecipientLastName:    "ชุติบุตร",
		RecipientPhoneNumber: "0970809292",
	}
	mockOrderRepository.On("CreateShipping", uid, oid, shippingInfo).Return(1, nil)

	mockOrderRepository.On("CreateOrderProduct", oid, submittedOrder.Cart[0].ProductID, submittedOrder.Cart[0].Quantity, productPrice).Return(nil)

	mockCartRepository := new(mockCartRepository)
	mockCartRepository.On("DeleteCart", uid, submittedOrder.Cart[0].ProductID).Return(nil)

	orderService := order.OrderService{
		ProductRepository:  mockProductRepository,
		OrderRepository:    mockOrderRepository,
		CartRepository:     mockCartRepository,
		PointService:       mockPointInterface,
		ShippingRepository: mockShippingRepository,
		OrderHelper:        mockOrderHelper,
		Clock:              func() time.Time { return fixedTime },
	}

	actual, err := orderService.CreateOrder(uid, submittedOrder)

	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)
}

func Test_CreateOrder_Input_Submitted_Order_Should_be_Return_Error_Points_not_Enough(t *testing.T) {
	expected := order.Order{}
	expectedErr := fmt.Errorf("points are not enough, please try again")

	uid := 1
	burnPoint := 100

	mockPointInterface := new(mockPointInterface)
	mockPointInterface.On("CheckBurnPoint", uid, -(burnPoint)).Return(false, fmt.Errorf("points are not enough, please try again"))

	submittedOrder := order.SubmitedOrder{
		Cart: []order.OrderProduct{
			{
				ProductID: 2,
				Quantity:  1,
			},
		},
		ShippingMethodID:     1,
		ShippingAddress:      "405/37 ถ.มหิดล",
		ShippingSubDistrict:  "ท่าศาลา",
		ShippingDistrict:     "เมือง",
		ShippingProvince:     "เชียงใหม่",
		ShippingZipCode:      "50000",
		RecipientFirstName:   "ณัฐญา",
		RecipientLastName:    "ชุติบุตร",
		RecipientPhoneNumber: "0970809292",
		PaymentMethodID:      1,
		BurnPoint:            burnPoint,
		SubTotalPrice:        100.00,
		DiscountPrice:        0,
		TotalPrice:           100.00,
	}

	orderService := order.OrderService{
		PointService: mockPointInterface,
	}

	actual, err := orderService.CreateOrder(uid, submittedOrder)

	assert.Equal(t, expected, actual)
	assert.Equal(t, expectedErr, err)
}

func Test_CreateOrder_Input_Submitted_Order_Should_be_Return_No_Product_in_Order_Error(t *testing.T) {
	expected := order.Order{}
	expectedErr := fmt.Errorf("There is no product in order, please try again")

	uid := 1

	submittedOrder := order.SubmitedOrder{
		Cart:                 []order.OrderProduct{},
		ShippingMethodID:     1,
		ShippingAddress:      "405/37 ถ.มหิดล",
		ShippingSubDistrict:  "ท่าศาลา",
		ShippingDistrict:     "เมือง",
		ShippingProvince:     "เชียงใหม่",
		ShippingZipCode:      "50000",
		RecipientFirstName:   "ณัฐญา",
		RecipientLastName:    "ชุติบุตร",
		RecipientPhoneNumber: "0970809292",
		PaymentMethodID:      1,
		BurnPoint:            0,
		SubTotalPrice:        100.00,
		DiscountPrice:        0,
		TotalPrice:           100.00,
	}

	mockPointInterface := new(mockPointInterface)
	mockPointInterface.On("CheckBurnPoint", uid, 0).Return(true, nil)

	orderService := order.OrderService{
		PointService: mockPointInterface,
	}

	actual, err := orderService.CreateOrder(uid, submittedOrder)

	assert.Equal(t, expected, actual)
	assert.Equal(t, expectedErr, err)
}

func Test_CreateOrder_Input_Submitted_Order_Should_be_Return_Create_Order_Error(t *testing.T) {
	expected := order.Order{}

	uid := 1
	oid := 8004359103
	productPrice := 12.95
	yearPrefix := "26"
	lastOrderNumber := "2601129522031"
	nextSeq := 32
	fixedTime := time.Date(2026, 1, 12, 0, 0, 0, 0, time.UTC)
	orderNumber := "2601129522032"

	submittedOrder := order.SubmitedOrder{
		Cart: []order.OrderProduct{
			{
				ProductID: 2,
				Quantity:  1,
			},
		},
		ShippingMethodID:     1,
		ShippingAddress:      "405/37 ถ.มหิดล",
		ShippingSubDistrict:  "ท่าศาลา",
		ShippingDistrict:     "เมือง",
		ShippingProvince:     "เชียงใหม่",
		ShippingZipCode:      "50000",
		RecipientFirstName:   "ณัฐญา",
		RecipientLastName:    "ชุติบุตร",
		RecipientPhoneNumber: "0970809292",
		PaymentMethodID:      1,
		BurnPoint:            0,
		SubTotalPrice:        100.00,
		DiscountPrice:        0,
		TotalPrice:           100.00,
	}

	mockPointInterface := new(mockPointInterface)
	mockPointInterface.On("CheckBurnPoint", uid, 0).Return(true, nil)

	mockProductRepository := new(mockProductRepository)
	mockProductRepository.On("GetProductByID", submittedOrder.Cart[0].ProductID).Return(product.ProductDetail{
		ID:           submittedOrder.Cart[0].ProductID,
		Name:         "43 Piece dinner Set",
		Price:        productPrice,
		PriceTHB:     0,
		PriceFullTHB: 0,
		Stock:        1,
		Brand:        "Coolkidz",
		Image:        "43_Piece_Dinner_Set.jpg",
	}, nil)

	mockShippingRepository := new(mockShippingRepository)
	mockShippingRepository.On("GetShippingMethodByID", submittedOrder.ShippingMethodID).Return(shipping.ShippingMethodDetail{
		ID:          1,
		Name:        "Kerry",
		Description: "",
		Fee:         50,
	}, nil)

	mockOrderRepository := new(mockOrderRepository)
	mockOrderRepository.On("GetLastOrderNumber", yearPrefix).Return(lastOrderNumber, nil)

	mockOrderHelper := new(mockOrderHelper)
	mockOrderHelper.On("GetNextSequence", lastOrderNumber).Return(nextSeq, nil)
	mockOrderHelper.On("GenerateOrderNumber", submittedOrder.PaymentMethodID, submittedOrder.ShippingMethodID, nextSeq, fixedTime).Return(orderNumber, nil)

	orderDetail := order.OrderDetail{
		OrderNumber:      orderNumber,
		ShippingMethodID: submittedOrder.ShippingMethodID,
		PaymentMethodID:  submittedOrder.PaymentMethodID,
		SubTotalPrice:    465.811034,
		DiscountPrice:    0,
		TotalPrice:       515.8110340000001,
		ShippingFee:      50,
		BurnPoint:        0,
		EarnPoint:        4,
	}
	mockOrderRepository.On("CreateOrder", uid, orderDetail).Return(oid, errors.New("CreateOrder Error"))

	orderService := order.OrderService{
		ProductRepository:  mockProductRepository,
		OrderRepository:    mockOrderRepository,
		PointService:       mockPointInterface,
		ShippingRepository: mockShippingRepository,
		OrderHelper:        mockOrderHelper,
		Clock:              func() time.Time { return fixedTime },
	}

	actual, err := orderService.CreateOrder(uid, submittedOrder)

	assert.Equal(t, expected, actual)
	assert.NotNil(t, err)
}

func Test_CreateOrder_Input_Submitted_Order_Should_be_Return_Create_Shipping_Error(t *testing.T) {
	expected := order.Order{}

	uid := 1
	oid := 8004359103
	productPrice := 12.95
	yearPrefix := "26"
	lastOrderNumber := "2612129522079"
	nextSeq := 80
	fixedTime := time.Date(2026, 12, 12, 0, 0, 0, 0, time.UTC)
	orderNumber := "2612129522080"

	submittedOrder := order.SubmitedOrder{
		Cart: []order.OrderProduct{
			{
				ProductID: 2,
				Quantity:  1,
			},
		},
		ShippingMethodID:     1,
		ShippingAddress:      "405/37 ถ.มหิดล",
		ShippingSubDistrict:  "ท่าศาลา",
		ShippingDistrict:     "เมือง",
		ShippingProvince:     "เชียงใหม่",
		ShippingZipCode:      "50000",
		RecipientFirstName:   "ณัฐญา",
		RecipientLastName:    "ชุติบุตร",
		RecipientPhoneNumber: "0970809292",
		PaymentMethodID:      1,
		BurnPoint:            0,
		SubTotalPrice:        100.00,
		DiscountPrice:        0,
		TotalPrice:           100.00,
	}

	mockPointInterface := new(mockPointInterface)
	mockPointInterface.On("CheckBurnPoint", uid, 0).Return(true, nil)

	mockProductRepository := new(mockProductRepository)
	mockProductRepository.On("GetProductByID", submittedOrder.Cart[0].ProductID).Return(product.ProductDetail{
		ID:           submittedOrder.Cart[0].ProductID,
		Name:         "43 Piece dinner Set",
		Price:        productPrice,
		PriceTHB:     0,
		PriceFullTHB: 0,
		Stock:        1,
		Brand:        "Coolkidz",
		Image:        "43_Piece_Dinner_Set.jpg",
	}, nil)

	mockShippingRepository := new(mockShippingRepository)
	mockShippingRepository.On("GetShippingMethodByID", submittedOrder.ShippingMethodID).Return(shipping.ShippingMethodDetail{
		ID:          1,
		Name:        "Kerry",
		Description: "",
		Fee:         50,
	}, nil)

	mockOrderRepository := new(mockOrderRepository)
	mockOrderRepository.On("GetLastOrderNumber", yearPrefix).Return(lastOrderNumber, nil)

	mockOrderHelper := new(mockOrderHelper)
	mockOrderHelper.On("GetNextSequence", lastOrderNumber).Return(nextSeq, nil)
	mockOrderHelper.On("GenerateOrderNumber", submittedOrder.PaymentMethodID, submittedOrder.ShippingMethodID, nextSeq, fixedTime).Return(orderNumber, nil)

	orderDetail := order.OrderDetail{
		OrderNumber:      orderNumber,
		ShippingMethodID: submittedOrder.ShippingMethodID,
		PaymentMethodID:  submittedOrder.PaymentMethodID,
		SubTotalPrice:    465.811034,
		DiscountPrice:    0,
		TotalPrice:       515.8110340000001,
		ShippingFee:      50,
		BurnPoint:        0,
		EarnPoint:        4,
	}

	mockOrderRepository.On("CreateOrder", uid, orderDetail).Return(oid, nil)

	shippingInfo := order.ShippingInfo{
		ShippingMethodID:     1,
		ShippingAddress:      "405/37 ถ.มหิดล",
		ShippingSubDistrict:  "ท่าศาลา",
		ShippingDistrict:     "เมือง",
		ShippingProvince:     "เชียงใหม่",
		ShippingZipCode:      "50000",
		RecipientFirstName:   "ณัฐญา",
		RecipientLastName:    "ชุติบุตร",
		RecipientPhoneNumber: "0970809292",
	}
	mockOrderRepository.On("CreateShipping", uid, oid, shippingInfo).Return(1, errors.New("CreateShipping Error"))

	orderService := order.OrderService{
		ProductRepository:  mockProductRepository,
		OrderRepository:    mockOrderRepository,
		PointService:       mockPointInterface,
		ShippingRepository: mockShippingRepository,
		OrderHelper:        mockOrderHelper,
		Clock:              func() time.Time { return fixedTime },
	}

	actual, err := orderService.CreateOrder(uid, submittedOrder)

	assert.Equal(t, expected, actual)
	assert.NotNil(t, err)
}

func Test_CreateOrder_Input_Submitted_Order_Should_be_Return_Create_Order_Product_Error(t *testing.T) {
	expected := order.Order{}

	uid := 1
	oid := 8004359103
	productPrice := 12.95
	yearPrefix := "26"
	lastOrderNumber := "2605159522178"
	nextSeq := 179
	fixedTime := time.Date(2026, 05, 15, 0, 0, 0, 0, time.UTC)
	orderNumber := "2605159522179"

	submittedOrder := order.SubmitedOrder{
		Cart: []order.OrderProduct{
			{
				ProductID: 2,
				Quantity:  1,
			},
		},
		ShippingMethodID:     1,
		ShippingAddress:      "405/37 ถ.มหิดล",
		ShippingSubDistrict:  "ท่าศาลา",
		ShippingDistrict:     "เมือง",
		ShippingProvince:     "เชียงใหม่",
		ShippingZipCode:      "50000",
		RecipientFirstName:   "ณัฐญา",
		RecipientLastName:    "ชุติบุตร",
		RecipientPhoneNumber: "0970809292",
		PaymentMethodID:      1,
		BurnPoint:            0,
		SubTotalPrice:        100.00,
		DiscountPrice:        0,
		TotalPrice:           100.00,
	}

	mockPointInterface := new(mockPointInterface)
	mockPointInterface.On("CheckBurnPoint", uid, submittedOrder.BurnPoint).Return(true, nil)

	mockProductRepository := new(mockProductRepository)
	mockProductRepository.On("GetProductByID", submittedOrder.Cart[0].ProductID).Return(product.ProductDetail{
		ID:           submittedOrder.Cart[0].ProductID,
		Name:         "43 Piece dinner Set",
		Price:        productPrice,
		PriceTHB:     0,
		PriceFullTHB: 0,
		Stock:        1,
		Brand:        "Coolkidz",
		Image:        "43_Piece_Dinner_Set.jpg",
	}, nil)

	mockShippingRepository := new(mockShippingRepository)
	mockShippingRepository.On("GetShippingMethodByID", submittedOrder.ShippingMethodID).Return(shipping.ShippingMethodDetail{
		ID:          1,
		Name:        "Kerry",
		Description: "",
		Fee:         50,
	}, nil)

	mockOrderRepository := new(mockOrderRepository)
	mockOrderRepository.On("GetLastOrderNumber", yearPrefix).Return(lastOrderNumber, nil)

	mockOrderHelper := new(mockOrderHelper)
	mockOrderHelper.On("GetNextSequence", lastOrderNumber).Return(nextSeq, nil)
	mockOrderHelper.On("GenerateOrderNumber", submittedOrder.PaymentMethodID, submittedOrder.ShippingMethodID, nextSeq, fixedTime).Return(orderNumber, nil)

	orderDetail := order.OrderDetail{
		OrderNumber:      orderNumber,
		ShippingMethodID: submittedOrder.ShippingMethodID,
		PaymentMethodID:  submittedOrder.PaymentMethodID,
		SubTotalPrice:    465.811034,
		DiscountPrice:    0,
		TotalPrice:       515.8110340000001,
		ShippingFee:      50,
		BurnPoint:        0,
		EarnPoint:        4,
	}

	mockOrderRepository.On("CreateOrder", uid, orderDetail).Return(oid, nil)

	shippingInfo := order.ShippingInfo{
		ShippingMethodID:     1,
		ShippingAddress:      "405/37 ถ.มหิดล",
		ShippingSubDistrict:  "ท่าศาลา",
		ShippingDistrict:     "เมือง",
		ShippingProvince:     "เชียงใหม่",
		ShippingZipCode:      "50000",
		RecipientFirstName:   "ณัฐญา",
		RecipientLastName:    "ชุติบุตร",
		RecipientPhoneNumber: "0970809292",
	}
	mockOrderRepository.On("CreateShipping", uid, oid, shippingInfo).Return(1, nil)

	mockOrderRepository.On("CreateOrderProduct", oid, submittedOrder.Cart[0].ProductID, submittedOrder.Cart[0].Quantity, productPrice).Return(errors.New("CreateOrderProduct Error"))

	orderService := order.OrderService{
		ProductRepository:  mockProductRepository,
		OrderRepository:    mockOrderRepository,
		PointService:       mockPointInterface,
		ShippingRepository: mockShippingRepository,
		OrderHelper:        mockOrderHelper,
		Clock:              func() time.Time { return fixedTime },
	}

	actual, err := orderService.CreateOrder(uid, submittedOrder)

	assert.Equal(t, expected, actual)
	assert.NotNil(t, err)
}

func Test_OrderBurnPoint_Input_Burn_Points_100_Should_be_Return_Total_Point_50(t *testing.T) {
	expected := point.TotalPoint{
		Point: 50,
	}

	uid := 1
	burnPoint := 100
	submitedPoint := point.SubmitedPoint{
		Amount: -(burnPoint),
	}

	mockPointInterface := new(mockPointInterface)
	mockPointInterface.On("DeductPoint", uid, submitedPoint).Return(point.TotalPoint{
		Point: 50,
	}, nil)

	orderService := order.OrderService{
		PointService: mockPointInterface,
	}

	actual, err := orderService.OrderBurnPoint(uid, burnPoint)

	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)
}

func Test_OrderBurnPoint_Input_Burn_Points_100_Should_be_Return_Totol_Point_Error(t *testing.T) {
	expected := point.TotalPoint{}

	uid := 1
	burnPoint := 100
	submitedPoint := point.SubmitedPoint{
		Amount: -(burnPoint),
	}

	mockPointInterface := new(mockPointInterface)
	mockPointInterface.On("DeductPoint", uid, submitedPoint).Return(point.TotalPoint{}, errors.New("DeductPoint Error"))

	orderService := order.OrderService{
		PointService: mockPointInterface,
	}

	actual, err := orderService.OrderBurnPoint(uid, burnPoint)

	assert.Equal(t, expected, actual)
	assert.NotNil(t, err)
}

func Test_GetOrderSummary_Should_Return_One_Product_If_OrderNumber_is_2601069522001(t *testing.T) {
	userID := 4
	orderID := 1
	trackingNumber := "KR-443947172"
	orderNumber := "2601069522001"
	updatedTime := time.Date(2026, 2, 28, 18, 58, 44, 0, time.UTC)
	expectedUpdateTime := "01-03-2026 01:58:44"

	orderDetail := order.OrderDetailWithTrackingNumber{
		ID:               orderID,
		OrderNumber:      orderNumber,
		UserID:           userID,
		ShippingMethodID: 1,
		PaymentMethodID:  1,
		SubTotalPrice:    4314.6,
		DiscountPrice:    0,
		TotalPrice:       4364.6,
		ShippingFee:      50,
		BurnPoint:        0,
		EarnPoint:        43,
		TransactionID:    "TXN202512250934",
		Status:           "paid",
		TrackingNumber:   trackingNumber,
		Updated:          updatedTime,
	}

	orderProduct := []order.OrderProductWithPrice{
		{
			ProductBrand: "SportsFun",
			ProductName:  "Balance Training Bicycle",
			Quantity:     1,
			Price:        119.95,
		},
	}

	userDetail := auth.UserPayload{
		UserID:    userID,
		FirstName: "Noppadon",
		LastName:  "Sookwattana",
		Username:  "noppadon.s",
	}

	expected := order.OrderSummary{
		OrderNumber:    orderNumber,
		FirstName:      userDetail.FirstName,
		LastName:       userDetail.LastName,
		TrackingNumber: trackingNumber,
		ShippingMethod: "Kerry",
		PaymentMethod:  "Credit Card / Debit Card",
		OrderProductList: []order.OrderSummaryProduct{
			{
				ProductBrand:  "SportsFun",
				ProductName:   "Balance Training Bicycle",
				Quantity:      1,
				PriceTHB:      4314.6,
				TotalPriceTHB: 4314.6,
			},
		},
		SubTotalPrice:  orderDetail.SubTotalPrice,
		DiscountPrice:  orderDetail.DiscountPrice,
		TotalPrice:     orderDetail.TotalPrice,
		ShippingFee:    orderDetail.ShippingFee,
		BurnPoint:      orderDetail.BurnPoint,
		ReceivingPoint: orderDetail.EarnPoint,
		IssuedDate:     expectedUpdateTime,
	}

	mockOrderRepository := new(mockOrderRepository)
	mockOrderRepository.On("GetOrderWithTrackingNumberByOrderNumber", orderNumber).Return(orderDetail, nil)
	mockOrderRepository.On("GetOrderProductWithPrice", orderID).Return(orderProduct, nil)

	mockUserRepository := new(mockUserRepository)
	mockUserRepository.On("FindByID", userID).Return(userDetail, nil)

	orderService := order.OrderService{
		OrderRepository: mockOrderRepository,
		UserRepository:  mockUserRepository,
	}

	actual, err := orderService.GetOrderSummary(orderNumber)
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func Test_GetOrderSummary_Should_Return_Two_Products_If_OrderOrderNumber_is_2601069522002(t *testing.T) {
	userID := 5
	orderID := 2
	trackingNumber := "KR-304590466"
	orderNumber := "2601069522002"
	updatedTime := time.Date(2026, 2, 14, 1, 40, 32, 0, time.UTC)
	expectedUpdateTime := "14-02-2026 08:40:32"

	orderDetail := order.OrderDetailWithTrackingNumber{
		ID:               orderID,
		OrderNumber:      orderNumber,
		UserID:           userID,
		ShippingMethodID: 1,
		PaymentMethodID:  1,
		SubTotalPrice:    5246.22,
		DiscountPrice:    0,
		TotalPrice:       5256.22,
		ShippingFee:      50,
		BurnPoint:        0,
		EarnPoint:        52,
		TransactionID:    "TXN202512251028",
		Status:           "paid",
		TrackingNumber:   trackingNumber,
		Updated:          updatedTime,
	}

	orderProduct := []order.OrderProductWithPrice{
		{
			ProductBrand: "SportsFun",
			ProductName:  "Balance Training Bicycle",
			Quantity:     1,
			Price:        119.95,
		},
		{
			ProductBrand: "CoolKidz",
			ProductName:  "43 Piece dinner Set",
			Quantity:     2,
			Price:        12.95,
		},
	}

	userDetail := auth.UserPayload{
		UserID:    userID,
		FirstName: "Pimmida",
		LastName:  "Katethong",
		Username:  "pimmida.k",
	}

	expected := order.OrderSummary{
		OrderNumber:    orderNumber,
		FirstName:      userDetail.FirstName,
		LastName:       userDetail.LastName,
		TrackingNumber: trackingNumber,
		ShippingMethod: "Kerry",
		PaymentMethod:  "Credit Card / Debit Card",
		OrderProductList: []order.OrderSummaryProduct{
			{
				ProductBrand:  "SportsFun",
				ProductName:   "Balance Training Bicycle",
				Quantity:      1,
				PriceTHB:      4314.6,
				TotalPriceTHB: 4314.6,
			},
			{
				ProductBrand:  "CoolKidz",
				ProductName:   "43 Piece dinner Set",
				Quantity:      2,
				PriceTHB:      465.81,
				TotalPriceTHB: 931.62,
			},
		},
		SubTotalPrice:  orderDetail.SubTotalPrice,
		DiscountPrice:  orderDetail.DiscountPrice,
		TotalPrice:     orderDetail.TotalPrice,
		ShippingFee:    orderDetail.ShippingFee,
		BurnPoint:      orderDetail.BurnPoint,
		ReceivingPoint: orderDetail.EarnPoint,
		IssuedDate:     expectedUpdateTime,
	}

	mockOrderRepository := new(mockOrderRepository)
	mockOrderRepository.On("GetOrderWithTrackingNumberByOrderNumber", orderNumber).Return(orderDetail, nil)
	mockOrderRepository.On("GetOrderProductWithPrice", orderID).Return(orderProduct, nil)

	mockUserRepository := new(mockUserRepository)
	mockUserRepository.On("FindByID", userID).Return(userDetail, nil)

	orderService := order.OrderService{
		OrderRepository: mockOrderRepository,
		UserRepository:  mockUserRepository,
	}

	actual, err := orderService.GetOrderSummary(orderNumber)
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}
