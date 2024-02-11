package api

import (
	"log"
	"net/http"

	"store-service/internal/order"

	"github.com/gin-gonic/gin"
)

type OrderAPI struct {
	OrderService order.OrderInterface
}

type OrderConfirmation struct {
	OrderID int `json:"order_id"`
}

func (api OrderAPI) SubmitOrderHandler(context *gin.Context) {
	uid := 1
	var request order.SubmitedOrder
	if err := context.BindJSON(&request); err != nil {
		context.String(http.StatusBadRequest, err.Error())
		log.Printf("bad request %s", err.Error())
		return
	}

	order, err := api.OrderService.CreateOrder(uid, request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, OrderConfirmation{
		OrderID: order.OrderID,
	})
}
