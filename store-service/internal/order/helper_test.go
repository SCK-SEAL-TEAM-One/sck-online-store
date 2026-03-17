package order_test

import (
	"store-service/internal/order"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateOrderNumber_Should_Return_16Digit_If_Order_with_Credit_Card_And_Kerry_UserID_1(t *testing.T) {
	// Format: YYMMDD PP SM UUU SEQ = 260106 95 22 001 001
	var expected int64 = 2601069522001001

	paymentMethodID := 1
	shippingMethodID := 1
	userID := 1
	seq := 1
	fixedTime := time.Date(2026, 1, 6, 0, 0, 0, 0, time.UTC)
	helper := order.OrderHelper{}

	actual, err := helper.GenerateOrderNumber(paymentMethodID, shippingMethodID, userID, seq, fixedTime)
	assert.Equal(t, expected, actual)
	assert.Equal(t, err, nil)
}

func Test_GenerateOrderNumber_Should_Return_16Digit_If_Order_with_Credit_Card_And_Thai_Post_UserID_2(t *testing.T) {
	// Format: YYMMDD PP SM UUU SEQ = 260106 95 33 002 002
	var expected int64 = 2601069533002002

	paymentMethodID := 1
	shippingMethodID := 2
	userID := 2
	seq := 2
	fixedTime := time.Date(2026, 1, 6, 0, 0, 0, 0, time.UTC)
	helper := order.OrderHelper{}

	actual, err := helper.GenerateOrderNumber(paymentMethodID, shippingMethodID, userID, seq, fixedTime)
	assert.Equal(t, expected, actual)
	assert.Equal(t, err, nil)
}

func Test_GenerateOrderNumber_Should_Return_16Digit_If_Order_with_Credit_Card_And_Lineman_UserID_3(t *testing.T) {
	// Format: YYMMDD PP SM UUU SEQ = 260106 95 44 003 003
	var expected int64 = 2601069544003003

	paymentMethodID := 1
	shippingMethodID := 3
	userID := 3
	seq := 3
	fixedTime := time.Date(2026, 1, 6, 0, 0, 0, 0, time.UTC)
	helper := order.OrderHelper{}

	actual, err := helper.GenerateOrderNumber(paymentMethodID, shippingMethodID, userID, seq, fixedTime)
	assert.Equal(t, expected, actual)
	assert.Equal(t, err, nil)
}

func Test_GenerateOrderNumber_Should_Return_16Digit_If_Order_with_Linepay_And_Kerry_UserID_4(t *testing.T) {
	// Format: YYMMDD PP SM UUU SEQ = 260106 98 22 004 004
	var expected int64 = 2601069822004004

	paymentMethodID := 2
	shippingMethodID := 1
	userID := 4
	seq := 4
	fixedTime := time.Date(2026, 1, 6, 0, 0, 0, 0, time.UTC)
	helper := order.OrderHelper{}

	actual, err := helper.GenerateOrderNumber(paymentMethodID, shippingMethodID, userID, seq, fixedTime)
	assert.Equal(t, expected, actual)
	assert.Equal(t, err, nil)
}

func Test_GenerateOrderNumber_Should_Return_16Digit_If_Order_with_Linepay_And_Thai_Post_UserID_5(t *testing.T) {
	// Format: YYMMDD PP SM UUU SEQ = 260106 98 33 005 005
	var expected int64 = 2601069833005005

	paymentMethodID := 2
	shippingMethodID := 2
	userID := 5
	seq := 5
	fixedTime := time.Date(2026, 1, 6, 0, 0, 0, 0, time.UTC)
	helper := order.OrderHelper{}

	actual, err := helper.GenerateOrderNumber(paymentMethodID, shippingMethodID, userID, seq, fixedTime)
	assert.Equal(t, expected, actual)
	assert.Equal(t, err, nil)
}

func Test_GenerateOrderNumber_Should_Return_16Digit_If_Order_with_Linepay_And_Lineman_UserID_6(t *testing.T) {
	// Format: YYMMDD PP SM UUU SEQ = 260106 98 44 006 006
	var expected int64 = 2601069844006006

	paymentMethodID := 2
	shippingMethodID := 3
	userID := 6
	seq := 6
	fixedTime := time.Date(2026, 1, 6, 0, 0, 0, 0, time.UTC)
	helper := order.OrderHelper{}

	actual, err := helper.GenerateOrderNumber(paymentMethodID, shippingMethodID, userID, seq, fixedTime)
	assert.Equal(t, expected, actual)
	assert.Equal(t, err, nil)
}

func Test_GenerateOrderNumber_Should_Throw_Error_If_SEQ_Is_Reach_Limit(t *testing.T) {
	var expected int64 = 0
	expectedErrorMessage := "Invalid sequence: 9999 (limit is 999)"

	paymentMethodID := 2
	shippingMethodID := 2
	userID := 1
	seq := 9999
	fixedTime := time.Date(2026, 1, 6, 0, 0, 0, 0, time.UTC)
	helper := order.OrderHelper{}

	actual, err := helper.GenerateOrderNumber(paymentMethodID, shippingMethodID, userID, seq, fixedTime)
	assert.Equal(t, expected, actual)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedErrorMessage)
}

func Test_GenerateOrderNumber_Should_Throw_Error_If_SEQ_Invalid(t *testing.T) {
	var expected int64 = 0
	expectedErrorMessage := "Invalid sequence: -1 (must be positive)"

	paymentMethodID := 2
	shippingMethodID := 2
	userID := 1
	seq := -1
	fixedTime := time.Date(2026, 1, 6, 0, 0, 0, 0, time.UTC)
	helper := order.OrderHelper{}

	actual, err := helper.GenerateOrderNumber(paymentMethodID, shippingMethodID, userID, seq, fixedTime)
	assert.Equal(t, expected, actual)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedErrorMessage)
}

func Test_GenerateOrderNumber_Should_Throw_Error_If_UserID_Is_Zero(t *testing.T) {
	var expected int64 = 0
	expectedErrorMessage := "Invalid userID: 0 (must be 1-999)"

	paymentMethodID := 1
	shippingMethodID := 1
	userID := 0
	seq := 1
	fixedTime := time.Date(2026, 1, 6, 0, 0, 0, 0, time.UTC)
	helper := order.OrderHelper{}

	actual, err := helper.GenerateOrderNumber(paymentMethodID, shippingMethodID, userID, seq, fixedTime)
	assert.Equal(t, expected, actual)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedErrorMessage)
}

func Test_GenerateOrderNumber_Should_Throw_Error_If_UserID_Exceeds_999(t *testing.T) {
	var expected int64 = 0
	expectedErrorMessage := "Invalid userID: 1000 (must be 1-999)"

	paymentMethodID := 1
	shippingMethodID := 1
	userID := 1000
	seq := 1
	fixedTime := time.Date(2026, 1, 6, 0, 0, 0, 0, time.UTC)
	helper := order.OrderHelper{}

	actual, err := helper.GenerateOrderNumber(paymentMethodID, shippingMethodID, userID, seq, fixedTime)
	assert.Equal(t, expected, actual)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedErrorMessage)
}

func Test_GenerateOrderNumber_Should_Handle_Max_UserID_999_And_Max_SEQ_999(t *testing.T) {
	// Format: YYMMDD PP SM UUU SEQ = 261231 98 44 999 999
	var expected int64 = 2612319844999999

	paymentMethodID := 2
	shippingMethodID := 3
	userID := 999
	seq := 999
	fixedTime := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
	helper := order.OrderHelper{}

	actual, err := helper.GenerateOrderNumber(paymentMethodID, shippingMethodID, userID, seq, fixedTime)
	assert.Equal(t, expected, actual)
	assert.Equal(t, err, nil)
}
