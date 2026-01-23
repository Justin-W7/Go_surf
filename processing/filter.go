package processing

import (
	"strings"

	"go_surf/models"
)

// SortAMPMForecasts filters a slice of SurfSpot models based on DateLocal.HH value.
// DateLocal.HH value uses 24-hour time.
//
// Parameters:
//   - forecasts: Slice of models.SurfSpot to filter.
//
// Returns:
//   - []models.SurfSpot: Slice containing forecasts for am times.
//   - []models.SurfSpot: Slice containing forecasts for pm times.
func SortAMPMForecasts(
	forecasts []models.SurfForecast,
) ([]models.SurfForecast, []models.SurfForecast) {
	var am []models.SurfForecast
	var pm []models.SurfForecast

	for _, i := range forecasts {
		hour := i.DateLocal.HH
		if hour >= 0 && hour <= 12 {
			am = append(am, i)
		}
		if hour >= 1 && hour <= 12 {
			pm = append(pm, i)
		}
	}
	return am, pm
}

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
