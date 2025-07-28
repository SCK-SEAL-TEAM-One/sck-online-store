package api

import (
	"log"
	"net/http"

	"store-service/internal/order"

	"github.com/gin-gonic/gin"
)

// @title Store Service API
// @version 1.0
// @description Store service API documentation
// @host localhost:8000
// @BasePath /api/v1
type OrderAPI struct {
	OrderService order.OrderInterface
}

// OrderConfirmation represents the response after order submission
// @Description Order confirmation response containing order ID
type OrderConfirmation struct {
	OrderID int `json:"order_id"`
}

// @Summary Submit new order
// @Description Creates a new order from the submitted order details
// @Tags order
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param order body order.SubmitedOrder true "Order details"
// @Success 200 {object} OrderConfirmation "Successfully created order"
// @Failure 400 {string} string "Bad Request - Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500
// @Router /api/v1/order [post]
func (api OrderAPI) SubmitOrderHandler(context *gin.Context) {
	uid := 1
	var request order.SubmitedOrder
	if err := context.BindJSON(&request); err != nil {
		context.String(http.StatusBadRequest, err.Error())
		log.Printf("bad request %s", err.Error())
		return
	}

	createdOrder, err := api.OrderService.CreateOrder(uid, request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, OrderConfirmation{
		OrderID: createdOrder.OrderID,
	})
}
