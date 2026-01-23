package main

import (
	"fmt"

	"go_surf/api"
	"go_surf/processing"
	"go_surf/utils"
)

// DataFilterKey is the location keyword used to filter surf spots.
// We can use a name of a city as this is how the api sorts spots by address.
// i.e. useing "newport" returns all of the spots in Newport Beach.
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

	// 4. Fetch Forecast JSON data for surf spots.
	rawForecastData, err := api.FetchSpitcastForecast(filteredSurfSpots)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}

	// 5. Build todays forecast.
	// NOTE: spotForecasts is of []models.SurfForecast type.
	spotForecasts, err := processing.ParseSpotForecast(rawForecastData, filteredSurfSpots)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}

	// NOTE: todaysForecasts represents data for each hour of todays date for each spot.
	// todaysForecasts is of []models.SurfForecast type.
	todaysForecasts := processing.ParseTodaysForecasts(spotForecasts)

	// append SpotWeather to forecasts.
	processing.AppendSpotWeather(todaysForecasts)

	// Append hourly period to each forecast.
	processing.AppendHourlyWeatherForecasts(todaysForecasts)

	// sort into am/pm forecasts
	amForecasts, pmForecasts := processing.SortAMPMForecasts(todaysForecasts)

	// Build todays summary forecast
	amTodaysForecasts := processing.SummarizeTodaysForecast(amForecasts)
	pmTodaysForecasts := processing.SummarizeTodaysForecast(pmForecasts)

	// 6. Write to JSON file.
	utils.ToJSONFile(amTodaysForecasts, "am_forecasts")
	utils.ToJSONFile(pmTodaysForecasts, "pm_forecasts")
}
