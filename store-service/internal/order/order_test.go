package order_test

import (
	"errors"
	"fmt"
	"store-service/internal/order"
	"store-service/internal/point"
	"store-service/internal/product"
	"store-service/internal/shipping"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CreateOrder_Input_Submitted_Order_Should_be_OrderID_8004359103(t *testing.T) {
	expected := order.Order{
		OrderID: 8004359103,
	}

	uid := 1
	oid := 8004359103
	productPrice := 12.95

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

	orderDetail := order.OrderDetail{
		ShippingMethodID: submittedOrder.ShippingMethodID,
		PaymentMethodID:  submittedOrder.PaymentMethodID,
		SubTotalPrice:    465.811034,
		DiscountPrice:    0,
		TotalPrice:       515.8110340000001,
		ShippingFee:      50,
		BurnPoint:        0,
		EarnPoint:        4,
	}
	mockOrderRepository := new(mockOrderRepository)
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

	orderDetail := order.OrderDetail{
		ShippingMethodID: submittedOrder.ShippingMethodID,
		PaymentMethodID:  submittedOrder.PaymentMethodID,
		SubTotalPrice:    465.811034,
		DiscountPrice:    0,
		TotalPrice:       515.8110340000001,
		ShippingFee:      50,
		BurnPoint:        0,
		EarnPoint:        4,
	}
	mockOrderRepository := new(mockOrderRepository)
	mockOrderRepository.On("CreateOrder", uid, orderDetail).Return(oid, errors.New("CreateOrder Error"))

	orderService := order.OrderService{
		ProductRepository:  mockProductRepository,
		OrderRepository:    mockOrderRepository,
		PointService:       mockPointInterface,
		ShippingRepository: mockShippingRepository,
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

	orderDetail := order.OrderDetail{
		ShippingMethodID: submittedOrder.ShippingMethodID,
		PaymentMethodID:  submittedOrder.PaymentMethodID,
		SubTotalPrice:    465.811034,
		DiscountPrice:    0,
		TotalPrice:       515.8110340000001,
		ShippingFee:      50,
		BurnPoint:        0,
		EarnPoint:        4,
	}
	mockOrderRepository := new(mockOrderRepository)
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

	orderDetail := order.OrderDetail{
		ShippingMethodID: submittedOrder.ShippingMethodID,
		PaymentMethodID:  submittedOrder.PaymentMethodID,
		SubTotalPrice:    465.811034,
		DiscountPrice:    0,
		TotalPrice:       515.8110340000001,
		ShippingFee:      50,
		BurnPoint:        0,
		EarnPoint:        4,
	}
	mockOrderRepository := new(mockOrderRepository)
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
