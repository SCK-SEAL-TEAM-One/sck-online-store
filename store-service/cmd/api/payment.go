package api

import (
	"errors"
	"log/slog"
	"net/http"
	"store-service/internal/order"
	"store-service/internal/payment"

	"github.com/gin-gonic/gin"
)

// PaymentAPI handles payment-related operations
type PaymentAPI struct {
	PaymentService payment.PaymentInterface
}

// @Summary Confirm payment
// @Description Process and confirm a payment for an order
// @Tags payment
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param payment body payment.SubmitedPayment true "Payment details to confirm"
// @Success 200 {object} payment.SubmitedPaymentResponse "Payment confirmation details"
// @Failure 400 {string} string "Bad Request - Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500
// @Router /api/v1/payment/confirm [post]
func (api PaymentAPI) ConfirmPaymentHandler(context *gin.Context) {
	uid := context.GetInt("userID")
	var request payment.SubmitedPayment
	if err := context.BindJSON(&request); err != nil {
		context.String(http.StatusBadRequest, err.Error())
		slog.Error("bad request", "error", err)
		return
	}

	ctx := context.Request.Context()
	confirmPayment, err := api.PaymentService.ConfirmPayment(ctx, uid, request)
	if err != nil {
		if errors.Is(err, order.ErrOrderNotFound) {
			slog.ErrorContext(ctx, "PaymentService.ConfirmPayment not found", "orderNumber", request.OrderNumber, "error", err)
			context.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, confirmPayment)
}
