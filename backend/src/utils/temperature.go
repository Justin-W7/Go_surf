package utils


func FahrenheitToCelsius(f float64) float64 {
	return (f - 32) * 5 / 9
}

func CelsiusToFahrenheit(c float64) float64 {
	return (c * 9 / 5) + 32
}