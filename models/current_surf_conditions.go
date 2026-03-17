package models

import "time"

type CurrentSurfSpotConditions struct {
	SpotID                int
	RecordedAt            time.Time
	WaveHeightM           *float64
	WindSpeedMetersPerSec *float64
	WindDirectionDegT     *float64
	AirTempDegC           *float64
	WaterTempDegC         *float64
	Precipitation         *float64
	ClodCoverage          *string
}
