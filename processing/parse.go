// Package processing parses, filters, and otherwise directs how the program handles data.
package processing

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go_surf/models"
	"go_surf/utils"
)

// FilterLocations filters a slice of SurfSpot models based on a search string.
//
// Parameters:
//   - locations: Slice of models.SurfSpot to filter.
//   - filter: Substring to match against SpotIDChar (case-insensitive).
//
// Returns:
//   - []models.SurfSpot: A slice of spots whose SpotIDChar contains the filter string.
func FilterLocations(locations []models.SurfSpot, filter string) []models.SurfSpot {
	var namedSpots []models.SurfSpot
	for i := range locations {
		if strings.Contains(strings.ToLower(locations[i].SpotIDChar), strings.ToLower(filter)) {
			namedSpots = append(namedSpots, locations[i])
		}
	}
	return namedSpots
}

// SummarizeTodaysForecast summarizes surf data for each spot in forecast. Calculates the
// avg wave height and quality for each spot.
//
// Parameters:
//   - forecast: Slice of SurfForecast data for multiple spots.
//
// Returns:
//   - []models.SumTodaysForecast: A slice summarizing data for each unique surf spot.
func SummarizeTodaysForecast(forecast []models.SurfForecast) []models.SumTodaysForecast {
	var summary []models.SumTodaysForecast
	var checkID []int
	seen := make(map[int]bool)

	for _, elem := range forecast {
		if !seen[elem.SpotID] {
			checkID = append(checkID, elem.SpotID)
			seen[elem.SpotID] = true
		}
	}
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

		summary = append(summary, spotSum)
	}
	return summary
}

// ParseSurfSpots parses raw JSON data into a slice of SurfSpot models.
//
// Paramters
//   - data: Raw JSON data that is returned by the Spitcast API.
//
// Returns:
//   - []models.SurfSpot: Slice of SurfSpot models.
//   - error: An error if the JSON unmarshalling fails.
func ParseSurfSpots(data []byte) ([]models.SurfSpot, error) {
	var spots []models.SurfSpot
	if err := json.Unmarshal(data, &spots); err != nil {
		return nil, err
	}
	return spots, nil
}

// ParseSpotForecast parses multiple spot forecasts from raw JSON data and
// merges them with filtered spot information.
//
// Paramters:
//   - data: Slices of raw JSON byte slices, each respresenting a spot forecast.
//   - filteredSpots: Slice of SurfSpot models.
//
// Returns:
//   - []models.SurfForecast: Combined forecast data with spot names and coordinates.
//   - error: An error if the JSON unmarshalling fails.
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
				forecasts[elem].Coordinates = filteredSpots[spot].Coordinates
			}
		}
	}
	return forecasts, nil
}

// ParseTodaysForecasts filters a slice of SurfForecast to only include
// today's forecasts.
//
// Paramters:
// - data: Slice of SurfForecasts.
//
// Returns:
// - []models.SurfForecast: Forecasts for surf spots on the current day.
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

// ParseWeatherForecast parses raw JSON data into a WeatherForecast model.
//
// Parameters:
//   - data: Raw JSON data returned by FetchWeatherForecast().
//
// Returns:
//   - models.WeatherForecast: Parsed weather forecast data.
//   - error: An error if the JSON unmarshalling fails.
func ParseWeatherForecast(data []byte) (models.WeatherForecast, error) {
	var weather models.WeatherForecast
	if err := json.Unmarshal(data, &weather); err != nil {
		fmt.Println("ERROR in ParseWeatherForecast(): ", err)
		return models.WeatherForecast{}, err
	}
	return weather, nil
}

// ParseWeatherGridForecast parses raw JSON data into a ForecastGridData models.
//
// Parameters:
//   - data: Raw JSON data returned by FetchWeahterGridForecast().
//
// Returns:
//   - models.ForecastGridData: Parsed ForecastGrid data.
//   - error: An error if the JOSN unmarshalling fails.
func ParseWeatherGridForecast(data []byte) (models.ForecastGridData, error) {
	var gridData models.ForecastGridData
	if err := json.Unmarshal(data, &gridData); err != nil {
		fmt.Println("ERROD in ParseWeatherGridForecast(): ", err)
		return models.ForecastGridData{}, err
	}
	return gridData, nil
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
