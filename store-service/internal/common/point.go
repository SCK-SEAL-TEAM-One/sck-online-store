package common

import "math"

func CalculatePoint(amount float64) int {
	if amount < 0 {
		return 0
	}

	points := int(math.Floor(amount / 100))
	return points
}
