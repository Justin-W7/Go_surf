package main

import (
	"log"

	"go_surf/api"
	"go_surf/processing"
	"go_surf/utils"
)

const DataFilterKey = "newport"

func main() {
	rawSpots, err := api.FetchSpitcastSpots(api.SpitcastSpotURL)
	if err != nil {
		log.Fatalf("FetchSpitcastSpot main error: %v", err)
		return
	}

	parsedSpots, err := processing.ParseSurfSpots(rawSpots)
	if err != nil {
		log.Fatalf("ParseSurfSpots error: %v", err)
	}

	filteredSpots := processing.FilterLocations(parsedSpots, DataFilterKey)

	rawForecasts, err := api.FetchSpitcastForecast(filteredSpots)
	if err != nil {
		log.Fatalf("FetchSpitcastForecast error: %v", err)
	}

	forecastData, err := processing.ParseSpotForecast(rawForecasts, filteredSpots)
	if err != nil {
		log.Fatalf("PareseSpotForecast error: %v", err)
	}

	todaysForecast := processing.ParseTodaysForecasts(forecastData)
	todaysSummary := processing.SummarizeTodaysForecast(todaysForecast)
	utils.ToJSONFile(todaysSummary, "todays_summary")
}
