package order_test

import (
	"store-service/internal/order"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateOrderNumber_Should_Return_2601069522001_If_Order_with_Credit_Card_And_Kerry(t *testing.T) {
	expected := "2601069522001"

	paymentMethodID := 1
	shippingMethodID := 1
	seq := 1
	fixedTime := time.Date(2026, 1, 6, 0, 0, 0, 0, time.UTC)
	helper := order.OrderHelper{}

	actual, err := helper.GenerateOrderNumber(paymentMethodID, shippingMethodID, seq, fixedTime)
	assert.Equal(t, expected, actual)
	assert.Equal(t, err, nil)
}

func Test_GenerateOrderNumber_Should_Return_2601069533002_If_Order_with_Credit_Card_And_Thai_Post(t *testing.T) {
	expected := "2601069533002"

	paymentMethodID := 1
	shippingMethodID := 2
	seq := 2
	fixedTime := time.Date(2026, 1, 6, 0, 0, 0, 0, time.UTC)
	helper := order.OrderHelper{}

	actual, err := helper.GenerateOrderNumber(paymentMethodID, shippingMethodID, seq, fixedTime)
	assert.Equal(t, expected, actual)
	assert.Equal(t, err, nil)
}

func Test_GenerateOrderNumber_Should_Return_2601069544003_If_Order_with_Credit_Card_And_Lineman(t *testing.T) {
	expected := "2601069544003"

	paymentMethodID := 1
	shippingMethodID := 3
	seq := 3
	fixedTime := time.Date(2026, 1, 6, 0, 0, 0, 0, time.UTC)
	helper := order.OrderHelper{}

	actual, err := helper.GenerateOrderNumber(paymentMethodID, shippingMethodID, seq, fixedTime)
	assert.Equal(t, expected, actual)
	assert.Equal(t, err, nil)
}

func Test_GenerateOrderNumber_Should_Return_2601069822004_If_Order_with_Linepay_And_Kerry(t *testing.T) {
	expected := "2601069822004"

	paymentMethodID := 2
	shippingMethodID := 1
	seq := 4
	fixedTime := time.Date(2026, 1, 6, 0, 0, 0, 0, time.UTC)
	helper := order.OrderHelper{}

	actual, err := helper.GenerateOrderNumber(paymentMethodID, shippingMethodID, seq, fixedTime)
	assert.Equal(t, expected, actual)
	assert.Equal(t, err, nil)
}

func Test_GenerateOrderNumber_Should_Return_2601069833005_If_Order_with_Linepay_And_Thai_Post(t *testing.T) {
	expected := "2601069833005"

	paymentMethodID := 2
	shippingMethodID := 2
	seq := 5
	fixedTime := time.Date(2026, 1, 6, 0, 0, 0, 0, time.UTC)
	helper := order.OrderHelper{}

	actual, err := helper.GenerateOrderNumber(paymentMethodID, shippingMethodID, seq, fixedTime)
	assert.Equal(t, expected, actual)
	assert.Equal(t, err, nil)
}

func Test_GenerateOrderNumber_Should_Return_2601069844006_If_Order_with_Linepay_And_Lineman(t *testing.T) {
	expected := "2601069844006"

	paymentMethodID := 2
	shippingMethodID := 3
	seq := 6
	fixedTime := time.Date(2026, 1, 6, 0, 0, 0, 0, time.UTC)
	helper := order.OrderHelper{}

	actual, err := helper.GenerateOrderNumber(paymentMethodID, shippingMethodID, seq, fixedTime)
	assert.Equal(t, expected, actual)
	assert.Equal(t, err, nil)
}

func Test_GenerateOrderNumber_Should_Throw_Error_If_SEQ_Is_Reach_Limit(t *testing.T) {
	expected := ""
	expectedErrorMessage := "Invalid sequence: 9999 (limit is 999)"

	paymentMethodID := 2
	shippingMethodID := 2
	seq := 9999
	fixedTime := time.Date(2026, 1, 6, 0, 0, 0, 0, time.UTC)
	helper := order.OrderHelper{}

	actual, err := helper.GenerateOrderNumber(paymentMethodID, shippingMethodID, seq, fixedTime)
	assert.Equal(t, expected, actual)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedErrorMessage)
}

func Test_GenerateOrderNumber_Should_Throw_Error_If_SEQ_Invalid(t *testing.T) {
	expected := ""
	expectedErrorMessage := "Invalid sequence: -1 (must be positive)"

	paymentMethodID := 2
	shippingMethodID := 2
	seq := -1
	fixedTime := time.Date(2026, 1, 6, 0, 0, 0, 0, time.UTC)
	helper := order.OrderHelper{}

	actual, err := helper.GenerateOrderNumber(paymentMethodID, shippingMethodID, seq, fixedTime)
	assert.Equal(t, expected, actual)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedErrorMessage)
}

func Test_GetNextSequence_Should_Return_2_If_Last_Order_Number_Is_2601149844001(t *testing.T) {
	expected := 2
	lastOrderNumber := "2601149855001"
	helper := order.OrderHelper{}

	actual, err := helper.GetNextSequence(lastOrderNumber)

	assert.Equal(t, expected, actual)
	assert.Equal(t, err, nil)
}

func Test_GetNextSequence_Should_Return_50_If_Last_Order_Number_Is_2601149522049(t *testing.T) {
	expected := 50
	lastOrderNumber := "2601149522049"
	helper := order.OrderHelper{}

	actual, err := helper.GetNextSequence(lastOrderNumber)

	assert.Equal(t, expected, actual)
	assert.Equal(t, err, nil)
}

func Test_GetNextSequence_Should_Return_100_If_Last_Order_Number_Is_2602289522099(t *testing.T) {
	expected := 100
	lastOrderNumber := "2602289522099"
	helper := order.OrderHelper{}

	actual, err := helper.GetNextSequence(lastOrderNumber)

	assert.Equal(t, expected, actual)
	assert.Equal(t, err, nil)
}

func Test_GetNextSequence_Should_Return_200_If_Last_Order_Number_Is_2603159533199(t *testing.T) {
	expected := 200
	lastOrderNumber := "2603159533199"
	helper := order.OrderHelper{}

	actual, err := helper.GetNextSequence(lastOrderNumber)

	assert.Equal(t, expected, actual)
	assert.Equal(t, err, nil)
}

func Test_GetNextSequence_Should_Return_Err_Invalid_Length_If_Last_Order_Number_Is_260315953301(t *testing.T) {
	expected := 0
	expectedErrorMessage := "Invalid length for order 260315953301"
	lastOrderNumber := "260315953301"
	helper := order.OrderHelper{}

	actual, err := helper.GetNextSequence(lastOrderNumber)

	assert.Equal(t, expected, actual)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedErrorMessage)
}

func Test_GetNextSequence_Should_Return_Err_Invalid_Length_If_Last_Order_Number_Is_26031595339999(t *testing.T) {
	expected := 0
	expectedErrorMessage := "Invalid length for order 26031595339999"
	lastOrderNumber := "26031595339999"
	helper := order.OrderHelper{}

	actual, err := helper.GetNextSequence(lastOrderNumber)

	assert.Equal(t, expected, actual)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedErrorMessage)
}

func Test_GetNextSequence_Should_Return_Err_Invalid_Length_If_Last_Order_Number_Is_2603159533XYZ(t *testing.T) {
	expected := 0
	expectedErrorMessage := "Invalid sequence number: XYZ"
	lastOrderNumber := "2603159533XYZ"
	helper := order.OrderHelper{}

	actual, err := helper.GetNextSequence(lastOrderNumber)

	assert.Equal(t, expected, actual)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedErrorMessage)
}

func Test_GetNextSequence_Should_Return_Err_Invalid_Length_If_Last_Order_Number_Is_2603159544000(t *testing.T) {
	expected := 0
	expectedErrorMessage := "Invalid sequence number: 000"
	lastOrderNumber := "2603159533000"
	helper := order.OrderHelper{}

	actual, err := helper.GetNextSequence(lastOrderNumber)

	assert.Equal(t, expected, actual)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedErrorMessage)
}

func Test_GetNextSequence_Should_Return_999_Invalid_Length_If_Last_Order_Number_Is_2611129544998(t *testing.T) {
	expected := 999
	lastOrderNumber := "2611129544998"
	helper := order.OrderHelper{}

	actual, err := helper.GetNextSequence(lastOrderNumber)

	assert.Equal(t, expected, actual)
	assert.Equal(t, err, nil)
}

func Test_GetNextSequence_Should_Return_Err_Annual_Order_Limit_999_Reached_If_Last_Order_Number_Is_2612129822999(t *testing.T) {
	expected := 0
	expectedErrorMessage := "Annual order limit (999) reached"
	lastOrderNumber := "2612129822999"
	helper := order.OrderHelper{}

	actual, err := helper.GetNextSequence(lastOrderNumber)

	assert.Equal(t, expected, actual)
	assert.Error(t, err)
	assert.EqualError(t, err, expectedErrorMessage)
}
