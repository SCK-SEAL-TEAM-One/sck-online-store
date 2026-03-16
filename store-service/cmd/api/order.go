package api

import (
	"errors"
	"fmt"
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
// @Description Order confirmation response containing order number
type OrderConfirmation struct {
	OrderNumber string `json:"order_number"`
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
	uid := context.GetInt("userID")
	var request order.SubmitedOrder
	if err := context.BindJSON(&request); err != nil {
		context.String(http.StatusBadRequest, err.Error())
		log.Printf("bad request %s", err.Error())
		return
	}

	ctx := context.Request.Context()
	createdOrder, err := api.OrderService.CreateOrder(ctx, uid, request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, OrderConfirmation{
		OrderNumber: createdOrder.OrderNumber,
	})
}

func (api OrderAPI) GetOrderSummaryHandler(context *gin.Context) {
	acceptHeader := context.GetHeader("Accept")
	allowedHeaders := []string{"application/pdf", "", "*/*", "application/json"}
	isAllowed := false
	for _, header := range allowedHeaders {
		if header == acceptHeader {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		// 406 Not Acceptable
		context.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	ctx := context.Request.Context()
	orderNumber := context.Param("id")

	orderSummary, err := api.OrderService.GetOrderSummary(ctx, orderNumber)
	if err != nil {
		if errors.Is(err, order.ErrOrderNotFound) {
			log.Printf("OrderService.GetOrderSummary not found Order Number: %s %s", orderNumber, err.Error())
			context.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		log.Printf("OrderService.GetOrderSummary internal error %s", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if acceptHeader == "application/json" {
		context.JSON(http.StatusOK, orderSummary)
		return
	}

	pdfData, err := api.OrderService.GeneratePDFFromData(orderSummary)
	if err != nil {
		log.Printf("OrderService.GeneratePDFFromData internal error %s", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	filename := fmt.Sprintf("Order_Summary_%s.pdf", orderNumber)
	contentDisposition := fmt.Sprintf("attachment; filename=\"%s\"", filename)

	context.Header("Access-Control-Expose-Headers", "Content-Disposition")
	context.Header(
		"Content-Disposition",
		contentDisposition,
	)

	context.Data(http.StatusOK, "application/pdf", pdfData)
}
