package main

import (
	"log"

	"go_surf/api"
	"go_surf/processing"
	"go_surf/utils"
)

// DataFilterKey is the location keyword used to filter surf spots.
const DataFilterKey = "newport"

func main() {
	// Fetch raw surf spot data from Spitcast.
	rawSpots, err := api.FetchSpitcastSpots(api.SpitcastSpotURL)
	if err != nil {
		log.Fatalf("FetchSpitcastSpot main error: %v", err)
		return
	}

	// Parse and filter spots by the configured location.
	parsedSpots, err := processing.ParseSurfSpots(rawSpots)
	if err != nil {
		log.Fatalf("ParseSurfSpots error: %v", err)
	}
	filteredSpots := processing.FilterLocations(parsedSpots, DataFilterKey)

	// Fetch and parse forecast data for the filtered spots
	rawForecasts, err := api.FetchSpitcastForecast(filteredSpots)
	if err != nil {
		log.Fatalf("FetchSpitcastForecast error: %v", err)
	}

	forecastData, err := processing.ParseSpotForecast(rawForecasts, filteredSpots)
	if err != nil {
		log.Fatalf("PareseSpotForecast error: %v", err)
	}

	// Get wind data and append direction and speed to forecastData[n]

	// Generate today's summary and persist it as JSON.
	todaysForecast := processing.ParseTodaysForecasts(forecastData)
	todaysSummary := processing.SummarizeTodaysForecast(todaysForecast)
	utils.ToJSONFile(todaysSummary, "todays_summary")
}
