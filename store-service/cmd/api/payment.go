package api

import (
	"log"
	"net/http"
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
	uid := 1
	var request payment.SubmitedPayment
	if err := context.BindJSON(&request); err != nil {
		context.String(http.StatusBadRequest, err.Error())
		log.Printf("bad request %s", err.Error())
		return
	}

	confirmPayment, err := api.PaymentService.ConfirmPayment(uid, request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, confirmPayment)
}
