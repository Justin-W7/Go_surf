package models

import "time"

type WeatherDatapoint struct {
	ID            int
	ObservedAt    time.Time
	RecordedAt    time.Time
	WindSpeed     *string
	WindDirection *string
	AirTemp       *float64
	Precipitation *float64
	CloudCoverage *string
}
