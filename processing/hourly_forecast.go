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
		sForecast := &surfForecasts[i]

		URL := sForecast.SpotWeather.Properties.ForecastHourly

		rawHourlyWeatherData, err := api.FetchHourlyWeatherForecast(URL)
		if err != nil {
			fmt.Println("ERROR: ", err)
		}

		hourlyWeather, err := ParseHourlyWeatherForecast(rawHourlyWeatherData)
		if err != nil {
			fmt.Println("ERROR: ", err)
		}

		for j := range hourlyWeather.Properties.Periods {
			period := hourlyWeather.Properties.Periods[j]
			start := period.StartTime

			t, err := time.Parse(time.RFC3339, start)
			if err != nil {
				fmt.Println("ERROR: ", err)
			}

			sForecastTime := sForecast.Timestamp
			tUnix := t.Unix()

			if tUnix/3600 == sForecastTime/3600 {
				sForecast.PeriodForecasts = append(sForecast.PeriodForecasts, period)
			}
		}
	}
}
