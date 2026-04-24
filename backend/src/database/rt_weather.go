package database

import (
	"fmt"
	"database/sql"
	"encoding/json"
	"go_surf/backend/src/config"
	"go_surf/backend/src/api"
)

type CurrentObservation struct {
	Properties Properties `json:"properties"`
}

type Properties struct {
	Timestamp               string        `json:"timestamp"`
	Temperature             Value         `json:"temperature"`
	WindSpeed               Value         `json:"windSpeed"`
	WindDirection           Value         `json:"windDirection"`
	PrecipitationLastHour   Value         `json:"precipitationLastHour"`
	CloudLayers             []CloudLayer  `json:"cloudLayers"`
}

type Value struct {
	Value *float64 `json:"value"` // pointer allows null
}

type CloudLayer struct {
	Amount string `json:"amount"`
}

func UpdateRTWeatherData(db *sql.DB) error {
	rows, err := db.Query(`SELECT id, weather_station FROM cities;`)
	if err != nil {
		return fmt.Errorf("Error querying cities table: %w", err)
	}
	defer rows.Close()

	var id int
	var station string
	m := make(map[int]string)

	for rows.Next() {
		err := rows.Scan(&id, &station)
		if err != nil {
			return err
		}
		m[id] = station
	}

	urlMap, err := buildUrls(m) 
	if err != nil {
		return fmt.Errorf("Error buildUrls() for realtime weather: %w", err)
	}

	// fetch raw data
	rawDataMap, err := buildRawDataMap(urlMap)
	if err != nil {
		return fmt.Errorf("Error buildRawDataMap() for realtime weather: %w", err)
	}

	// parse data into structs
	idObservationMap, err := buildCurrentObservationStructs(rawDataMap)
	if err != nil {
		return fmt.Errorf("Error buildCurrentObservationStructs() for realtime weather: %w", err)
	}

	if err := updateRTWeatherTable(idObservationMap, db); err != nil {
		return fmt.Errorf("Error insertRTDataToTable() for realtime weather: %w", err)
	}	

	return nil
}

func buildUrls(m map[int]string) (map[int]string, error) {
	if len(m) == 0 {
		return m, nil
	}

	urls := make(map[int]string)
	for key, value := range m {
		url := fmt.Sprintf(config.RT_WEATHER_URL, value)
		urls[key] = url
	}
	return urls, nil
}

func buildRawDataMap(m map[int]string) (map[int][]byte, error) {
	if len(m) == 0 {
		return map[int][]byte{}, nil
	}

	rawDataMap := make(map[int][]byte)
	for key, value := range m {
		data, err := api.FetchURL(value)
		if err != nil {
			fmt.Printf("fetch failed for %d: %v\n", key, err)
			continue
		}
		
		rawDataMap[key] = data
	}
	return rawDataMap, nil
}

func buildCurrentObservationStructs(m map[int][]byte) (map[int]CurrentObservation, error) {
	if len(m) == 0 {
		return nil, nil
	}

	observationMap := make(map[int]CurrentObservation)
	for key, value := range m {
		observation, err := parseObservation(value)
		if err != nil {
			return nil, err
		}

		observationMap[key] = observation
	}
	return observationMap, nil
}

func parseObservation(data []byte) (CurrentObservation, error) {
	var observation CurrentObservation

	if err := json.Unmarshal(data, &observation); err != nil {
		return CurrentObservation{}, fmt.Errorf("Could not unmarshal CurrentObservation data: %w", err)
	}
	return observation, nil
}

func updateRTWeatherTable(m map[int]CurrentObservation, db *sql.DB) error {
	if len(m) == 0 {
		return nil
	}

	_, err := db.Exec(`TRUNCATE TABLE current_weather`)
	if err != nil {
		return err
	}

	stmnt, err := db.Prepare(`
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
	`)
	if err != nil {
		return fmt.Errorf("prepare failed: %w", err)
	}
	defer stmnt.Close()

	for CityId, observation := range m {
		props := observation.Properties

		var cloud string
		if len(props.CloudLayers) > 0 {
			cloud = props.CloudLayers[0].Amount
		}

		_, err := stmnt.Exec(
			CityId,
			props.Timestamp,                         // recorded_at
			props.WindSpeed.Value,
			props.WindDirection.Value,
			props.Temperature.Value,
			props.PrecipitationLastHour.Value,
			cloud,
			props.Timestamp, // observed_at (same for now)
		)

		if err != nil {
			return fmt.Errorf("insert failed for city %d: %w", CityId, err)
		}
	}
	return nil
}








