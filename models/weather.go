package models

import (
	"time"
)

// -------------------------------------------------
// -- Weather structs for initial spot api call
// -------------------------------------------------

type SpotWeather struct {
	Context    []any           `json:"@context"`
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	Geometry   PointGeometry   `json:"geometry"`
	Properties PointProperties `json:"properties"`
}

type PointGeometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type PointProperties struct {
	ID                  string                `json:"@id"`
	Type                string                `json:"@type"`
	CWA                 string                `json:"cwa"`
	ForecastOffice      string                `json:"forecastOffice"`
	GridID              string                `json:"gridId"`
	GridX               int                   `json:"gridX"`
	GridY               int                   `json:"gridY"`
	Forecast            string                `json:"forecast"`
	ForecastHourly      string                `json:"forecastHourly"`
	ForecastGridData    string                `json:"forecastGridData"`
	ObservationStations string                `json:"observationStations"`
	RelativeLocation    PointRelativeLocation `json:"relativeLocation"`
	ForecastZone        string                `json:"forecastZone"`
	TimeZone            string                `json:"timeZone"`
	RadarStation        string                `json:"radarStation"`
}

type PointRelativeLocation struct {
	Geometry           PointGeometry      `json:"geometry"`
	LocationProperties LocationProperties `json:"properties"`
}

type LocationProperties struct {
	City     string   `json:"city"`
	State    string   `json:"state"`
	Distance Distance `json:"distance"`
	Bearing  Bearing  `json:"bearing"`
}

type Distance struct {
	UnitCode string  `json:"unitCode"`
	Value    float64 `json:"value"`
}

type Bearing struct {
	UnitCode string  `json:"unitCode"`
	Value    float64 `json:"value"`
}

//--------------------------------------------------------------------
// Weather Forecast Structs WeatherForecast
//--------------------------------------------------------------------

type WeatherForecast struct {
	Context    []any              `json:"@context"`
	Type       string             `json:"type"`
	Geometry   ForecastGeometry   `json:"geometry"`
	Properties ForecastProperties `json:"properties"`
}

type ForecastGeometry struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

type ForecastProperties struct {
	Units             string           `json:"units"`
	ForecastGenerator string           `json:"forecastGenerator"`
	GeneratedAt       string           `json:"generatedAt"`
	UpdateTime        string           `json:"updateTime"`
	ValidTimes        string           `json:"validTimes"`
	Elevation         QuantValue       `json:"elevation"`
	Periods           []ForecastPeriod `json:"periods"`
}

type QuantValue struct {
	UnitCode string   `json:"unitCode"`
	Value    *float64 `json:"value"`
}

type ForecastPeriod struct {
	Number                     int        `json:"number"`
	Name                       string     `json:"name"`
	StartTime                  time.Time  `json:"startTime"`
	EndTime                    time.Time  `json:"endTime"`
	IsDaytime                  bool       `json:"isDaytime"`
	Temperature                int        `json:"temperature"`
	TemperatureUnit            string     `json:"temperatureUnit"`
	TemperatureTrend           *string    `json:"temperatureTrend"`
	ProbabilityOfPrecipitation QuantValue `json:"probabilityOfPrecipitation"`
	WindSpeed                  string     `json:"windSpeed"`
	WindDirection              string     `json:"windDirection"`
	Icon                       string     `json:"icon"`
	ShortForecast              string     `json:"shortForecast"`
	DetailedForecast           string     `json:"detailedForecast"`
}

//--------------------------------------------------------------------
//-- Weather Forecast Structs ForecastGridData
//--------------------------------------------------------------------

type ForecastGridData struct {
	Context    []any          `json:"@context"`
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Geometry   GridGeometry   `json:"geometry"`
	Properties GridProperties `json:"properties"`
}

type GridGeometry struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

type GridTimeValue struct {
	ValidTime TimeInterval `json:"validTime"`
	Value     float64      `json:"value"`
}

type GridUnitValue struct {
	UnitCode string  `json:"unitcode"`
	Value    float64 `json:"value"`
}

type GridWeatherParameter struct {
	UOM    string          `json:"uom,omitempty"`
	Values []GridTimeValue `json:"values"`
}

type GridProperties struct {
	ID         string `json:"@id"`
	Type       string `json:"@type"`
	UpdateTime string `json:"updateTime"`
	ValidTimes string `json:"validTimes"`

	Elevation GridUnitValue `json:"elevation"`

	ForecastOffice string `json:"forecastOffice"`
	GridID         string `json:"gridId"`
	GridX          int    `json:"gridX"`
	GridY          int    `json:"gridY"`

	Temperature                GridWeatherParameter `json:"temperature"`
	Dewpoint                   GridWeatherParameter `json:"dewpoint"`
	MaxTemperature             GridWeatherParameter `json:"maxTemperature"`
	MinTemperature             GridWeatherParameter `json:"minTemperature"`
	RelativeHumidity           GridWeatherParameter `json:"relativeHumidity"`
	ApparentTemperature        GridWeatherParameter `json:"apparentTemperature"`
	WetBlubGlobeTemperature    GridWeatherParameter `json:"wetBulbGlobeTemperature"`
	HeatIndex                  GridWeatherParameter `json:"heatIndex"`
	WindChill                  GridWeatherParameter `json:"windChill"`
	SkyCover                   GridWeatherParameter `json:"skyCover"`
	WindDirection              GridWeatherParameter `json:"windDirection"`
	WindSpeed                  GridWeatherParameter `json:"windSpeed"`
	WindGust                   GridWeatherParameter `json:"windGust"`
	ProbabilityOfPrecipitation GridWeatherParameter `json:"probabilityOfPrecipitation"`
	QuantitativePrecipitation  GridWeatherParameter `json:"quantitativePrecipitation"`

	TransportWindSpeed     GridWeatherParameter `json:"transportWindSpeed"`
	TransportWindDirection GridWeatherParameter `json:"transportWindDirection"`

	WaveHeight              GridWeatherParameter `json:"waveHeight"`
	WavePeriod              GridWeatherParameter `json:"wavePeriod"`
	WaveDirection           GridWeatherParameter `json:"waveDirection"`
	PrimarySwellHeight      GridWeatherParameter `json:"primarySwellHeight"`
	PrimarySwellDirection   GridWeatherParameter `json:"primarySwellDirection"`
	SecondarySwellHeight    GridWeatherParameter `json:"secondarySwellHeight"`
	SecondarySwellDirection GridWeatherParameter `json:"secondarySwellDirection"`
	WindWaveHeight          GridWeatherParameter `json:"windWaveHeight"`

	Pressure GridWeatherParameter `json:"pressure"`
}

//--------------------------------------------------------------------
//-- Weather Forecast Structs HourlyWeatherForecast
//-------------------------------------------------------------------

type HourlyWeatherForecast struct {
	Context    []any            `json:"@context"`
	ID         string           `json:"id"`
	Type       string           `json:"type"`
	Geometry   HourlyGeometry   `json:"geometry"`
	Properties HourlyProperties `json:"properties"`
}

type HourlyValueUnit struct {
	UnitCode string  `json:"unitCode"`
	Value    float64 `json:"value"`
}

type HourlyGeometry struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

type HourlyPeriods struct {
	Number                     int             `json:"number"`
	Name                       string          `json:"name"`
	StartTime                  time.Time       `json:"startTime"`
	EndTime                    time.Time       `json:"endTime"`
	IsDayTime                  bool            `json:"isDayTime"`
	Temperature                int             `json:"temperature"`
	TemperatureUnit            string          `json:"temperatureUnit"`
	ProbabilityOfPrecipitation HourlyValueUnit `json:"probabilityOfPrecipitation"`
	RelativeHumidity           HourlyValueUnit `json:"relativeHumidity"`
	WindSpeed                  string          `json:"windSpeed"`
	WindDirection              string          `json:"windDirection"`
	ShortForecast              string          `json:"shortForecast"`
	DetailedForecast           string          `json:"detailedForecast"`
}

type HourlyProperties struct {
	Units             string          `json:"units"`
	ForecastGenerator string          `json:"forecastGenerator"`
	GeneratedAt       time.Time       `json:"generatedAt"`
	UpdateTime        time.Time       `json:"updateTime"`
	ValidTimes        string          `json:"validTimes"`
	Elevation         HourlyValueUnit `json:"elevation"`
	Periods           []HourlyPeriods `json:"periods"`
}
