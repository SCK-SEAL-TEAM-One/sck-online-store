package common_test

import (
	"store-service/internal/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CalculatePoint_Input_0_Should_be_Return_0(t *testing.T) {
	expected := 0
	price := 0.0

	actual := common.CalculatePoint(price)

	assert.Equal(t, expected, actual)
}

func Test_CalculatePoint_Input_100_Should_be_Return_1(t *testing.T) {
	expected := 1
	price := 100.0

	actual := common.CalculatePoint(price)

	assert.Equal(t, expected, actual)
}

func Test_CalculatePoint_Input_599_Dot_999_Should_be_Return_5(t *testing.T) {
	expected := 5
	price := 599.999

	actual := common.CalculatePoint(price)

	assert.Equal(t, expected, actual)
}

func Test_CalculatePoint_Input_Minus_599_Dot_999_Should_be_Return_0(t *testing.T) {
	expected := 0
	price := -599.999

	actual := common.CalculatePoint(price)

	assert.Equal(t, expected, actual)
}

