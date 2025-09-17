package api

import (
	"log"
	"net/http"
	"strconv"

	"store-service/internal/cart"

	"github.com/gin-gonic/gin"
)

type CartAPI struct {
	CartService cart.CartInterface
}

// @Summary Get cart by user ID
// @Description Retrieves the shopping cart for a specific user
// @Tags cart
// @Accept json
// @Produce json
// @Success 200 {object} cart.CartResult
// @Failure 500
// @Router /api/v1/cart [get]
func (api CartAPI) GetCartHandler(context *gin.Context) {
	uidString := context.GetHeader("uid")
	var cartDetails cart.CartResult
	var err error

	if uidString == "" {
		cartDetails, err = api.CartService.InitCart()
	} else {
		uid, uidErr := strconv.Atoi(uidString)
		if uidErr != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid UID header",
			})
			return
		}
		cartDetails, err = api.CartService.GetCart(uid)
	}

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"uidErr": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, cartDetails)
}

// @Summary Add items to cart
// @Description Add new items to user's shopping cart
// @Tags cart
// @Accept json
// @Produce json
// @Param request body cart.SubmitedCart true "Cart items to add"
// @Success 200 {object} cart.CartResult
// @Failure 400 {string} string "Bad request error"
// @Failure 500
// @Router /api/v1/cart [post]
func (api CartAPI) AddCartHandler(context *gin.Context) {
	uidString := context.GetHeader("uid")

	var request cart.SubmitedCart
	var addedCart cart.CartResult
	var err error

	if err := context.BindJSON(&request); err != nil {
		context.String(http.StatusBadRequest, err.Error())
		log.Printf("bad request %s", err.Error())
		return
	}

	if uidString == "" {
		addedCart, err = api.CartService.AssignAndAddCart(request)

	} else {
		uid, uidErr := strconv.Atoi(context.GetHeader("uid"))
		if uidErr != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid UID header",
			})
			return
		}

		addedCart, err = api.CartService.AddCart(uid, request)
	}

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, addedCart)
}

// @Summary Update cart
// @Description Update items in user's shopping cart
// @Tags cart
// @Accept json
// @Produce json
// @Param request body cart.SubmitedCart true "Updated cart items"
// @Success 200 {object} cart.CartResult
// @Failure 400 {string} string "Bad request error"
// @Failure 500
// @Router /api/v1/cart [put]
func (api CartAPI) UpdateCartHandler(context *gin.Context) {
	var request cart.SubmitedCart
	if err := context.BindJSON(&request); err != nil {
		context.String(http.StatusBadRequest, err.Error())
		log.Printf("bad request %s", err.Error())
		return
	}

	uid, uidErr := strconv.Atoi(context.GetHeader("uid"))
	if uidErr != nil {
		uid = 1
	}

	updatedCart, err := api.CartService.UpdateCart(uid, request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, updatedCart)
}
