package api

import (
	"log/slog"
	"net/http"

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
	var cartDetails cart.CartResult
	var err error

	ctx := context.Request.Context()
	uid := context.GetInt("userID")
	cartDetails, err = api.CartService.GetCart(ctx, uid)

	if err != nil {
		slog.ErrorContext(ctx, "CartService.GetCart failed",
			"log_type", "error",
			"error_code", "CART_QUERY_FAILED",
			"error_message", err.Error(),
			"user_id", uid,
		)
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
	ctx := context.Request.Context()
	uid := context.GetInt("userID")

	var request cart.SubmitedCart
	var addedCart cart.CartResult
	var err error

	if err := context.BindJSON(&request); err != nil {
		slog.ErrorContext(ctx, "Add cart bad request",
			"log_type", "error",
			"error_code", "INVALID_REQUEST",
			"error_message", err.Error(),
			"user_id", uid,
		)
		context.String(http.StatusBadRequest, err.Error())
		return
	}

	addedCart, err = api.CartService.AddCart(ctx, uid, request)

	if err != nil {
		slog.ErrorContext(ctx, "CartService.AddCart failed",
			"log_type", "error",
			"error_code", "CART_ADD_FAILED",
			"error_message", err.Error(),
			"user_id", uid,
			slog.Any("request", map[string]any{"product_id": request.ProductID, "quantity": request.Quantity}),
		)
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	slog.InfoContext(ctx, "Cart item added",
		"log_type", "state_change",
		"entity_type", "cart",
		"entity_id", request.ProductID,
		"changed_by", uid,
		slog.Any("after", map[string]any{
			"product_id": request.ProductID,
			"quantity":   request.Quantity,
		}),
		slog.Any("changed_fields", []string{"product_id", "quantity"}),
	)

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
	ctx := context.Request.Context()
	uid := context.GetInt("userID")

	var request cart.SubmitedCart
	if err := context.BindJSON(&request); err != nil {
		slog.ErrorContext(ctx, "Update cart bad request",
			"log_type", "error",
			"error_code", "INVALID_REQUEST",
			"error_message", err.Error(),
			"user_id", uid,
		)
		context.String(http.StatusBadRequest, err.Error())
		return
	}

	updatedCart, err := api.CartService.UpdateCart(ctx, uid, request)
	if err != nil {
		slog.ErrorContext(ctx, "CartService.UpdateCart failed",
			"log_type", "error",
			"error_code", "CART_UPDATE_FAILED",
			"error_message", err.Error(),
			"user_id", uid,
			slog.Any("request", map[string]any{"product_id": request.ProductID, "quantity": request.Quantity}),
		)
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	action := "cart_item_updated"
	if request.Quantity <= 0 {
		action = "cart_item_removed"
	}
	slog.InfoContext(ctx, "Cart item changed",
		"log_type", "state_change",
		"entity_type", "cart",
		"entity_id", request.ProductID,
		"changed_by", uid,
		slog.Any("after", map[string]any{
			"product_id": request.ProductID,
			"quantity":   request.Quantity,
			"action":     action,
		}),
		slog.Any("changed_fields", []string{"quantity"}),
	)

	context.JSON(http.StatusOK, updatedCart)
}
