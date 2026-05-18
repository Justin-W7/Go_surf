package dbLib

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"go_surf/backend/src/database"
	"go_surf/backend/src/spacial"
	"io"
	"log"
	"os"
	"strconv"
)

const (
	constStr   = "user=postgres password=pass dbname=surfdate sslmode=disable"
	psqlDriver = "postgers"

	dbBouysList    = "backend/src/database/buoys.csv"
	dbCitiesList   = "backend/src/database/cities.csv"
	dbSurfSpotList = "backend/src/database/surfspots.csv"
)

type DataClient struct {
	db *sql.DB
}

func NewDBClient() (*DataClient, error) {
	db, err := sql.Open(psqlDriver, constStr)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	c := &DataClient{
		db: db,
	}
	return c, nil
}

func (c *DataClient) Close() error {
	return c.db.Close()
}

/*

// Need to build static tables
func UpdateStaticBuoyTable()     {}
func UpdateStaticSurfSpotTable() {}
func UpdateStaticCitiesTable()   {}

// Need to update real time data tables
func UpdateRTBuoyData()            {}
func UpdateRTWeatherData()         {}
func UpdateCurrentSurfConditions() {}

// Need to update forecasted data
func UpdateForecastedBuoyData()    {}
func UpdateForecastedWeatherData() {}
func UpdateForecastedConditions()  {}
*/

func (c *DataClient) UpdateStaticBouyTable() error {
	sqlStmnt, err := c.db.Prepare(`
		INSERT INTO bouys (id, name, latitude, longitude)
		VALUES ($1, $2, $3, $4)
	`)
	if err != nil {
		return fmt.Errorf("Error preparing sql statement in UpdateStaticBouyTable(): %w", err)
	}
	defer sqlStmnt.Close()

	_, err = c.db.Exec(`TRUNCATE TABLE bouys RESTART IDENTITY CASCADE`)
	if err != nil {
		return fmt.Errorf("Error truncating static bouys table in UpdateStaticBouyTable(): %w", err)
	}

	file, err := os.Open(dbBouysList)
	if err != nil {
		return fmt.Errorf("Error opening dbBouysList in UpdateStaticBouyTable(): %w", err)
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
			return fmt.Errorf("Error reading csv reader in UpdateStaticBouyTable(): %w", err)
		}
		if linenumber == 0 {
			linenumber++
			continue
		}

		id, err := strconv.Atoi(record[0])
		lat, err := strconv.ParseFloat(record[2], 64)
		lon, err := strconv.ParseFloat(record[3], 64)

		_, err = sqlStmnt.Exec(id, record[1], lat, lon)
		linenumber++
	}
	return nil
}

func (c *DataClient) UpdateStaticCitiesTable() error {
	sqlStmnt, err := c.db.Prepare(`
			INSERT INTO cities (id, name, latitude, longitude, country, state, county)
			VALUES($1, $2, $3, $4, $5, $6, $7)
		`)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlStmnt.Close()

	_, err = c.db.Exec(`TRUNCATE TABLE cities RESTART IDENTITY CASCADE`)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(dbCitiesList)
	if err != nil {
		return fmt.Errorf("Error opening file in UpdaateStaticCitiesTable(): %w", err)
	}
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
		lon, err := strconv.ParseFloat(record[2], 64)

		_, err = sqlStmnt.Exec(id, record[1], lat, lon, record[4], record[5], record[6])
		if err != nil {
			return fmt.Errorf("Could not execute sql statement in UpdateStaticCitiesTable(): %w", err)
		}
		linenumber++
	}

	if err := database.UpdateCityWeatherStationId(c.db); err != nil {
		return fmt.Errorf("Error updating weatherstation id in UpdateStaticCitiesTable(): %w", err)
	}
	return nil
}

func (c *DataClient) UpdateStaticSurfSpotTable() {
	sqlStmnt, err := c.db.Prepare(`
			INSERT INTO surfspot (id, name, latitude, longitude, city_id, break_type, orientation, nearest_buoy)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8)
		`)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlStmnt.Close()

	_, err = c.db.Exec(`TRUNCATE TABLE surfspot RESTART IDENTITY CASCADE`)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(dbCitiesList)
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
			log.Fatal(err)
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

		nearestBouy := spacial.NearestBuoy(lat, lon, c.db)

		_, err = sqlStmnt.Exec(id, record[1], lat, lon, city_id, record[5], orientation, nearestBouy)
		if err != nil {
			log.Fatalf("line %d: insert failed: %v", linenumber, err)
		}
		linenumber++
	}
}
