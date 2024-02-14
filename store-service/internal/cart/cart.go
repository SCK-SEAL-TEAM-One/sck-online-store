package cart

import (
	"database/sql"
	"log"
	"store-service/internal/common"
)

type CartInterface interface {
	GetCart(uid int) (CartResult, error)
	AddCart(uid int, submitedCart SubmitedCart) (CartResult, error)
	UpdateCart(uid int, submitedCart SubmitedCart) (CartResult, error)
}

type CartService struct {
	CartRepository CartRepository
}

func (cartService CartService) GetCart(uid int) (CartResult, error) {
	carts, err := cartService.CartRepository.GetCartDetail(uid)
	if err != nil {
		log.Printf("CartRepository.GetCartDetail internal error %s", err.Error())
	}

	totalPrice := 0.0
	for i := range carts {
		c := &carts[i]
		digit := common.ConvertToThb(c.Price)

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

func (cartService CartService) AddCart(uid int, submitedCart SubmitedCart) (CartResult, error) {
	cart, err := cartService.CartRepository.GetCartByProductID(uid, submitedCart.ProductID)

	if err == sql.ErrNoRows {
		cartService.CartRepository.CreateCart(uid, submitedCart.ProductID, submitedCart.Quantity)
		cartResult, _ := cartService.GetCart(uid)
		return cartResult, nil
	}
	err = cartService.CartRepository.UpdateCart(uid, submitedCart.ProductID, submitedCart.Quantity+cart.Quantity)
	if err != nil {
		log.Printf("CartRepository.UpdateCart internal error %s", err.Error())
		return CartResult{
			Carts:   []CartDetail{},
			Summary: CartSummary{},
		}, err
	}

	cartResult, _ := cartService.GetCart(uid)
	return cartResult, nil
}

func (cartService CartService) UpdateCart(uid int, submitedCart SubmitedCart) (CartResult, error) {
	if submitedCart.Quantity <= 0 {
		err := cartService.CartRepository.DeleteCart(uid, submitedCart.ProductID)
		if err != nil {
			log.Printf("CartRepository.DeleteCart internal error %s", err.Error())
			return CartResult{
				Carts:   []CartDetail{},
				Summary: CartSummary{},
			}, err
		}
	} else {
		err := cartService.CartRepository.UpdateCart(uid, submitedCart.ProductID, submitedCart.Quantity)
		if err != nil {
			log.Printf("CartRepository.UpdateCart internal error %s", err.Error())
			return CartResult{
				Carts:   []CartDetail{},
				Summary: CartSummary{},
			}, err
		}
	}
	cartResult, _ := cartService.GetCart(uid)
	return cartResult, nil

}
