// Package api provides functions for making https requests to
// external surf and weather API's.
package api

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"go_surf/models"
)

// fetchURL performs an HTTP GET request for the given URL and returns
// the response body.
//
// Parameters:
//   - url: The API endpoint URl to fetch.
//
// Returns:
// - []byte: The raw response body.
// - error: An error if the request fails.
func fetchURL(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed request to %s: %w", url, err)
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response from %s: %w", url, err)
	}

	return data, nil
}

// FetchSpitcastSpots sends an HTTP GET request to the provided url.
//
// Parameters:
//   - url: A Spitcast endpoint URL that returns
//     surf spot data (a JSON lsit of spot metadata).
//
// Returns:
//   - []byte: A byte slice containing the raw response body.
//   - error: An error if the http request fails.
func FetchSpitcastSpots(url string) ([]byte, error) {
	return fetchURL(url)
}

// FetchSpitcastForecast retrieves surf forecasts for a list of surf spots.
//
// Parameters:
//   - spots: A slice of SurfSpot models representing the surf spots.
//     SpotID, which is used to construct a Spitcast forecast request.
//
// The function uses the current date (year, month, and day) when building
// forecast URLs.
//
// Returns:
//   - [][]byte: A slice of byte slices, where each element is the raw forecast response
//     for a corresponding surf spot.
//   - error: An error if any HTTP request fails.
func FetchSpitcastForecast(spots []models.SurfSpot) ([][]byte, error) {
	year := time.Now().Year()
	month := int(time.Now().Month())
	day := time.Now().Day()
	locations := make([][]byte, 0)

	for i := range spots {
		url := fmt.Sprintf(SpitcastForecastURL, spots[i].SpotID, year, month, day)
		data, err := fetchURL(url)
		if err != nil {
			return nil, err
		}
		locations = append(locations, data)
	}
	return locations, nil
}

// FetchWeatherPoint sends an HTTP GET request to the National Weather Service
// API to retrieve weather point data for a given longitude and latitude.
//
// Parameters:
//   - long: Longitude of the location.
//   - lat: Latitude of the location.
//
// Returns:
//   - []byte: A byte slice containing the raw response body.
//   - error: An error if the HTTP request fails.
func FetchWeatherPoint(lat float64, long float64) ([]byte, error) {
	url := fmt.Sprintf(NWSWeatherURL, long, lat)
	return fetchURL(url)
}

// FetchWeatherForecast sends an HTTP GET request to retrieve a weather forecast
// from a given NWS (National Weather Service) URL.
//
// Parameters:
//   - url: The API endpoint URL for the forecast.
//
// Returns:
//   - []byte: A byte slice containing the raw response body.
//   - error: An error if the HTTP request fails.
func FetchWeatherForecast(url string) ([]byte, error) {
	return fetchURL(url)
}

// FetchHourlyWeatherForecast retrieves hourly weather forecast data
// from a specified NWS URL.
//
// Parameters:
//   - url: The full API endpoint URL for the hourly forecast.
//
// Returns:
//   - []byte: The raw response body containing hourly forecast data.
//   - error: An error if the HTTP request fails.
func FetchHourlyWeatherForecast(url string) ([]byte, error) {
	return fetchURL(url)
}

// FetchWeatherGridForecast retrieves grid-based weather forecast data
// from a specific NWS url.
//
// Paramters:
//   - url: The API endpoint URL for the grid forecast.
//
// Returns:
//
//   - []byte: The raw response body containing grid forecast data.
//   - error: An errif the HTTP request fails.
func FetchWeatherGridForecast(url string) ([]byte, error) {
	return fetchURL(url)
}
