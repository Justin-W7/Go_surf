package main

import (
	"fmt"

	"go_surf/api"
	"go_surf/models"
	"go_surf/processing"
)

// DataFilterKey is the location keyword used to filter surf spots.
const DataFilterKey = "newport"

func main() {
	// 1. Get data from Pitacst API
	rawSpotsData, err := api.FetchSpitcastSpots(api.SpitcastSpotURL)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// 2. Parse data into []models.SurfSpot.
	spotsArr, err := processing.ParseSurfSpots(rawSpotsData)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// 3. Filter surf spots based on DataFilterKey
	filteredSurfSpots := processing.FilterLocations(spotsArr, DataFilterKey)

	// 3.1. Fetch Forecast JSON data for surf spots.
	rawForecastData, err := api.FetchSpitcastForecast(filteredSurfSpots)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	// 3.2. Build todays forecast
	spotForecasts, err := processing.ParseSpotForecast(rawForecastData, filteredSurfSpots)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}

	// 4. Get weather data for each spot.
	var weatherSpotSlice []models.SpotWeather
	for _, i := range spotForecasts {
		// Get api response for weather point and parse into json
		tempWeatherPoint, _ := api.FetchWeatherPoint(i.Coordinates[0], i.Coordinates[1])
		tempSpotWeather, _ := processing.ParseSpotWeather(tempWeatherPoint)

		weatherSpotSlice = append(weatherSpotSlice, tempSpotWeather)
	}

	// 4.1. Append wind data to each spotForecast
	for _, i := range weatherSpotSlice {
		weatherForecastURL := i.Properties.Forecast
		weatherForecast, _ := api.FetchWeatherForecast(weatherForecastURL)

	}

	// 5. Build todays forecast for each spot.

	// 6. Write to JSON file.
}
