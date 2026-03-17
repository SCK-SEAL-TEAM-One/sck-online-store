package cart

import (
	"context"
	"database/sql"
	"log/slog"
	"store-service/internal/common"
)

type CartInterface interface {
	GetCart(ctx context.Context, uid int) (CartResult, error)
	AddCart(ctx context.Context, uid int, submitedCart SubmitedCart) (CartResult, error)
	UpdateCart(ctx context.Context, uid int, submitedCart SubmitedCart) (CartResult, error)
}

type CartService struct {
	CartRepository CartRepository
}

func (cartService CartService) GetCart(ctx context.Context, uid int) (CartResult, error) {
	carts, err := cartService.CartRepository.GetCartDetail(ctx, uid)
	if err != nil {
		slog.ErrorContext(ctx, "CartRepository.GetCartDetail failed",
			"log_type", "error", "error_code", "CART_QUERY_FAILED", "error_message", err.Error(), "user_id", uid)
	}

	totalPrice := 0.0
	for i := range carts {
		c := &carts[i]
		digit := common.ConvertToThb(c.Price)
		if c.ProductID == 8 {
			digit.ShortDecimal += 0.01
			digit.LongDecimal += 0.01
		}

		c.PriceTHB = digit.ShortDecimal
		c.PriceFullTHB = digit.LongDecimal
		totalPrice = totalPrice + (c.Price * float64(c.Quantity))
	}

	decimal := common.ConvertToThb(totalPrice)
	totalPriceTHB := decimal.ShortDecimal
	totalPriceFullTHB := decimal.LongDecimal

	if len(carts) == 0 {
		return CartResult{
			Carts:   []CartDetail{},
			Summary: CartSummary{},
		}, err
	}
	return CartResult{
		Carts: carts,
		Summary: CartSummary{
			TotalPrice:        totalPrice,
			TotalPriceTHB:     totalPriceTHB,
			TotalPriceFullTHB: totalPriceFullTHB,
			ReceivePoint:      common.CalculatePoint(totalPriceTHB),
		},
	}, err
}

func (cartService CartService) AddCart(ctx context.Context, uid int, submitedCart SubmitedCart) (CartResult, error) {
	cart, err := cartService.CartRepository.GetCartByProductID(ctx, uid, submitedCart.ProductID)

	if err == sql.ErrNoRows {
		_, err := cartService.CartRepository.CreateCart(ctx, uid, submitedCart.ProductID, submitedCart.Quantity)
		if err != nil {
			slog.ErrorContext(ctx, "CartRepository.CreateCart failed",
				"log_type", "error", "error_code", "CART_INSERT_FAILED", "error_message", err.Error(),
				"user_id", uid, "product_id", submitedCart.ProductID)
			return CartResult{Carts: []CartDetail{}, Summary: CartSummary{}}, err
		}
		cartResult, err := cartService.GetCart(ctx, uid)
		if err != nil {
			return CartResult{Carts: []CartDetail{}, Summary: CartSummary{}}, err
		}
		return cartResult, nil
	}

	err = cartService.CartRepository.UpdateCart(ctx, uid, submitedCart.ProductID, submitedCart.Quantity+cart.Quantity)
	if err != nil {
		slog.ErrorContext(ctx, "CartRepository.UpdateCart failed",
			"log_type", "error", "error_code", "CART_UPDATE_FAILED", "error_message", err.Error(),
			"user_id", uid, "product_id", submitedCart.ProductID)
		return CartResult{Carts: []CartDetail{}, Summary: CartSummary{}}, err
	}

	cartResult, err := cartService.GetCart(ctx, uid)
	if err != nil {
		return CartResult{Carts: []CartDetail{}, Summary: CartSummary{}}, err
	}
	return cartResult, nil
}

func (cartService CartService) UpdateCart(ctx context.Context, uid int, submitedCart SubmitedCart) (CartResult, error) {
	if submitedCart.Quantity <= 0 {
		err := cartService.CartRepository.DeleteCart(ctx, uid, submitedCart.ProductID)
		if err != nil {
			slog.ErrorContext(ctx, "CartRepository.DeleteCart failed",
				"log_type", "error", "error_code", "CART_DELETE_FAILED", "error_message", err.Error(),
				"user_id", uid, "product_id", submitedCart.ProductID)
			return CartResult{
				Carts:   []CartDetail{},
				Summary: CartSummary{},
			}, err
		}
	} else {
		err := cartService.CartRepository.UpdateCart(ctx, uid, submitedCart.ProductID, submitedCart.Quantity)
		if err != nil {
			slog.ErrorContext(ctx, "CartRepository.UpdateCart failed",
				"log_type", "error", "error_code", "CART_UPDATE_FAILED", "error_message", err.Error(),
				"user_id", uid, "product_id", submitedCart.ProductID)
			return CartResult{
				Carts:   []CartDetail{},
				Summary: CartSummary{},
			}, err
		}
	}
	cartResult, err := cartService.GetCart(ctx, uid)
	if err != nil {
		return CartResult{Carts: []CartDetail{}, Summary: CartSummary{}}, err
	}
	return cartResult, nil
}
