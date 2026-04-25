package utils

import "math"

func RoundToTenth(val float64) float64 {
	return math.Round(val*10) / 10
}

func KphToMph(kph float64) float64 {
	return kph * 0.621371
}