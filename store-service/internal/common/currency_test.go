package common_test

import (
	"store-service/internal/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Decimal struct {
	ShortDecimal float64 `json:"short_digit"`
	LongDecimal  float64 `json:"long_digit"`
}

func Test_ConvertToThb_Input_123_Should_be_4424_Point_31_and_Point_305572(t *testing.T) {
	expected := Decimal{
		ShortDecimal: 4424.31,
		LongDecimal:  4424.305572,
	}

	actual := common.ConvertToThb(123)

	assert.Equal(t, expected.ShortDecimal, actual.ShortDecimal)
	assert.Equal(t, expected.LongDecimal, actual.LongDecimal)
}

func Test_ConvertToThb_Input_0_Should_be_0(t *testing.T) {
	expected := Decimal{
		ShortDecimal: 0,
		LongDecimal:  0,
	}

	actual := common.ConvertToThb(0)

	assert.Equal(t, expected.ShortDecimal, actual.ShortDecimal)
	assert.Equal(t, expected.LongDecimal, actual.LongDecimal)
}

func Test_ConvertToThb_Input_Minus_123_Should_be_Minus_4424_Point_31_and_Point_305572(t *testing.T) {
	expected := Decimal{
		ShortDecimal: -4424.31,
		LongDecimal:  -4424.305572,
	}

	actual := common.ConvertToThb(-123)

	assert.Equal(t, expected.ShortDecimal, actual.ShortDecimal)
	assert.Equal(t, expected.LongDecimal, actual.LongDecimal)
}
