package models

import "time"

type WeatherDatapoint struct {
	CityID        int
	RecordedAt    time.Time
	WindSpeed     *string
	WindDirection *string
	AirTemp       *float64
	Precipitation *float64
	CloudCoverage *string
	ObservedAt    time.Time
}


