// Package processing parses, filters, and otherwise directs how the program handles data.
package processing

import (
	"encoding/json"
	"fmt"
	"go_surf/backend/src/models"
)

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

// ParseHourlyWeatherForecast parses raw JSON data
// into a HourlyWeahterForecast model.
//
// Parameters:
//   - data: Raw JSON data returned by FetchHourlyWeatherForecast().
//
// Returns:
//   - models.ParseHourlyWeatherForecast: Parse HourlyWeahterForecast data.
//   - error: AN error if the JSON unmarshalling fails.
func ParseHourlyWeatherForecast(data []byte) (models.HourlyWeatherForecast, error) {
	var hourlyWeather models.HourlyWeatherForecast
	if err := json.Unmarshal(data, &hourlyWeather); err != nil {
		fmt.Println("ERROR in ParseHourlyWeatherForecast(): ", err)
		return models.HourlyWeatherForecast{}, err
	}
	return hourlyWeather, nil
}

func ParseWeatherObservationStations(data []byte) (models.ObservationStationCollection, error) {
	var observationStations models.ObservationStationCollection
	if err := json.Unmarshal(data, &observationStations); err != nil {
		return models.ObservationStationCollection{}, fmt.Errorf("Could not unmarshal ObservationStationCollection: %w", err)
	}
	return observationStations, nil
}
