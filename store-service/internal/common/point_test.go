package common_test

import (
	"store-service/internal/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CalculatePoint_Input_599_Point_999_Should_be_Return_5(t *testing.T) {
	expected := 5

	actual := common.CalculatePoint(599.999)

	assert.Equal(t, expected, actual)
}

func Test_CalculatePoint_Input_Minus_599_Point_999_Should_be_Return_0(t *testing.T) {
	expected := 0

	actual := common.CalculatePoint(-599.999)

	assert.Equal(t, expected, actual)
}

func Test_CalculatePoint_Input_0_Point_999_Should_be_Return_0(t *testing.T) {
	expected := 0

	actual := common.CalculatePoint(0)

	assert.Equal(t, expected, actual)
}
