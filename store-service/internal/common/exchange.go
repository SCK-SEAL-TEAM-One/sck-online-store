package common

import "math"

type Digit struct {
	ShortDigit float64 `json:"short_digit"`
	LongDigit  float64 `json:"long_digit"`
}

func ConvertToThb(amount float64) Digit {
	rate := 35.969964
	result := amount * rate
	factor2 := math.Pow(10, 2)
	factor6 := math.Pow(10, 6)

	return Digit{
		ShortDigit: math.Round(result*factor2) / factor2,
		LongDigit:  math.Round(result*factor6) / factor6,
	}
}
