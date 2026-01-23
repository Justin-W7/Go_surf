package processing

import (
	"fmt"
	"time"

	"go_surf/api"
	"go_surf/models"
)

// AppendHourlyWeatherForecasts appends hourly weather forecasts to
// models.SurfForecast.PeriodForecasts to the provided models.SurfForecast.
//
// Parameters:
//   - surfForecasts: Slice of models.SurfForecast.
//
// Returns:
//   - NONE: only completes an action.
func AppendHourlyWeatherForecasts(surfForecasts []models.SurfForecast) {
	for i := range surfForecasts {
		URL := surfForecasts[i].SpotWeather.Properties.ForecastHourly

		rawHourlyWeatherData, err := api.FetchHourlyWeatherForecast(URL)
		if err != nil {
			fmt.Println("ERROR: ", err)
		}

		hourlyWeather, err := ParseHourlyWeatherForecast(rawHourlyWeatherData)
		if err != nil {
			fmt.Println("ERROR: ", err)
		}

		for j := range hourlyWeather.Properties.Periods {
			start := hourlyWeather.Properties.Periods[j].StartTime

			t, err := time.Parse(time.RFC3339, start)
			if err != nil {
				fmt.Println("ERROR: ", err)
			}

			surfSpotTime := surfForecasts[i].Timestamp

			if t.Unix()/3600 == surfSpotTime/3600 {
				surfForecasts[i].PeriodForecasts = append(
					surfForecasts[i].PeriodForecasts,
					hourlyWeather.Properties.Periods[j],
				)
			}
		}
	}
}
