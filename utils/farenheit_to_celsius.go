package utils

import "math"

func FarenheitToCelsius(f float64) float64 {
	return math.Round((f - 32) * 5 / 9)
}
