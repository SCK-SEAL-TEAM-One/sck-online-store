//go:build integration
// +build integration

package payment_test

import (
	"store-service/internal/payment"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Payment_Input_PaymentDetail_CardNumber_4719700591590995_Should_Be_TransactionID_Not_Empty(t *testing.T) {

	paymentDetail := payment.PaymentDetail{
		CardNumber:   "4719700591590995",
		CVV:          752,
		ExpiredMonth: 12,
		ExpiredYear:  27,
		CardName:     "SCK ShuHaRi",
		Amount:       90.00,
		Currency:     "USD",
		MerchantID:   1,
	}

	gateway := payment.BankGateway{
		BankEndpoint: "http://localhost:8882",
	}
	actual, err := gateway.Payment(paymentDetail)

	assert.NotEmpty(t, actual)
	assert.Equal(t, nil, err)
}

func Test_GetCardDetail_Input_PaymentDetail_CardNumber_4719700591590995_Should_Be_TransactionID_Not_Empty(t *testing.T) {
	expected := payment.CardDetail{
		CardNumber:   "4719700591590995",
		CVV:          752,
		ExpiredMonth: 12,
		ExpiredYear:  27,
		CardName:     "SCK ShuHaRi",
	}

	uid := 1
	orgID := 1

	gateway := payment.BankGateway{
		BankEndpoint: "http://localhost:8882",
	}
	actual, err := gateway.GetCardDetail(orgID, uid)

	assert.Equal(t, expected, actual)
	assert.Equal(t, nil, err)
}
