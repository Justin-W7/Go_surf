package processing

import (
	"go_surf/models"
	"go_surf/utils"
)

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

	for elem := range checkID {
		var spotSum models.SumTodaysForecast

		// calculate avg waveheight data
		count := 0.0
		waveHTotal := 0.0
		qualityTotal := 0.0
		// windSpeedTotal := 0
		name := ""

		for _, spot := range forecast {
			if elem == spot.SpotID {

				// calculate wave data totals
				waveHTotal += spot.SizeFt
				qualityTotal += spot.Shape
				// calculate windspeed total
			}
		}

		spotSum.AvgWaveHeight = utils.RoundToTenth(waveHTotal / count)
		spotSum.Quality = utils.RoundToTenth(qualityTotal / count)

		spotSum.SpotName = name

		summary = append(summary, spotSum)
	}
	return summary
}
