package dbLib

import (
	meteo "Go_surf_redesign/src/backend/api"
	"Go_surf_redesign/src/backend/data"
	"Go_surf_redesign/src/backend/spacial"
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

const (
	conntStr   = "user=postgres password=pass dbname=surfdb sslmode=disable"
	psqlDriver = "postgres"

	dbBuoysList    = "buoys.csv"
	dbCitiesList   = "cities.csv"
	dbSurfSpotList = "surfspots.csv"

	rtBuoyDataURL = "https://www.ndbc.noaa.gov/data/realtime2/%s.txt"
)

type DataClient struct {
	DB *sql.DB
}

func NewDBClient() *DataClient {
	DB, err := sql.Open(psqlDriver, conntStr)
	if err != nil {
		log.Fatalf("Failed to instantiate database client: %v", err)
		return nil
	}
	c := &DataClient{
		DB: DB,
	}
	return c
}

func (c *DataClient) PingDB() error {
	if err := c.DB.Ping(); err != nil {
		return err
	}
	return nil
}

func (c *DataClient) Close() error {
	fmt.Println("Database disconnected.")
	return c.DB.Close()
}

func (c *DataClient) LoadStaticData() {

	err := c.UpdateStaticCitiesTable()
	if err != nil {
		log.Println("Error: ", err)
	}

	err = c.UpdateStaticBuoyTable()
	if err != nil {
		log.Println("Error: ", err)
	}

	err = c.UpdateStaticSurfSpotTable()
	if err != nil {
		log.Println("Error: ", err)
	}
}

func (c *DataClient) UpdateStaticBuoyTable() error {
	sqlStmnt, err := c.DB.Prepare(`
		INSERT INTO buoys (id, name, latitude, longitude)
		VALUES ($1, $2, $3, $4)
	`)
	if err != nil {
		return fmt.Errorf("Error preparing sql statement in UpdateStaticBuoyTable(): %w", err)
	}
	defer sqlStmnt.Close()

	_, err = c.DB.Exec(`TRUNCATE TABLE buoys RESTART IDENTITY CASCADE`)
	if err != nil {
		return fmt.Errorf("Error truncating static buoys table in UpdateStaticBuoyTable(): %w", err)
	}

	file, err := os.Open(filepath.Join(data.FilePathBuilder(), dbBuoysList))
	if err != nil {
		return fmt.Errorf("could not open dbBuoysList in UpdateStaticBuoyTable(): %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	linenumber := 0
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("Error reading csv reader in UpdateStaticBuoyTable(): %w", err)

		}
		if linenumber == 0 {
			linenumber++
			continue
		}

		id, err := strconv.Atoi(record[0])
		if err != nil {
			return fmt.Errorf("invalid id on linenumber %d: %w", linenumber, err)
		}

		lat, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return fmt.Errorf("invalid latitude on linenumber %d: %w", linenumber, err)
		}

		lon, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return fmt.Errorf("invalid longitude on linenumber %d: %w", linenumber, err)
		}

		_, err = sqlStmnt.Exec(id, record[1], lat, lon)
		if err != nil {
			return fmt.Errorf("Could not update record %d: %w", id, err)
		}
		linenumber++
	}
	return nil
}

func (c *DataClient) UpdateStaticCitiesTable() error {
	sqlStmnt, err := c.DB.Prepare(`
			INSERT INTO cities (id, name, latitude, longitude, country, state, county)
			VALUES($1, $2, $3, $4, $5, $6, $7)
		`)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlStmnt.Close()

	_, err = c.DB.Exec(`TRUNCATE TABLE cities RESTART IDENTITY CASCADE`)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(filepath.Join(data.FilePathBuilder(), dbCitiesList))
	if err != nil {
		return fmt.Errorf("Error opening file in UpdaateStaticCitiesTable(): %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	linenumber := 0

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("Error reading file ")
		}
		if linenumber == 0 {
			linenumber++
			continue
		}
		id, err := strconv.Atoi(record[0])
		lat, err := strconv.ParseFloat(record[2], 64)
		lon, err := strconv.ParseFloat(record[3], 64)

		_, err = sqlStmnt.Exec(id, record[1], lat, lon, record[4], record[5], record[6])
		if err != nil {
			return fmt.Errorf("Could not execute sql statement in UpdateStaticCitiesTable(): %w", err)
		}
		linenumber++
	}

	if err := updateCityWeatherStationId(c); err != nil {
		return fmt.Errorf("Error updating weatherstation id in UpdateStaticCitiesTable(): %w", err)
	}
	return nil
}

func (c *DataClient) UpdateStaticSurfSpotTable() error {
	sqlStmnt, err := c.DB.Prepare(`
			INSERT INTO surfspot (id, name, latitude, longitude, city_id, break_type, orientation, nearest_buoy)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8)
		`)
	if err != nil {
		return err
	}
	defer sqlStmnt.Close()

	_, err = c.DB.Exec(`TRUNCATE TABLE surfspot RESTART IDENTITY CASCADE`)
	if err != nil {
		return err
	}

	file, err := os.Open(filepath.Join(data.FilePathBuilder(), dbSurfSpotList))
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(file)
	linenumber := 0

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if linenumber == 0 {
			linenumber++
			continue
		}

		id, err := strconv.Atoi(record[0])
		lat, err := strconv.ParseFloat(record[2], 64)
		lon, err := strconv.ParseFloat(record[3], 64)
		city_id, err := strconv.Atoi(record[4])
		orientation, err := strconv.ParseFloat(record[6], 64)

		nearestBuoy := spacial.NearestBuoy(lat, lon, c.DB)

		_, err = sqlStmnt.Exec(id, record[1], lat, lon, city_id, record[5], orientation, nearestBuoy)
		if err != nil {
			return fmt.Errorf("line %d: insert failed: %w", linenumber, err)
		}
		linenumber++
	}
	return nil
}

func (c *DataClient) UpdateRTBuoyData(ctx context.Context, api *meteo.Client) error {
	// read bouy ids from static buoy table
	ids, err := c.GetBuoyIds()
	if err != nil {
		return fmt.Errorf("could net get buoy ids: %w", err)
	}

	// Truncate table to clear for new data.
	_, err = c.DB.Exec(`TRUNCATE real_time_buoy_data_points`)
	if err != nil {
		return err
	}

	// for each bouy id, fetch buoy data
	idsMap := make(map[int]*meteo.BouyObservation)
	for _, id := range ids {
		obs, err := api.RTBouy.GetObservation(ctx, strconv.Itoa(id))
		if err != nil {
			continue
		}
		idsMap[id] = obs
	}

	// read new data into DB. Use a seperate helper function.
	for buoyId, obs := range idsMap {
		if err := c.insertRTBouyData(strconv.Itoa(buoyId), obs); err != nil {
			fmt.Println("could not insert buoy data: ", err)
			continue
		}
	}
	fmt.Println("Realtime buoy data updated.")
	return nil
}

func (c *DataClient) insertRTBouyData(buoyId string, obs *meteo.BouyObservation) error {
	sqlStmnt, err := c.DB.Prepare(`
		INSERT INTO real_time_buoy_data_points (
			buoy_id,
			recorded_at,
			winddir_degt,
			windspeed_m_pers,
			windgust_m_pers,
			waveh_m,
			domwp_sec,
			avgwavep_sec,
			meanwavedir_degt,
			airt_degc,
			watert_degc,
			inserted_at
		)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`)
	if err != nil {
		return fmt.Errorf("could not prepare statement: %w", err)
	}
	defer sqlStmnt.Close()

	_, err = sqlStmnt.Exec(
		buoyId,
		obs.RecordedAt,
		obs.WindDirectionDegT,
		obs.WindSpeedMetersPerSec,
		obs.WindGustMetersPerSec,
		obs.WaveHeightM,
		obs.DominantWavePeriodSec,
		obs.AvgWavePeriodSec,
		obs.MeanWaveDirectionDegT,
		obs.AirTempDegC,
		obs.WaterTempDegC,
		obs.InsertedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

// GetBuoyData returns all the static buoy table ids in a slice.
func (c *DataClient) GetBuoyIds() ([]int, error) {
	rows, err := c.DB.Query(`SELECT id FROM buoys`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (c *DataClient) UpdateRTWeatherData(ctx context.Context, api *meteo.Client) {
	// iterate through each city for weather station
	weatherStations, err := c.GetWeatherStations()
	if err != nil {
		fmt.Printf("could not get weather stations: %v", err)
	}
	// clear rt weather data from table
	_, err = c.DB.Exec(`TRUNCATE current_weather`)
	if err != nil {
		fmt.Printf("could not truncate current_weather: %v", err)
	}
	// for each station get data from api &
	// insert into table
	for _, weatherStations := range weatherStations {
		obs, err := api.RTWeather.GetObservation(ctx, weatherStations.station)
		if err = c.insertRTWeatherData(weatherStations.cityId, obs); err != nil {
			fmt.Printf("could not insert current weather observations into table: %v", err)
			continue
		}
	}
	fmt.Println("Realtime weather data updated.")
}

func (c *DataClient) insertRTWeatherData(cityId int, obs *meteo.WeatherObservation) error {
	sqlStmnt, err := c.DB.Prepare(`
		INSERT INTO current_weather(
			city_id,
			recorded_at,
			wind_speed,
			wind_direction,
			air_temp_c,
			precipitation,
			cloud_coverage,
			observed_at
		)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)
	`)
	if err != nil {
		return fmt.Errorf("could not prepare statement: %w", err)
	}

	// Check for empty values.
	var cloudLayersAmount string
	if len(obs.Properties.CloudLayers) > 0 {
		cloudLayersAmount = obs.Properties.CloudLayers[0].Amount
	} else {
		cloudLayersAmount = "unknown"
	}

	// Formatting data
	data := obs.Properties.WindSpeed.Value
	strWindSpeed := "NA"
	if data != nil {
		windSpeedKMH := *obs.Properties.WindSpeed.Value
		intWindSpeedKMPH := int(math.Round(windSpeedKMH))
		fWindSpeedMPH := KMHToMPH(float64(intWindSpeedKMPH))
		strWindSpeed = strconv.Itoa(int(fWindSpeedMPH))
	}

	_, err = sqlStmnt.Exec(
		cityId,
		obs.RecordedAt,
		strWindSpeed,
		obs.Properties.WindDirection.Value,
		obs.Properties.Temperature.Value,
		obs.Properties.Precipitation.Value,
		cloudLayersAmount,
		obs.Properties.Timestamp,
	)
	if err != nil {
		return err
	}
	return nil
}

type cityWeatherStation struct {
	cityId  int
	station string
}

func (c *DataClient) GetWeatherStations() ([]cityWeatherStation, error) {
	rows, err := c.DB.Query(`SELECT id, weather_station FROM cities`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var weatherStations []cityWeatherStation
	for rows.Next() {
		var id int
		var station string
		if err := rows.Scan(&id, &station); err != nil {
			return nil, err
		}
		cityWeatherStation := cityWeatherStation{id, station}
		weatherStations = append(weatherStations, cityWeatherStation)
	}
	return weatherStations, nil
}

type CurrentSurfSpotConditions struct {
	SpotId                int
	RecordedAt            time.Time
	DomSwellHeightM       *float64 // from buoy data
	DomSwellDir           *float64 // from buoy data
	WindSpeedMph          *string  // from city weather data
	WindDirection         *string  // from city weather data
	AirTempDegC           *float64
	WaterTempDegC         *float64 // from buoy data
	Precipitation         *float64 // from city weather data
	CloudCoverage         *string  // from city weather data
	DominantWavePeriodSec *float64 // from buoy data
	BuoyId                int
}

func (c *DataClient) UpdateCurrentSurfConditions(api *meteo.Client) {
	// Iterate through surfspots
	surfSpots, err := c.GetSurfSpots()
	if err != nil {
		fmt.Printf("could not get surfspot ids: %v", err)
	}
	// clear current conditions data from table
	_, err = c.DB.Exec(`TRUNCATE current_surf_spot_conditions`)
	if err != nil {
		fmt.Printf("could not truncate current_surf_spot_conditions: %v", err)
	}

	// Build []CuurentSurfSpotConditions from surfSpots.
	conditions, err := c.buildCurrentConditions(surfSpots)
	if err != nil {
		fmt.Printf("could not build current surf conditions: %v", err)
	}

	// For each data set of conditions, insert into database
	for _, data := range conditions {
		if err = c.insertCurrentSurfConditions(data); err != nil {
			fmt.Printf("could not insert surf conditions data into table: %v", err)
		}
	}
	fmt.Println("Current surf conditions updated.")
}

func (c *DataClient) insertCurrentSurfConditions(data CurrentSurfSpotConditions) error {
	sqlStmnt, err := c.DB.Prepare(`
		INSERT INTO current_surf_spot_conditions (
		spot_id,
		recorded_at,
		dom_swell_height_m,
		dom_swell_dir,
		wind_speed_mph,
		wind_direction,
		air_temp_deg_c,
		water_temp_deg_c,
		precipitation,
		cloud_coverage,
		domwp_sec,
		nearest_buoy
		)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`)
	if err != nil {
		return fmt.Errorf("could not prepare statment %w", err)
	}
	defer sqlStmnt.Close()

	_, err = sqlStmnt.Exec(
		data.SpotId,
		data.RecordedAt,
		data.DomSwellHeightM,
		data.DomSwellDir,
		data.WindSpeedMph,
		data.WindDirection,
		data.AirTempDegC,
		data.WaterTempDegC,
		data.Precipitation,
		data.CloudCoverage,
		data.DominantWavePeriodSec,
		data.BuoyId,
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *DataClient) buildCurrentConditions(surfSpots []surfSpot) ([]CurrentSurfSpotConditions, error) {
	var conditionsSlice []CurrentSurfSpotConditions
	var conditions CurrentSurfSpotConditions

	for _, surfSpot := range surfSpots {
		// for each surfspot, get all the correlating data to build a current surf spot conditions struct.
		conditions.SpotId = surfSpot.ID
		conditions.BuoyId = surfSpot.NearestBuoy
		buoyQuery := `
	 		SELECT recorded_at, waveh_m, meanwavedir_degt, domwp_sec, watert_degc
			FROM real_time_buoy_data_points
			WHERE buoy_id = $1
		`
		// use surfSpot.NearestBuoy to get current buoy data.
		buoyData, err := c.DB.Query(buoyQuery, surfSpot.NearestBuoy)
		if err != nil {
			return nil, err
		}
		defer buoyData.Close()

		var recordedAt time.Time
		var domSwellHeightM *float64
		var domSwellDir *float64
		var domWavePeriod *float64
		var waterTemp *float64

		for buoyData.Next() {
			if err := buoyData.Scan(
				&recordedAt,
				&domSwellHeightM,
				&domSwellDir,
				&domWavePeriod,
				&waterTemp,
			); err != nil {
				fmt.Printf("could not scan rows for buoyData while building CurrentSurfSpotConditions: %v", err)
				continue
			}

			conditions.RecordedAt = recordedAt.UTC()
			conditions.DomSwellHeightM = domSwellHeightM
			conditions.DomSwellDir = domSwellDir
			conditions.DominantWavePeriodSec = domWavePeriod
			conditions.WaterTempDegC = waterTemp
		}

		weatherQuery := `
			SELECT wind_speed, wind_direction, air_temp_c, precipitation, cloud_coverage
			FROM current_weather
			WHERE city_id = $1
		`

		// use surfSpot.CityId to get current weather data
		weatherData, err := c.DB.Query(weatherQuery, surfSpot.CityId)
		if err != nil {
			return nil, err
		}
		defer weatherData.Close()

		var windSpeed *string
		var windDir *string
		var airTemp *float64
		var precipitation *float64
		var cloudCoverage *string

		for weatherData.Next() {
			if err := weatherData.Scan(
				&windSpeed,
				&windDir,
				&airTemp,
				&precipitation,
				&cloudCoverage,
			); err != nil {
				fmt.Printf("could not scan rows for weatherData while building CurrentSurfSpotConditions: %v", err)
				continue
			}
			conditions.WindSpeedMph = windSpeed
			conditions.WindDirection = windDir
			conditions.AirTempDegC = airTemp
			conditions.Precipitation = precipitation
			conditions.CloudCoverage = cloudCoverage
		}
		conditionsSlice = append(conditionsSlice, conditions)
	}
	return conditionsSlice, nil
}

type surfSpot struct {
	ID          int
	Name        string
	CityId      int
	NearestBuoy int
}

func (c *DataClient) GetSurfSpots() ([]surfSpot, error) {
	rows, err := c.DB.Query(`SELECT id, name, city_id, nearest_buoy FROM surfspot`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var surfSpots []surfSpot
	var id int
	var name string
	var cityId int
	var buoy int

	for rows.Next() {
		if err := rows.Scan(&id, &name, &cityId, &buoy); err != nil {
			return nil, err
		}
		surfSpot := surfSpot{id, name, cityId, buoy}
		surfSpots = append(surfSpots, surfSpot)
	}
	return surfSpots, nil
}

// UTILITY FUNCITONS

// These functions may be used in other files within db_lib
func fetchURL(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed request to %s: %w", url, err)
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response from %s: %w", url, err)
	}
	return data, nil
}

func KMHToMPH(kmh float64) float64 {
	return kmh * 0.621371
}

/*
 TODOS:
// Need to update forecasted data
func UpdateForecastedBuoyData()    {}
func UpdateForecastedWeatherData() {}
func UpdateForecastedConditions()  {}
*/
