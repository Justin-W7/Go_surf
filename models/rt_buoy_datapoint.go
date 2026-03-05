package models

import "time"

type BuoyDataPoint struct {
	BuoyID                int
	RecordedAt            time.Time
	WindDirectionDegT     *float64
	WindSpeedMetersPerSec *float64
	WindGustMetersPerSec  *float64
	WaveHeightM           *float64
	DominantWavePeriodSec *float64
	AvgWavePeriodSec      *float64
	MeanWaveDirectionDegT *float64
	AirTempDegC           *float64
	WaterTempDegC         *float64
}
