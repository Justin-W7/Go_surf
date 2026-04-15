package config

// URLs
const (
	SpitcastSpotURL     = "https://api.spitcast.com/api/spot"
	SpitcastForecastURL = "https://api.spitcast.com/api/spot_forecast/%d/%d/%d/%d"
	NWSWeatherURL       = "https://api.weather.gov/points/%f,%f"

	// Real Time Buoy Data
	// EXAMPLE: https://www.ndbc.noaa.gov/data/realtime2/4038.txt
	NDBCBouyDataURL = "https://www.ndbc.noaa.gov/data/realtime2/%s.txt"
	STATION_ID_FILE = "backend/src/api/station_lists/ndbcstations_CA.txt"
)

// Database files
const (
	DATABASE_BUOYS_FILE        = "backend/src/database/buoys.csv"
	DATABASE_CITIES_FILE       = "backend/src/database/cities.csv"
	DATABASE_SURFSPOTS_FILE    = "backend/src/database/surfspots.csv"
	DATABASE_BUOYS_RT_RAW_DATA = "backend/src/database/raw_data/NDBC_buoy_data"
	OLD_BUOY_DATA_PATH         = "backend/src/database/raw_data/old_NDBC_buoy_data"
)
