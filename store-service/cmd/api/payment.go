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
	ctx := context.Request.Context()

	var request payment.SubmitedPayment
	if err := context.BindJSON(&request); err != nil {
		slog.ErrorContext(ctx, "Payment confirm bad request",
			"log_type", "error",
			"error_code", "INVALID_REQUEST",
			"error_message", err.Error(),
			"user_id", uid,
		)
		context.String(http.StatusBadRequest, err.Error())
		return
	}

	confirmPayment, err := api.PaymentService.ConfirmPayment(ctx, uid, request)
	if err != nil {
		if errors.Is(err, order.ErrOrderNotFound) {
			slog.ErrorContext(ctx, "PaymentService.ConfirmPayment order not found",
				"log_type", "error",
				"error_code", "ORDER_NOT_FOUND",
				"error_message", err.Error(),
				"user_id", uid,
				slog.Any("request", map[string]any{"order_number": request.OrderNumber}),
			)
			context.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		slog.ErrorContext(ctx, "PaymentService.ConfirmPayment failed",
			"log_type", "error",
			"error_code", "PAYMENT_FAILED",
			"error_message", err.Error(),
			"user_id", uid,
			slog.Any("request", map[string]any{"order_number": request.OrderNumber}),
		)
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	slog.InfoContext(ctx, "Payment confirmed",
		"log_type", "business",
		"event", "payment_confirmed",
		"entity_type", "payment",
		"entity_id", confirmPayment.OrderNumber,
		"actor_id", uid,
		slog.Any("metadata", map[string]any{
			"tracking_number":   confirmPayment.TrackingNumber,
			"shipping_method_id": confirmPayment.ShippingMethodID,
			"payment_date":      confirmPayment.PaymentDate,
		}),
	)

	context.JSON(http.StatusOK, confirmPayment)
}
