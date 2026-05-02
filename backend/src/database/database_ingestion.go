package database

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"go_surf/backend/src/config"
	"go_surf/backend/src/models"

	//"go_surf/backend/src/processing"
	"go_surf/backend/src/spacial"
	//"go_surf/backend/src/utils"

	_ "github.com/lib/pq"
)

// ConnectDatabase establishes a connection to the PostgreSQL database "surftest".
func ConnectDatabase() *sql.DB {
	connStr := "user=justin password=password dbname=surftest sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to database successfuly.")
	fmt.Println()

	return db
}

// DisconnectDatabase disconnects the program from the local database on program shutdown.
func DisconnectDatabase(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Println("Error closing database: ", err)
		return
	}
	fmt.Println("Database disconnect succesful.")
	fmt.Println()
}

// MoveOldBuoyData moves old buoy data from active folder to cold folder.
func MoveOldBuoyData() {
	srcDir := config.Path(config.DATABASE_BUOYS_RT_RAW_DATA)
	dstDir := config.Path(config.OLD_BUOY_DATA_PATH)

	files, err := os.ReadDir(config.DATABASE_BUOYS_RT_RAW_DATA)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		srcPath := filepath.Join(srcDir, file.Name())
		dstPath := filepath.Join(dstDir, file.Name())

		os.Rename(srcPath, dstPath)
	}
}

// ClearRTData truncates (deletes) all the real time data tables within
// the database and resets their respective sequence counters.
// This should mostly be used in testing and development.
// NOTE: this function does not back up or move current data in any way.
func ClearRTData(db *sql.DB) {
	_, err := db.Exec(`TRUNCATE real_time_buoy_data_points`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`TRUNCATE current_weather`)
	if err != nil {
		log.Fatal(err)
	}

	// reset sequence counters
	_, err = db.Exec(`ALTER SEQUENCE real_time_buoy_data_points_id_seq RESTART WITH 1`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`ALTER SEQUENCE current_weather_id_seq RESTART WITH 1`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Tables cleared.")
	fmt.Println()
}

// UpdateBuoyTable reads buoy data from a CSV file (path defined in api.DATABASE_BUOYS_FILE)
// and updates the "buoys" table in the database. The function:
func UpdateBuoyTable(db *sql.DB) {
	filepath := config.Path(config.DATABASE_BUOYS_FILE)

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lineNumber := 0

	sqlStmnt, err := db.Prepare(`
		INSERT INTO buoys (id, name, latitude, longitude)
		VALUES ($1, $2, $3, $4)
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`TRUNCATE TABLE buoys RESTART IDENTITY CASCADE`)
	if err != nil {
		log.Fatal(err)
	}

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		if lineNumber == 0 {
			lineNumber++
			continue
		}

		// Convert csv fields to proper types
		id, err := strconv.Atoi(record[0])
		lat, err := strconv.ParseFloat(record[2], 64)
		long, err := strconv.ParseFloat(record[3], 64)

		_, err = sqlStmnt.Exec(id, record[1], lat, long)
		lineNumber++
	}

	defer sqlStmnt.Close()
}

// UpdateSurfSpotTable reads surf spot data from a CSV file (path defined in api.DATABASE_SURFSPOTS_FILE)
// and updates the "surfspot" table in the database.
func UpdateSurfSpotTable(db *sql.DB) {
	filepath := config.Path(config.DATABASE_SURFSPOTS_FILE)

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("UpdateSurfSpotTable failed: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lineNumber := 0

	sqlStmnt, err := db.Prepare(`
		INSERT INTO surfspot (id, name, latitude, longitude, city_id, break_type, orientation, nearest_buoy)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlStmnt.Close()

	_, err = db.Exec(`TRUNCATE TABLE surfspot RESTART IDENTITY CASCADE`)
	if err != nil {
		log.Fatal(err)
	}

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		if lineNumber == 0 {
			lineNumber++
			continue
		}

		// convert csv fields to proper types
		id, err := strconv.Atoi(record[0])
		lat, err := strconv.ParseFloat(record[2], 64)
		lon, err := strconv.ParseFloat(record[3], 64)
		city_id, err := strconv.Atoi(record[4])
		orientation, err := strconv.ParseFloat(record[6], 64)

		// find nearest buoy
		nearestBuoy := spacial.NearestBuoy(lat, lon, db)

		_, err = sqlStmnt.Exec(id, record[1], lat, lon, city_id, record[5], orientation, nearestBuoy)
		if err != nil {
			log.Fatalf("line %d: insert failed: %v", lineNumber, err)
		}
		lineNumber++
	}
}

// UpdateCitiesTable reads city data from a CSV file (path defined in api.DATABASE_CITIES_FILE)
// and updates the "cities" table in the database.
func UpdateCitiesTable(db *sql.DB) error {
	filepath := config.Path(config.DATABASE_CITIES_FILE)

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Could not open config.DATABASE_CITIES_FILE: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lineNumber := 0

	sqlStmnt, err := db.Prepare(`
		INSERT INTO cities (id, name, latitude, longitude, country, state, county)
		VALUES($1, $2, $3, $4, $5, $6, $7)
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlStmnt.Close()

	_, err = db.Exec(`TRUNCATE TABLE cities RESTART IDENTITY CASCADE`)
	if err != nil {
		log.Fatal(err)
	}

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		if lineNumber == 0 {
			lineNumber++
			continue
		}
		// split record into tokens

		// convert csv fields to proper types
		id, err := strconv.Atoi(record[0])
		lat, err := strconv.ParseFloat(record[2], 64)
		long, err := strconv.ParseFloat(record[3], 64)

		_, err = sqlStmnt.Exec(id, record[1], lat, long, record[4], record[5], record[6])
		if err != nil {
			log.Fatalf("line %d: insert failed: %v", lineNumber, err)
		}
		lineNumber++
	}

	if err := UpdateCityWeatherStationId(db); err != nil {
		return fmt.Errorf("updateCityWeatherStationId failed: %w", err)
	}
	return nil
}

// UpdateRTBuoyDataTable updates the real-time buoy data table in the database.
// The funcion does the following:
//  1. Clears table to new data.
//  2. Reads all the raw data files from the directory api.DATABASE_BUOYS_RAW_DATA.
//  3. Iterates through each file, gets the buoy id from the file name, and calls
//     processRTBuoyFile to insert data into the database.
func UpdateRTBuoyDataTable(db *sql.DB) error {
	// Clear table for new data.
	_, err := db.Exec(`TRUNCATE real_time_buoy_data_points`)
	if err != nil {
		return err
	}

	folder := config.Path(config.DATABASE_BUOYS_RT_RAW_DATA)

	files, err := os.ReadDir(folder)
	if err != nil {
		return err
	}

	for _, file := range files {
		name := file.Name()
		start := len(name) - 9
		buoystr := name[start : start+5]

		// convert buoystr to int
		buoyID, err := strconv.Atoi(buoystr)
		if err != nil {
			return err
		}

		filepath := filepath.Join(folder, name)

		processRTBuoyFile(db, filepath, buoyID)
	}

	return nil
}

func processRTBuoyFile(db *sql.DB, filepath string, buoyID int) {
	// open file
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// scan line with data
	scanner := bufio.NewScanner(f)

	year := strconv.Itoa(time.Now().Year())
	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, year) {
			continue
		}

		buoyDataPoint, err := parseRTBuoyLine(line, buoyID)
		if err != nil {
			log.Fatal(err)
		}

		err = insertBuoyData(db, buoyDataPoint)
		if err != nil {
			log.Fatal(err)
		}
		break
	}
}

// parseRTBuoyLine returns a pointer to a models.BuoyDataPoint and an error.
func parseRTBuoyLine(line string, buoyID int) (*models.BuoyDataPoint, error) {
	data := strings.Fields(line)

	// parse time
	timeLayout := "2006 01 02 15 04"

	recordedAt, err := time.Parse(timeLayout, strings.Join(data[:5], " "))
	if err != nil {
		log.Fatal(err)
	}

	// parse data types
	windDirection, _ := parseDataFloat(data[5])
	windSpeed, _ := parseDataFloat(data[6])
	windGust, _ := parseDataFloat(data[7])
	waveHeightM, _ := parseDataFloat(data[8])
	dominantWavePeriod, _ := parseDataFloat(data[9])
	avgWavePeriod, _ := parseDataFloat(data[10])
	meanWaveDirection, _ := parseDataFloat(data[11])
	airTemperature, _ := parseDataFloat(data[13])
	waterTemperature, _ := parseDataFloat(data[14])

	// build models.BuoyDataPoint struct
	p := &models.BuoyDataPoint{
		BuoyID:                buoyID,
		RecordedAt:            recordedAt,
		WindDirectionDegT:     windDirection,
		WindSpeedMetersPerSec: windSpeed,
		WindGustMetersPerSec:  windGust,
		WaveHeightM:           waveHeightM,
		DominantWavePeriodSec: dominantWavePeriod,
		AvgWavePeriodSec:      avgWavePeriod,
		MeanWaveDirectionDegT: meanWaveDirection,
		AirTempDegC:           airTemperature,
		WaterTempDegC:         waterTemperature,
		InsertedAt:            time.Now().UTC(),
	}
	return p, nil
}

func insertBuoyData(db *sql.DB, p *models.BuoyDataPoint) error {
	_, err := db.Exec(`
			INSERT INTO real_time_buoy_data_points (
				buoy_id,
				recorded_at,
				windDir_degt,
				windSpeed_m_pers,
				windGust_m_pers,
				waveH_m,
				domWP_sec,
				avgWaveP_sec,
				meanWaveDir_degt,
				airT_degc,
				waterT_degc,
				inserted_at
				)
			Values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`,
		p.BuoyID,
		p.RecordedAt,
		p.WindDirectionDegT,
		p.WindSpeedMetersPerSec,
		p.WindGustMetersPerSec,
		p.WaveHeightM,
		p.DominantWavePeriodSec,
		p.AvgWavePeriodSec,
		p.MeanWaveDirectionDegT,
		p.AirTempDegC,
		p.WaterTempDegC,
		p.InsertedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

// THIS HAS BEEN REPLACED WITH UpdateRTWeatherData() located in rt_weather.go
//
/*
func UpdateRTWeatherTable(db *sql.DB) error {
	// Clear table for new data.
	_, err := db.Exec(`TRUNCATE current_weather`)
	if err != nil {
		return err
	}

	// for each record in cities table get latitude and longitude.
	rows, err := db.Query(`SELECT id, latitude, longitude FROM cities`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var lat float64
		var lon float64

		err = rows.Scan(&id, &lat, &lon)
		if err != nil {
			return err
		}

		// build url for fetching weather data.
		url := fmt.Sprintf(config.NWSWeatherURL, lat, lon)
		fmt.Println(url);

		// get weather forcast for lat lon.
		data, err := api.FetchWeatherForecast(url)
		if err != nil {
			return fmt.Errorf("FetchWeatherForecast: api call failed: %w", err)
		}

		// parse spot weather for hourly forecast url
		forecast, err := processing.ParseSpotWeather(data)
		if err != nil {
			return fmt.Errorf("ParseSpotWeather: failed to parse weather forecast data: %w", err)
		}

		hourlyUrl := forecast.Properties.ForecastHourly

		// get hourly forecast
		rawData, err := api.FetchHourlyWeatherForecast(hourlyUrl)
		if err != nil {
			return fmt.Errorf("Error 1: %w", err)
		}

		// parse raw hourly weather data
		hourlyForecast, err := processing.ParseHourlyWeatherForecast(rawData)
		if err != nil {
			return fmt.Errorf("Error 2: %w", err)
		}

		// build models.WeatherDataPoint
		dataPoint, err := parseRTWeatherData(id, &hourlyForecast)
		if err != nil {
			return fmt.Errorf("Error 3: %w", err)
		}

		// insert weather data into database
		err = insertRTWeatherData(db, dataPoint)
		if err != nil {
			return fmt.Errorf("Error 4: %w", err)
		}
	}
	return nil
}


func parseRTWeatherData(id int, data *models.HourlyWeatherForecast) (*models.WeatherDatapoint, error) {
	forecast := *data

	t := time.Now().UTC()

	st := forecast.Properties.Periods[0].StartTime

	startTime, err := time.Parse(time.RFC3339, st)
	if err != nil {
		log.Fatal(err)
	}
	utcStartTime := startTime.UTC()
	observedAt := utcStartTime
	recordedAt := t
	windSpeed := forecast.Properties.Periods[0].WindSpeed
	windDir := forecast.Properties.Periods[0].WindDirection
	airTempC := utils.FarenheitToCelsius(float64(forecast.Properties.Periods[0].Temperature))
	precipitation := forecast.Properties.Periods[0].ProbabilityOfPrecipitation.Value
	cloudCoverage := forecast.Properties.Periods[0].ShortForecast

	// parse into struct to be passed to an insert function
	p := &models.WeatherDatapoint{
		CityID:        id,
		ObservedAt:    observedAt,
		RecordedAt:    recordedAt,
		WindSpeed:     &windSpeed,
		WindDirection: &windDir,
		AirTemp:       &airTempC,
		Precipitation: &precipitation,
		CloudCoverage: &cloudCoverage,
	}
	return p, nil
}

func insertRTWeatherData(db *sql.DB, p *models.WeatherDatapoint) error {
	_, err := db.Exec(`
			INSERT INTO current_weather (
				city_id,
				recorded_at,
				wind_speed,
				wind_direction,
				air_temp_c,
				precipitation,
				cloud_coverage,
				observed_at
				)
			Values ($1, $2, $3, $4, $5, $6, $7, $8)
		`,
		p.CityID,
		p.RecordedAt,
		p.WindSpeed,
		p.WindDirection,
		p.AirTemp,
		p.Precipitation,
		p.CloudCoverage,
		p.ObservedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
*/

func UpdateCurrentSurfConditions(db *sql.DB) error {
	// Clear table for new data
	_, err := db.Exec(`TRUNCATE current_surf_spot_conditions`)
	if err != nil {
		return err
	}

	// Get all cities
	cities, err := db.Query(`SELECT
				id,
				name,
				latitude,
				longitude
			FROM cities`)
	if err != nil {
		return err
	}
	defer cities.Close()

	// Get all weather data
	weather, err := db.Query(`SELECT
				city_id,
				recorded_at,
				wind_speed,
				wind_direction,
				air_temp,
				precipitation,
				cloud_coverage,
				observed_at
			FROM current_weather`)
	if err != nil {
		return err
	}
	defer weather.Close()

	// Get all buoys
	buoys, err := db.Query(`SELECT
				id,
				latitude,
				longitude
			FROM buoys`)
	if err != nil {
		return err
	}
	defer buoys.Close()

	// Get data from all buoys
	buoyData, err := db.Query(`SELECT
				buoy_id,
				recorded_at,
				waveh_m,
				domwp_sec,
				avgwavep_sec,
				meanwavedir_degt,
				watert_degc
			FROM real_time_buoy_data_points`)
	if err != nil {
		return err
	}
	defer buoyData.Close()

	// Get all surfspots
	spots, err := db.Query(`SELECT
				id,
				name,
				latitude,
				longitude,
				city_id,
				nearest_buoy
			FROM surfspot`)
	if err != nil {
		return err
	}
	defer spots.Close()

	// build maps for each dataset
	// build city map
	cityMap := make(map[int]models.City)
	for cities.Next() {
		var c models.City

		err := cities.Scan(
			&c.ID,
			&c.Name,
			&c.Latitude,
			&c.Longitude,
		)
		if err != nil {
			return err
		}
		cityMap[c.ID] = c
	}

	// build weather map
	weatherMap := make(map[int]models.WeatherDatapoint)
	for weather.Next() {
		var w models.WeatherDatapoint

		err := weather.Scan(
			&w.CityID,
			&w.RecordedAt,
			&w.WindSpeed,
			&w.WindDirection,
			&w.AirTemp,
			&w.Precipitation,
			&w.CloudCoverage,
			&w.ObservedAt,
		)
		if err != nil {
			return err
		}
		weatherMap[w.CityID] = w
	}

	// build buoy map
	buoyMap := make(map[int]models.Buoy)
	for buoys.Next() {
		var b models.Buoy

		err := buoys.Scan(
			&b.ID,
			&b.Latitude,
			&b.Longitude,
		)
		if err != nil {
			return err
		}
		buoyMap[b.ID] = b
	}

	// build buoyData map
	buoyDataMap := make(map[int]models.BuoyData)
	for buoyData.Next() {
		var bd models.BuoyData

		err := buoyData.Scan(
			&bd.BuoyID,
			&bd.RecordedAt,
			&bd.WaveHeightM,
			&bd.DominantWavePeriodSec,
			&bd.AvgWavePeriodSec,
			&bd.MeanWaveDirectionDegT,
			&bd.WaterTempDegC,
		)
		if err != nil {
			return err
		}
		buoyDataMap[bd.BuoyID] = bd
	}

	// build surfspot map
	surfSpotsMap := make(map[int]models.StaticSurfSpot)
	for spots.Next() {
		var s models.StaticSurfSpot

		err := spots.Scan(
			&s.ID,
			&s.Name,
			&s.Latitude,
			&s.Longitude,
			&s.CityID,
			&s.NearestBuoy,
		)
		if err != nil {
			return err
		}
		surfSpotsMap[s.ID] = s
	}

	// build CurrentSurfSpotConditions struct
	// buildCurrentSurfConditions()
	conditions := buildCurrentSurfConditions(surfSpotsMap, buoyMap, buoyDataMap, weatherMap)

	// insert Conditions into table
	for _, p := range conditions {
		err = insertCurrentSurfConditions(p, db)
	}
	if err != nil {
		return err
	}

	return nil
}

func buildCurrentSurfConditions(
	surfSpots map[int]models.StaticSurfSpot, buoys map[int]models.Buoy,
	buoyData map[int]models.BuoyData, weather map[int]models.WeatherDatapoint,
) []*models.CurrentSurfSpotConditions {
	// slice to hold *models.CurrentSurfSpotConditions
	var conditions []*models.CurrentSurfSpotConditions

	for _, spot := range surfSpots {
		buoy := buoys[spot.NearestBuoy]
		bData, ok := buoyData[buoy.ID]
		if !ok {
			continue
		}
		wData := weather[spot.CityID]

		// build models.CurrentSurfSpotConditions
		p := &models.CurrentSurfSpotConditions{
			SpotId:                spot.ID,
			RecordedAt:            bData.RecordedAt,
			DomSwellHeightM:       bData.WaveHeightM,
			DomSwellDir:           bData.MeanWaveDirectionDegT,
			WindSpeedMph:          wData.WindSpeed,
			WindDirection:         wData.WindDirection,
			AirTempDegF:           wData.AirTemp,
			WaterTempDegC:         bData.WaterTempDegC,
			Precipitation:         wData.Precipitation,
			CloudCoverage:         wData.CloudCoverage,
			DominantWavePeriodSec: bData.DominantWavePeriodSec,
		}

		conditions = append(conditions, p)
	}

	return conditions
}

func insertCurrentSurfConditions(p *models.CurrentSurfSpotConditions, db *sql.DB) error {
	_, err := db.Exec(`
			INSERT INTO current_surf_spot_conditions (
			spot_id,
			recorded_at,
			dom_swell_height_m,
			dom_swell_dir,
			wind_speed_mph,
			wind_direction,
			air_temp_deg_f,
			water_temp_deg_c,
			precipitation,
			cloud_coverage,
			domwp_sec
			)
		Values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`,
		p.SpotId,
		p.RecordedAt,
		p.DomSwellHeightM,
		p.DomSwellDir,
		p.WindSpeedMph,
		p.WindDirection,
		p.AirTempDegF,
		p.WaterTempDegC,
		p.Precipitation,
		p.CloudCoverage,
		p.DominantWavePeriodSec,
	)
	if err != nil {
		return err
	}

	return nil
}

// UTILITY FUNCTIONS
// parseDataFloat returns a pointer so it can return multiple states.
// This enables it to return a number, a missing value ("MM") or nil to mean no value, instead of 0
// which may be a valid value for the data.
func parseDataFloat(value string) (*float64, error) {
	if value == "MM" {
		return nil, nil
	}

	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Fatal(err)
	}
	return &result, nil
}
