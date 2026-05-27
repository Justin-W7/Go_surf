package dbLib

import (
	"Go_surf_redesign/src/backend/models"
	"Go_surf_redesign/src/backend/spacial"
	"encoding/json"
	"fmt"
	"log"
	"math"
)

const (
	nwsWeatherURL = "https://api.weather.gov/points/%f,%f"
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
func updateCityWeatherStationId(c *DataClient) error {
	rows, err := c.DB.Query(`SELECT id, latitude, longitude FROM cities`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var cities []city

	// Scan database rows and build citites structs.
	for rows.Next() {
		var city city
		err = rows.Scan(&city.Id, &city.Latitude, &city.Longitude)
		if err != nil {
			return err
		}
		cities = append(cities, city)
	}

	m := make(map[int]string)
	for _, city := range cities {
		stationId, err := resolveStationsForCity(city)
		if err != nil {
			return fmt.Errorf("Could not resolve weather stations for cities: %w", err)
		}
		m[city.Id] = stationId

	}
	if err = insertToCitiesTable(m, c); err != nil {
		return fmt.Errorf("Error inserting weather_stations into table: %w", err)
	}
	return nil
}

func resolveStationsForCity(city city) (string, error) {
	url, err := buildNWSWeatherURL(nwsWeatherURL, city.Latitude, city.Longitude)
	if err != nil {
		return "", err
	}

	// Fetch weather data
	rawData, err := fetchURL(url)
	if err != nil {
		return "", err
	}

	city.WeatherData, err = parseSpotWeather(rawData)
	if err != nil {
		return "", err
	}

	rawData, err = fetchURL(city.WeatherData.Properties.ObservationStations)
	if err != nil {
		return "", err
	}
	// parse raw data
	obsvStations, err := parseWeatherObservationStations(rawData)
	if err != nil {
		return "", fmt.Errorf("could not parse weather observations station: %w", err)
	}

	// find nearest city's nearest station.
	stationId := findNearestStation(city, obsvStations)
	return stationId, nil
}

func buildNWSWeatherURL(aString string, num1 float64, num2 float64) (string, error) {
	return fmt.Sprintf(aString, num1, num2), nil
}

func findNearestStation(city city, obsvStations observationStationCollection) string {
	// out of all the features coordinates, find the one closest to city coordinates.
	distance := math.MaxFloat64
	var nearestStation string
	var current float64

	for _, f := range obsvStations.Features {
		current = spacial.Haversine(city.Latitude, city.Longitude, f.Geometry.Coordinates[1], f.Geometry.Coordinates[0])
		if current < distance {
			distance = current
			nearestStation = f.Properties.StationIdentifier
		}
	}
	return nearestStation
}

func insertToCitiesTable(m map[int]string, c *DataClient) error {
	if len(m) == 0 {
		return nil
	}

	query := "UPDATE cities AS c SET weather_station = v.station_id FROM (VALUES "
	args := []any{}

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

	_, err := c.DB.Exec(query, args...)
	if err != nil {
		fmt.Println("db.Exec ERROR: ", err)
		return err
	}
	return nil
}

type observationStationCollection struct {
	Features []Feature `json:"features"`
}

type Feature struct {
	Geometry   FeatureGeometry   `json:"geometry"`
	Properties FeatureProperties `json:"properties"`
}

type FeatureGeometry struct {
	Coordinates []float64 `json:"coordinates"` // [ longitude, latitude ]
}

type FeatureProperties struct {
	StationIdentifier string `json:"stationIdentifier"`
}

func parseWeatherObservationStations(data []byte) (observationStationCollection, error) {
	var observationStations observationStationCollection
	if err := json.Unmarshal(data, &observationStations); err != nil {
		return observationStationCollection{}, fmt.Errorf("Could not unmarshal observationStationCollection: %w", err)
	}
	return observationStations, nil
}

func parseSpotWeather(data []byte) (models.SpotWeather, error) {
	var weatherData models.SpotWeather
	if err := json.Unmarshal(data, &weatherData); err != nil {
		fmt.Println("Unmarshal error in processing.ParseSpotWeather.", err)
		return models.SpotWeather{}, err
	}
	return weatherData, nil
}
