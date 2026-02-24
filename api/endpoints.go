package api

const (
	SpitcastSpotURL     = "https://api.spitcast.com/api/spot"
	SpitcastForecastURL = "https://api.spitcast.com/api/spot_forecast/%d/%d/%d/%d"
	NWSWeatherURL       = "https://api.weather.gov/points/%f,%f"

	// Real Time Buoy Data
	// EXAMPLE: https://www.ndbc.noaa.gov/data/realtime2/4038.txt
	NDBCBouyDataURL = "https://www.ndbc.noaa.gov/data/realtime2/%s.txt"

	STATION_ID_FILE = "api/station_lists/ndbcstations_CA.txt"
)
