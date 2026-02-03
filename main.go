package main

import (
	"fmt"

	"go_surf/api"
	"go_surf/models"
	"go_surf/processing"
	"go_surf/utils"
)

// DataFilterKey is the location keyword used to filter surf spots.
// We can use a name of a city as this is how the api sorts spots by address.
// i.e. useing "newport" returns all of the spots in Newport Beach.
const DataFilterKey = "newport"

func main() {
	rawSpotsData, err := api.FetchSpitcastSpots(api.SpitcastSpotURL)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	spotsArr, err := processing.ParseSurfSpots(rawSpotsData)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	filteredSurfSpots := processing.FilterLocations(spotsArr, DataFilterKey)

	rawForecastData, err := api.FetchSpitcastForecast(filteredSurfSpots)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}

	// NOTE: spotForecasts is of []models.SurfForecast type.
	spotForecasts, err := processing.ParseSpotForecast(rawForecastData, filteredSurfSpots)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}

	// NOTE: todaysForecasts represents data for each hour of todays date for each spot.
	// todaysForecasts is of []models.SurfForecast type.
	todaysForecasts := processing.ParseTodaysForecasts(spotForecasts)

	todaysForecasts = processing.AppendSpotWeather(todaysForecasts)

	processing.AppendHourlyWeatherForecasts(todaysForecasts)

	// Only include forecasts that have hourly weather data
	var todaysFullForecasts []models.SurfForecast
	for i := range todaysForecasts {
		if todaysForecasts[i].PeriodForecasts != nil {
			todaysFullForecasts = append(todaysFullForecasts, todaysForecasts[i])
		}
	}

	amForecasts, pmForecasts := processing.SortAMPMForecasts(todaysFullForecasts)

	amTodaysForecasts := processing.SummarizeTodaysForecast(amForecasts)
	pmTodaysForecasts := processing.SummarizeTodaysForecast(pmForecasts)

	utils.ToJSONFile(amTodaysForecasts, "amTodaysforecasts")
	utils.ToJSONFile(pmTodaysForecasts, "pmTodaysforecasts")
}
