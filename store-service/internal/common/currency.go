package common

import "math"

type Decimal struct {
	ShortDecimal float64 `json:"short_digit"`
	LongDecimal  float64 `json:"long_digit"`
}

func ConvertToThb(amount float64) Decimal {
	rate := 35.969964
	result := amount * rate
	factor2 := math.Pow(10, 2)
	factor6 := math.Pow(10, 6)

	return Decimal{
		ShortDecimal: math.Round(result*factor2) / factor2,
		LongDecimal:  math.Round(result*factor6) / factor6,
	}
}
