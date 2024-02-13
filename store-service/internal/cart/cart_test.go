package cart_test

import (
	"database/sql"
	"store-service/internal/cart"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetCart_Should_be_Have_Data_and_Receive_Point_4(t *testing.T) {
	expected := cart.CartResult{
		Carts: []cart.CartDetail{
			{
				ID:           1,
				UserID:       1,
				ProductID:    2,
				Quantity:     1,
				Name:         "43 Piece dinner Set",
				Price:        12.95,
				PriceTHB:     465.81,
				PriceFullTHB: 465.811034,
				Image:        "/43_Piece_dinner_Set.png",
				Stock:        10,
				Brand:        "CoolKidz",
			},
		},
		Summary: cart.CartSummary{
			TotalPrice:        12.95,
			TotalPriceTHB:     465.81,
			TotalPriceFullTHB: 465.811034,
			ReceivePoint:      4,
		},
	}

	uid := 1
	res := []cart.CartDetail{
		{
			ID:           1,
			UserID:       1,
			ProductID:    2,
			Quantity:     1,
			Name:         "43 Piece dinner Set",
			Price:        12.95,
			PriceTHB:     0,
			PriceFullTHB: 0,
			Image:        "/43_Piece_dinner_Set.png",
			Stock:        10,
			Brand:        "CoolKidz",
		},
	}
	mockCartRepository := new(mockCartRepository)
	mockCartRepository.On("GetCartDetail", uid).Return(res, nil)

	cartService := cart.CartService{
		CartRepository: mockCartRepository,
	}
	actual, err := cartService.GetCart(uid)

	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)
}

func Test_GetCart_Should_be_Empty(t *testing.T) {
	expected := cart.CartResult{
		Carts:   []cart.CartDetail{},
		Summary: cart.CartSummary{},
	}
	uid := 1
	res := []cart.CartDetail{}
	mockCartRepository := new(mockCartRepository)
	mockCartRepository.On("GetCartDetail", uid).Return(res, nil)

	cartService := cart.CartService{
		CartRepository: mockCartRepository,
	}
	actual, err := cartService.GetCart(uid)

	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)
}

func Test_AddCart_Input_Submitted_First_Product_Should_be_Have_1_Quantity_and_Receive_Point_43(t *testing.T) {
	expected := cart.CartResult{
		Carts: []cart.CartDetail{
			{
				ID:           1,
				UserID:       1,
				ProductID:    1,
				Quantity:     1,
				Name:         "Balance Training Bicycle",
				Price:        119.95,
				PriceTHB:     4314.6,
				PriceFullTHB: 4314.597182,
				Image:        "/Balance_Training_Bicycle.png",
				Stock:        100,
				Brand:        "SportsFun",
			},
		},
		Summary: cart.CartSummary{
			TotalPrice:        119.95,
			TotalPriceTHB:     4314.6,
			TotalPriceFullTHB: 4314.597182,
			ReceivePoint:      43,
		},
	}
	submitedCart := cart.SubmitedCart{
		ProductID: 1,
		Quantity:  1,
	}
	uid := 1
	mockCartRepository := new(mockCartRepository)
	mockCartRepository.On("GetCartByProductID", uid, submitedCart.ProductID).Return(cart.Cart{}, sql.ErrNoRows)
	mockCartRepository.On("CreateCart", uid, submitedCart.ProductID, submitedCart.Quantity).Return(1, nil)

	res := []cart.CartDetail{
		{
			ID:           1,
			UserID:       1,
			ProductID:    1,
			Quantity:     1,
			Name:         "Balance Training Bicycle",
			Price:        119.95,
			PriceTHB:     0,
			PriceFullTHB: 0,
			Image:        "/Balance_Training_Bicycle.png",
			Stock:        100,
			Brand:        "SportsFun",
		},
	}
	mockCartRepository.On("GetCartDetail", uid).Return(res, nil)

	cartService := cart.CartService{
		CartRepository: mockCartRepository,
	}
	actual, err := cartService.AddCart(uid, submitedCart)

	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)
}

func Test_AddCart_Input_Submitted_More_Product_Should_be_Have_2_Quantity_and_Receive_Point_86(t *testing.T) {
	expected := cart.CartResult{
		Carts: []cart.CartDetail{
			{
				ID:           1,
				UserID:       1,
				ProductID:    1,
				Quantity:     2,
				Name:         "Balance Training Bicycle",
				Price:        119.95,
				PriceTHB:     4314.6,
				PriceFullTHB: 4314.597182,
				Image:        "/Balance_Training_Bicycle.png",
				Stock:        100,
				Brand:        "SportsFun",
			},
		},
		Summary: cart.CartSummary{
			TotalPrice:        239.9,
			TotalPriceTHB:     8629.19,
			TotalPriceFullTHB: 8629.194364,
			ReceivePoint:      86,
		},
	}
	submitedCart := cart.SubmitedCart{
		ProductID: 1,
		Quantity:  1,
	}
	uid := 1
	mockCartRepository := new(mockCartRepository)
	mockCartRepository.On("GetCartByProductID", uid, submitedCart.ProductID).Return(cart.Cart{
		ID:        1,
		UserID:    1,
		ProductID: 1,
		Quantity:  1,
	}, nil)
	mockCartRepository.On("UpdateCart", uid, submitedCart.ProductID, 2).Return(nil)

	res := []cart.CartDetail{
		{
			ID:           1,
			UserID:       1,
			ProductID:    1,
			Quantity:     2,
			Name:         "Balance Training Bicycle",
			Price:        119.95,
			PriceTHB:     0,
			PriceFullTHB: 0,
			Image:        "/Balance_Training_Bicycle.png",
			Stock:        100,
			Brand:        "SportsFun",
		},
	}
	mockCartRepository.On("GetCartDetail", uid).Return(res, nil)

	cartService := cart.CartService{
		CartRepository: mockCartRepository,
	}
	actual, err := cartService.AddCart(uid, submitedCart)

	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)
}

func Test_UpdateCart_Input_Submitted_Quantity_2_Should_be_Have_2_Quantity_and_Receive_Point_9(t *testing.T) {
	expected := cart.CartResult{
		Carts: []cart.CartDetail{
			{
				ID:           1,
				UserID:       1,
				ProductID:    2,
				Quantity:     2,
				Name:         "43 Piece dinner Set",
				Price:        12.95,
				PriceTHB:     465.81,
				PriceFullTHB: 465.811034,
				Image:        "/43_Piece_dinner_Set.png",
				Stock:        200,
				Brand:        "CoolKidz",
			},
		},
		Summary: cart.CartSummary{
			TotalPrice:        25.9,
			TotalPriceTHB:     931.62,
			TotalPriceFullTHB: 931.622068,
			ReceivePoint:      9,
		},
	}
	submitedCart := cart.SubmitedCart{
		ProductID: 2,
		Quantity:  2,
	}
	uid := 1
	mockCartRepository := new(mockCartRepository)
	mockCartRepository.On("UpdateCart", uid, submitedCart.ProductID, submitedCart.Quantity).Return(nil)

	res := []cart.CartDetail{
		{
			ID:           1,
			UserID:       1,
			ProductID:    2,
			Quantity:     2,
			Name:         "43 Piece dinner Set",
			Price:        12.95,
			PriceTHB:     0,
			PriceFullTHB: 0,
			Image:        "/43_Piece_dinner_Set.png",
			Stock:        200,
			Brand:        "CoolKidz",
		},
	}
	mockCartRepository.On("GetCartDetail", uid).Return(res, nil)

	cartService := cart.CartService{
		CartRepository: mockCartRepository,
	}
	actual, err := cartService.UpdateCart(uid, submitedCart)

	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)
}

func Test_UpdateCart_Input_Submitted_Quantity_0_Should_be_Have_0_Quantity_and_Receive_Point_0(t *testing.T) {
	expected := cart.CartResult{
		Carts:   []cart.CartDetail{},
		Summary: cart.CartSummary{},
	}
	submitedCart := cart.SubmitedCart{
		ProductID: 1,
		Quantity:  0,
	}
	uid := 1
	mockCartRepository := new(mockCartRepository)
	mockCartRepository.On("DeleteCart", uid, submitedCart.ProductID).Return(nil)

	res := []cart.CartDetail{}
	mockCartRepository.On("GetCartDetail", uid).Return(res, nil)

	cartService := cart.CartService{
		CartRepository: mockCartRepository,
	}
	actual, err := cartService.UpdateCart(uid, submitedCart)

	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)
}
