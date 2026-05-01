package database

import (
	"database/sql"
	"fmt"
	"go_surf/backend/src/api"
	"go_surf/backend/src/config"
	"go_surf/backend/src/models"
	"go_surf/backend/src/processing"
	"go_surf/backend/src/spacial"
	"log"
	"math"
)

type city struct {
	Id          int
	Latitude    float64
	Longitude   float64
	WeatherData models.SpotWeather
}

// UpdateCityWeatherSationId() should only be run when
// city weather station ids need to be updated, or when
// rebuilding the database in part or whole.
func UpdateCityWeatherStationId(db *sql.DB) error {
	rows, err := db.Query(`SELECT id, latitude, longitude FROM cities`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var cities []city

	// Scan database rows and build citites structs.
	for rows.Next() {
		var c city
		err = rows.Scan(&c.Id, &c.Latitude, &c.Longitude)
		if err != nil {
			return err
		}
		cities = append(cities, c)
	}

	m := make(map[int]string)
	for _, c := range cities {
		stationId, err := resolveStationsForCity(c)
		if err != nil {
			return fmt.Errorf("Could not resolve weather stations for cities: %w", err)
		}
		m[c.Id] = stationId

	}
	if err = insertToCitiesTable(m, db); err != nil {
		return fmt.Errorf("Error inserting weather_stations into table: %w", err)
	}
	return nil
}

func resolveStationsForCity(c city) (string, error) {
	url, err := buildNWSWeatherURL(config.NWSWeatherURL, c.Latitude, c.Longitude)
	if err != nil {
		return "", err
	}

	// Fetch weather data
	rawData, err := api.FetchURL(url)
	if err != nil {
		return "", err
	}

	c.WeatherData, err = processing.ParseSpotWeather(rawData)
	if err != nil {
		return "", err
	}

	rawData, err = api.FetchURL(c.WeatherData.Properties.ObservationStations)
	if err != nil {
		return "", err
	}
	// parse raw data
	oStations, err := processing.ParseWeatherObservationStations(rawData)
	if err != nil {
		return "", err
	}

	// find nearest city's nearest station.
	stationId := findNearestStation(c, oStations)
	return stationId, nil
}

func buildNWSWeatherURL(aString string, num1 float64, num2 float64) (string, error) {
	return fmt.Sprintf(aString, num1, num2), nil
}

func findNearestStation(city city, oStations models.ObservationStationCollection) string {
	// out of all the features coordinates, find the one closest to city coordinates.
	distance := math.MaxFloat64
	var nearestStation string
	var current float64

	for _, i := range oStations.Features {
		current = spacial.Haversine(city.Latitude, city.Longitude, i.Geometry.Coordinates[1], i.Geometry.Coordinates[0])
		if current < distance {
			distance = current
			nearestStation = i.Properties.StationIdentifier
		}
	}
	return nearestStation
}

func insertToCitiesTable(m map[int]string, db *sql.DB) error {
	if len(m) == 0 {
		return nil
	}

	query := "UPDATE cities AS c SET weather_station = v.station_id FROM (VALUES "
	args := []interface{}{}

	i := 1
	first := true
	for cityId, station := range m {
		if !first {
			query += ", "
		}
		first = false

		query += fmt.Sprintf("($%d::int, $%d)", i, i+1)
		args = append(args, cityId, station)
		i += 2
	}
	query += ") AS v(id, station_id) WHERE c.id = v.id;"

	_, err := db.Exec(query, args...)
	if err != nil {
		fmt.Println("db.Exec ERROR: ", err)
		return err
	}
	return nil
}
