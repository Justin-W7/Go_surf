package processing

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"go_surf/models"
	"go_surf/utils"
)

func SummarizeTodaysForecast(forecast []models.SurfForecast) []models.SumTodaysForecast {
	var summaries []models.SumTodaysForecast
	seenSpotIDs := make(map[int]Values)

	// Iterate through forecast values - store unique SpotIDs and their indicies.
	for index, spot := range forecast {
		v, seen := seenSpotIDs[spot.SpotID]
		if !seen {
			v = Values{
				Seen:     true,
				Indicies: []int{},
			}
		}
		v.Indicies = append(v.Indicies, index)
		seenSpotIDs[spot.SpotID] = v
	}

	// Start of building forecast for each ID.
	for _, id := range seenSpotIDs {
		var summary models.SumTodaysForecast

		summary.SpotName = forecast[id.Indicies[0]].SpotName
		totalWaveHeight := 0.0
		totalQuality := 0.0
		// totalSwellPeriod := 0.0
		// totalSwellSize := 0.0
		// totalWaterTemp := 0.0
		totalAirTemp := 0.0
		totalWindSpeed := 0.0
		var arrWindDirection []string

		// iterate through all forecasts for the current ID
		for _, index := range id.Indicies {
			f := forecast[index]

			totalWaveHeight += f.SizeFt
			totalQuality += f.Shape
			totalAirTemp += float64(f.PeriodForecasts[0].Temperature)

			windFields := (strings.Fields(f.PeriodForecasts[0].WindSpeed))[0]
			windSpeed, err := strconv.Atoi(windFields)
			if err != nil {
				fmt.Println(err)
			}
			totalWindSpeed += float64(windSpeed)
			arrWindDirection = append(arrWindDirection, f.PeriodForecasts[0].WindDirection)
		}

		forecastCount := float64(len(id.Indicies))
		summary.AvgWaveHeight = utils.RoundToTenth(totalWaveHeight / forecastCount)
		summary.Quality = utils.RoundToTenth(totalQuality / forecastCount)
		summary.AirTemp = utils.RoundToTenth(totalAirTemp / forecastCount)
		summary.Wind.Direction, _ = avgWindDirection(arrWindDirection)
		summary.Wind.WindSpeed = utils.RoundToTenth(totalWindSpeed / forecastCount)

		summaries = append(summaries, summary)

	}

	return summaries
}

type Values struct {
	Seen     bool
	Indicies []int
}

// avgWindDirection calculates the average direction of the of the wind in a given forecast summary.
//
// Parameters:
//   - arrWindDirection: an array of wind directions. Ie: N, S, E, W
//
// Returns:
//   - closestCompassLabel: a string with a compass label
//   - avgDeg: the average bearing degree of the wind direction.
func avgWindDirection(arrWindDirection []string) (string, float64) {
	totalX := 0.0
	totalY := 0.0

	for _, bearing := range arrWindDirection {
		if bearing != "" {
			degree := float64(models.BearingMap[bearing])
			// convert degree to radians
			radian := degree * math.Pi / 180
			totalX += math.Cos(radian)
			totalY += math.Sin(radian)
		}
	}
	// calculate average degree
	avgRadian := math.Atan2(totalY, totalX)
	avgDeg := math.Mod(avgRadian*180/math.Pi+360, 360)

	// find closest bearing. String abv. eg: N, NW, S, NE ect.
	closestCompassLabel := ""
	minDiff := 360.0

	for label, bearing := range models.BearingMap {
		distance := math.Abs(float64(bearing) - avgDeg)

		if distance > 180 {
			distance = 360 - distance
		}
		if distance < minDiff {
			minDiff = distance
			closestCompassLabel = label
		}
	}
	return closestCompassLabel, avgDeg
}
