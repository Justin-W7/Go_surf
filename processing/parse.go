// Package processing parses, filters, and otherwise
// directs how the program formats data.
package processing

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go_surf/models"
	"go_surf/utils"
)

func ParseSurfSpots(data []byte) ([]models.SurfSpot, error) {
	var spots []models.SurfSpot
	if err := json.Unmarshal(data, &spots); err != nil {
		return nil, err
	}
	return spots, nil
}

func FilterLocations(locations []models.SurfSpot, filter string) []models.SurfSpot {
	var namedSpots []models.SurfSpot
	for i := range locations {
		if strings.Contains(strings.ToLower(locations[i].SpotIDChar), strings.ToLower(filter)) {
			namedSpots = append(namedSpots, locations[i])
		}
	}
	return namedSpots
}

func ParseSpotForecast(data [][]byte, filteredSpots []models.SurfSpot) ([]models.SurfForecast, error) {
	var forecasts []models.SurfForecast
	for spot := range data {
		var spotForecast []models.SurfForecast
		if err := json.Unmarshal(data[spot], &spotForecast); err != nil {
			return nil, err
		}
		forecasts = append(forecasts, spotForecast...)
	}

	for spot := range filteredSpots {
		for elem := range forecasts {
			if filteredSpots[spot].SpotID == forecasts[elem].SpotID {
				forecasts[elem].SpotName = filteredSpots[spot].SpotName
			}
		}
	}
	return forecasts, nil
}

func ParseTodaysForecasts(data []models.SurfForecast) []models.SurfForecast {
	var todaysForecast []models.SurfForecast

	today := time.Now()
	tYear := today.Year()
	tMonth := int(today.Month())
	tDay := today.Day()

	for spot := range data {
		t := time.Unix(int64(data[spot].Timestamp), 0)
		year, month, day := t.Date()

		if tDay == day && tMonth == int(month) && tYear == year {
			todaysForecast = append(todaysForecast, data[spot])
		}
	}
	return todaysForecast
}

func SummarizeTodaysForecast(forecast []models.SurfForecast) []models.SumTodaysForecast {
	var summary []models.SumTodaysForecast
	// create slice of spot IDs (no duplicates)
	var checkID []int
	seen := make(map[int]bool)

	for _, elem := range forecast {
		if !seen[elem.SpotID] {
			checkID = append(checkID, elem.SpotID)
			seen[elem.SpotID] = true
		}
	}
	// Build summary for each spot
	for _, elem := range checkID {
		var spotSum models.SumTodaysForecast

		waveHTotal := 0.0
		count := 0.0
		qualityTotal := 0.0
		name := ""

		for _, spot := range forecast {
			if elem == spot.SpotID {
				waveHTotal += spot.SizeFt
				count += 1
				qualityTotal += spot.Shape
				name = spot.SpotName
			}
		}

		spotSum.SpotName = name
		spotSum.AvgWaveHeight = utils.RoundToTenth(waveHTotal / count)
		spotSum.Quality = utils.RoundToTenth(qualityTotal / count)

		// Append to final summary slice
		summary = append(summary, spotSum)
	}
	return summary
}

func ParseSpotWeather(data []byte) (models.SpotWeather, error) {
	var weatherData models.SpotWeather
	if err := json.Unmarshal(data, &weatherData); err != nil {
		fmt.Println("Unmarshal error in processing.ParseSpotWeather.", err)
		return models.SpotWeather{}, err
	}
	return weatherData, nil
}

func ParseWeatherForecast(data []byte) (models.WeatherForecast, error) {
	var weather models.WeatherForecast
	if err := json.Unmarshal(data, &weather); err != nil {
		fmt.Println("ERROR in ParseWeatherForecast(): ", err)
		return models.WeatherForecast{}, err
	}
	return weather, nil
}

func ParseWeatherGridForecast(data []byte) (models.ForecastGridData, error) {
	var gridData models.ForecastGridData
	if err := json.Unmarshal(data, &gridData); err != nil {
		fmt.Println("ERROD in ParseWeatherGridForecast(): ", err)
		return models.ForecastGridData{}, err
	}
	return gridData, nil
}

func ParseHourlyWeatherForecast(data []byte) (models.HourlyWeatherForecast, error) {
	var hourlyWeather models.HourlyWeatherForecast
	if err := json.Unmarshal(data, &hourlyWeather); err != nil {
		fmt.Println("ERROR in ParseHourlyWeatherForecast(): ", err)
		return models.HourlyWeatherForecast{}, err
	}
	return hourlyWeather, nil
}
