package processing

import (
	"go_surf/backend/src/api"
	"go_surf/backend/src/models"
)

// AppendSpotWeather appends WeatherPoint data to models.SurfForecast.SpotWeather
//
// Parameters:
//   - surfForecasts: Slice of models.SurfForecast
//
// Return:
//   - NONE: only completes an action
func AppendSpotWeather(forecasts []models.SurfForecast) []models.SurfForecast {
	for i := range forecasts {
		weatherPoint, _ := api.FetchWeatherPoint(
			forecasts[i].Coordinates[0],
			forecasts[i].Coordinates[1],
		)
		spotWeather, _ := ParseSpotWeather(weatherPoint)
		forecasts[i].SpotWeather = spotWeather
	}
	return forecasts
}
