package processing

import (
	"encoding/json"
	"fmt"

	"go_surf/api"
	"go_surf/models"
)

// AppendSpotWeather appends WeatherPoint data to models.SurfForecast.SpotWeather
//
// Parameters:
//   - surfForecasts: Slice of models.SurfForecast
//
// Return:
//   - NONE: only completes an action
func AppendSpotWeather(forecasts []models.SurfForecast) {
	for i := range forecasts {
		weatherPoint, _ := api.FetchWeatherPoint(
			forecasts[i].Coordinates[0],
			forecasts[i].Coordinates[1],
		)
		spotWeather, _ := ParseSpotWeather(weatherPoint)
		forecasts[i].SpotWeather = spotWeather
	}
}

// ParseSpotWeather parses raw JSON data into a SpotWeather model.
//
// Parameters:
//   - data: Raw JSON data returned by FetchWeatherPoint().
//
// Returns:
//   - models.SpotWeather: Parsed SpotWeather data.
//
// - error: An error if the JSON unmarshalling fails.
func ParseSpotWeather(data []byte) (models.SpotWeather, error) {
	var weatherData models.SpotWeather
	if err := json.Unmarshal(data, &weatherData); err != nil {
		fmt.Println("Unmarshal error in processing.ParseSpotWeather.", err)
		return models.SpotWeather{}, err
	}
	return weatherData, nil
}
