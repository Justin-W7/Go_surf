package processing

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"go_surf/models"
	"go_surf/utils"
)

// SummarizeTodaysForecast summarizes surf data for each spot in forecast. Calculates the
// averages for the characteristics of each spot.
//
// Parameters:
//   - forecast: Slice of SurfForecast data for multiple spots.
//
// Returns:
//   - []models.SumTodaysForecast: A slice summarizing data for each unique surf spot.
func SummarizeTodaysForecast(forecast []models.SurfForecast) []models.SumTodaysForecast {
	var summaries []models.SumTodaysForecast
	var spotIDs []int
	seenSpotIDs := make(map[int]bool)

	for _, surfForecast := range forecast {
		if !seenSpotIDs[surfForecast.SpotID] {
			spotIDs = append(spotIDs, surfForecast.SpotID)
			seenSpotIDs[surfForecast.SpotID] = true
		}
	}

	// make a summary for each location designated by the spotID.
	for _, spotID := range spotIDs {
		var spotSum models.SumTodaysForecast

		count := 0.0
		totalWaveHeight := 0.0
		totalQuality := 0.0
		totalWindSpeed := 0.0
		var arrWindDirection []string
		spotName := ""

		for i := range forecast {
			if spotID == forecast[i].SpotID {
				count += 1
				totalWaveHeight += forecast[i].SizeFt
				totalQuality += forecast[i].Shape
				spotName = forecast[i].SpotName

				// Get windspeed for each weather period
				for j := range forecast[i].PeriodForecasts {
					fmt.Println("Inside j := range forecast[i] loop")
					// parse wind speed "5 mph"
					windData := forecast[i].PeriodForecasts[j].WindSpeed
					windDataFields := strings.Fields(windData)

					// windSpeed is the "5", converted to an int
					windSpeed, err := strconv.Atoi(windDataFields[0])
					if err != nil {
						fmt.Println("ERROR in SummarizeTodaysForecast ",
							"windSpeed string conversion: ",
							err,
						)
					}
					totalWindSpeed += float64(windSpeed)

					// windDataFields[1] is the direction, eg: "NW"
					// collect them in arrWindDirection
					arrWindDirection = append(arrWindDirection, windDataFields[1])
				}
			}
		}
		// determine average wind direction
		avgWindDirection, _ := avgWindDirection(arrWindDirection)

		spotSum.AvgWaveHeight = utils.RoundToTenth(totalWaveHeight / count)
		spotSum.Wind.WindSpeed = utils.RoundToTenth(totalWindSpeed / count)
		spotSum.Wind.Direction = avgWindDirection
		spotSum.Quality = utils.RoundToTenth(totalQuality / count)
		spotSum.SpotName = spotName

		summaries = append(summaries, spotSum)
	}
	if len(summaries) == 0 {
		fmt.Println("summries forecast is empty")
	}
	return summaries
}

// avgWindDirection calculates the average direction of the of the wind in a given forecast.
//
// Parameters:
//   - arrWindDirection: an array of wind directions. Ie: N, S, E, W
//
// Returns:
// -
func avgWindDirection(arrWindDirection []string) (string, float64) {
	totalX := 0.0
	totalY := 0.0

	for _, bearing := range arrWindDirection {
		degree := float64(models.BearingMap[bearing])
		// convert degree to radians
		radian := degree * math.Pi / 180
		totalX += math.Cos(radian)
		totalY += math.Sin(radian)
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

func buildLocationSummary() {}

func avgWaveHeight() {}

func aveWaveQuality() {}

func avgWindSpeed() {}
