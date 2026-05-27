package models

import "time"

type CurrentSurfSpotConditions struct {
	ID                    int
	SpotId                int
	RecordedAt            time.Time
	DomSwellHeightM       *float64 // from buoy data
	DomSwellDir           *float64 // from buoy data
	WindSpeedMph          *string  // from city weather data
	WindDirection         *string  // from city weather data
	AirTempDegC           *float64 // from city weather data
	WaterTempDegC         *float64 // from buoy data
	Precipitation         *float64 // from city weather data
	CloudCoverage         *string  // from city weather data
	DominantWavePeriodSec *float64 // from buoy data
}

type Buoy struct {
	ID        int
	Latitude  float64
	Longitude float64
}

type BuoyData struct {
	BuoyID                int
	RecordedAt            time.Time
	WaveHeightM           *float64
	DominantWavePeriodSec *float64
	AvgWavePeriodSec      *float64
	MeanWaveDirectionDegT *float64
	WaterTempDegC         *float64
}

type City struct {
	ID        int
	Name      string
	Latitude  float64
	Longitude float64
}

type StaticSurfSpot struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	CityID      int     `json:"cityId"`
	NearestBuoy int     `json:"nearestBuoy"`
}
