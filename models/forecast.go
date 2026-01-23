// Package models provides struct definitions for parsed json data.
package models

// SurfForecast struct and components
type SurfForecast struct {
	SpotID          int         `json:"spot_id"`
	SpotName        string      ""
	ID              string      `json:"_id"`
	DateGmt         DateInfo    `json:"date_gmt"`
	DateLocal       DateInfo    `json:"date_local"`
	IsDom           bool        `json:"is_dom"`
	Shape           float64     `json:"shape"`
	ShapeList       []ShapeItem `json:"shape_list"`
	Size            float64     `json:"size"`
	SizeFt          float64     `json:"size_ft"`
	SizeList        []SizeItem  `json:"size_list"`
	Timestamp       int64       `json:"timestamp"`
	Warnings        []string    `json:"warnings"`
	Coordinates     []float64
	SpotWeather     SpotWeather
	PeriodForecasts []HourlyPeriods
}

type DateInfo struct {
	DD int `json:"dd"`
	HH int `json:"hh"`
	MM int `json:"mm"`
	YY int `json:"yy"`
}

type ShapeItem struct {
	Influence int     `json:"influence"`
	Source    string  `json:"source"`
	Value     float64 `json:"value"`
}

type SizeItem struct {
	Shape      float64 `json:"shape"`
	SizeMeters float64 `json:"size_meters"`
}

// SumTodaysForecast struct
type SumTodaysForecast struct {
	SpotName      string
	AvgWaveHeight float64
	Quality       float64
	SwellPeriod   float64
	SwellSize     float64
	WaterTemp     int
	Wind          []SumWind
	AirTemp       float32
}

type SumWind struct {
	Time      string
	Direction string
	WindSpeed float64
}

// WeatherPointsResponse struct and components.
type WeatherPointsResponse struct {
	ID         string     `json:"id"`
	Type       string     `json:"type"`
	Geometry   Geometry   `json:"geometry"`
	Properties Properties `json:"properties"`
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Properties struct {
	CWA                 string           `json:"cwa"`
	ForecastOffice      string           `json:"forecastOffice"`
	GridID              string           `json:"gridId"`
	GridX               int              `json:"gridX"`
	GridY               int              `json:"gridY"`
	Forecast            string           `json:"forecast"`
	ForecastHourly      string           `json:"forecastHourly"`
	ForecastGridData    string           `json:"forecastGridData"`
	ObservationStations string           `json:"obsetvationStations"`
	RelativeLocation    RelativeLocation `json:"relativeLocation"`
	ForecastZone        string           `json:"forecastZone"`
	County              string           `json:"county"`
	FireWeatherZone     string           `json:"fireWeatherZone"`
	TimeZone            string           `json:"timeZone"`
	Radarstation        string           `json:"radarStation"`
}

type RelativeLocation struct {
	Type       string          `json:"type"`
	Geometry   Geometry        `json:"geometry"`
	Properties LocationDetails `json:"properties"`
	Distance   Measurement     `json:"distance"`
	Bearing    Measurement     `json:"bearing"`
}

type LocationDetails struct {
	City  string `json:"city"`
	State string `json:"state"`
}

type Measurement struct {
	Value    float64 `json:"value"`
	UnitCode string  `json:"unitCode"`
}
