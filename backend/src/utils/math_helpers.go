package utils

import "math"

func RoundToTenth(val float64) float64 {
	return math.Round(val*10) / 10
}
