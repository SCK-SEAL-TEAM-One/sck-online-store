package order_test

import (
	"errors"
	"fmt"
	"store-service/internal/order"
	"store-service/internal/point"
	"store-service/internal/product"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CreateOrder_Input_Submitted_Order_Should_be_OrderID_8004359103(t *testing.T) {
	expected := order.Order{
		OrderID: 8004359103,
	}

	uid := 1
	oid := 8004359103
	pid := 2
	qty := 1
	productPrice := 12.95

	mockPointServiceInterface := new(mockPointServiceInterface)
	mockPointServiceInterface.On("CheckBurnPoint", uid, 0).Return(true, nil)

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
		DiscountPrice:        10.00,
		TotalPrice:           90.00,
	}
	mockOrderRepository := new(mockOrderRepository)
	mockOrderRepository.On("CreateOrder", uid, submittedOrder).Return(oid, nil)

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

	mockProductRepository := new(mockProductRepository)
	mockProductRepository.On("GetProductByID", pid).Return(product.ProductDetail{
		ID:    pid,
		Name:  "43 Piece dinner Set",
		Price: productPrice,
		Stock: 1,
		Brand: "Coolkidz",
		Image: "43_Piece_Dinner_Set.jpg",
	}, nil)

	mockOrderRepository.On("CreateOrderProduct", oid, pid, qty, productPrice).Return(nil)

	mockCartRepository := new(mockCartRepository)
	mockCartRepository.On("DeleteCart", uid, pid).Return(nil)

	orderService := order.OrderService{
		ProductRepository: mockProductRepository,
		OrderRepository:   mockOrderRepository,
		CartRepository:    mockCartRepository,
		PointService:      mockPointServiceInterface,
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

	mockPointServiceInterface := new(mockPointServiceInterface)
	mockPointServiceInterface.On("CheckBurnPoint", uid, -(burnPoint)).Return(false, fmt.Errorf("points are not enough, please try again"))

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
		DiscountPrice:        10.00,
		TotalPrice:           90.00,
	}

	orderService := order.OrderService{
		PointService: mockPointServiceInterface,
	}

	actual, err := orderService.CreateOrder(uid, submittedOrder)

	assert.Equal(t, expected, actual)
	assert.Equal(t, expectedErr, err)
}

func Test_CreateOrder_Input_Submitted_Order_Should_be_Return_Create_Order_Error(t *testing.T) {
	expected := order.Order{}

	uid := 1
	oid := 8004359103

	mockPointServiceInterface := new(mockPointServiceInterface)
	mockPointServiceInterface.On("CheckBurnPoint", uid, 0).Return(true, nil)

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
		DiscountPrice:        10.00,
		TotalPrice:           90.00,
	}
	mockOrderRepository := new(mockOrderRepository)
	mockOrderRepository.On("CreateOrder", uid, submittedOrder).Return(oid, errors.New("CreateOrder Error"))

	orderService := order.OrderService{
		OrderRepository: mockOrderRepository,
		PointService:    mockPointServiceInterface,
	}

	actual, err := orderService.CreateOrder(uid, submittedOrder)

	assert.Equal(t, expected, actual)
	assert.NotNil(t, err)
}

func Test_CreateOrder_Input_Submitted_Order_Should_be_Return_Create_Shipping_Error(t *testing.T) {
	expected := order.Order{}

	uid := 1
	oid := 8004359103

	mockPointServiceInterface := new(mockPointServiceInterface)
	mockPointServiceInterface.On("CheckBurnPoint", uid, 0).Return(true, nil)

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
		DiscountPrice:        10.00,
		TotalPrice:           90.00,
	}
	mockOrderRepository := new(mockOrderRepository)
	mockOrderRepository.On("CreateOrder", uid, submittedOrder).Return(oid, nil)

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
		OrderRepository: mockOrderRepository,
		PointService:    mockPointServiceInterface,
	}

	actual, err := orderService.CreateOrder(uid, submittedOrder)

	assert.Equal(t, expected, actual)
	assert.NotNil(t, err)
}

func Test_CreateOrder_Input_Submitted_Order_Should_be_Return_Create_Order_Product_Error(t *testing.T) {
	expected := order.Order{}

	uid := 1
	oid := 8004359103
	pid := 2
	qty := 1
	productPrice := 12.95

	mockPointServiceInterface := new(mockPointServiceInterface)
	mockPointServiceInterface.On("CheckBurnPoint", uid, 0).Return(true, nil)

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
		DiscountPrice:        10.00,
		TotalPrice:           90.00,
	}
	mockOrderRepository := new(mockOrderRepository)
	mockOrderRepository.On("CreateOrder", uid, submittedOrder).Return(oid, nil)

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

	mockProductRepository := new(mockProductRepository)
	mockProductRepository.On("GetProductByID", pid).Return(product.ProductDetail{
		ID:    pid,
		Name:  "43 Piece dinner Set",
		Price: productPrice,
		Stock: 1,
		Brand: "Coolkidz",
		Image: "43_Piece_Dinner_Set.jpg",
	}, nil)

	mockOrderRepository.On("CreateOrderProduct", oid, pid, qty, productPrice).Return(errors.New("CreateOrderProduct Error"))

	orderService := order.OrderService{
		ProductRepository: mockProductRepository,
		OrderRepository:   mockOrderRepository,
		PointService:      mockPointServiceInterface,
	}

	actual, err := orderService.CreateOrder(uid, submittedOrder)

	assert.Equal(t, expected, actual)
	assert.NotNil(t, err)
}

func Test_OrderBurnPoint_Input_Burn_Points_100_Should_be_Return_Totol_Point_50(t *testing.T) {
	expected := point.TotalPoint{
		Point: 50,
	}

	uid := 1
	burnPoint := 100
	submitedPoint := point.SubmitedPoint{
		Amount: -(burnPoint),
	}

	mockPointServiceInterface := new(mockPointServiceInterface)
	mockPointServiceInterface.On("DeductPoint", uid, submitedPoint).Return(point.TotalPoint{
		Point: 50,
	}, nil)

	orderService := order.OrderService{
		PointService: mockPointServiceInterface,
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

	mockPointServiceInterface := new(mockPointServiceInterface)
	mockPointServiceInterface.On("DeductPoint", uid, submitedPoint).Return(point.TotalPoint{}, errors.New("DeductPoint Error"))

	orderService := order.OrderService{
		PointService: mockPointServiceInterface,
	}

	actual, err := orderService.OrderBurnPoint(uid, burnPoint)

	assert.Equal(t, expected, actual)
	assert.NotNil(t, err)
}
