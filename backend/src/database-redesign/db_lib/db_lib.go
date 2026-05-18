package dbLib

import (
	"database/sql"
	"log"
)

const (
	constStr   = "user=postgres password=pass dbname=surfdate sslmode=disable"
	psqlDriver = "postgers"
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
func BuildStaticBuoyTable()     {}
func BuildStaticSurfSpotTable() {}
func BuildStaticCitiesTable()   {}

// Need to update real time data tables
func UpdateRTBuoyData()            {}
func UpdateRTWeatherData()         {}
func UpdateCurrentSurfConditions() {}

// Need to update forecasted data
func UpdateForecastedBuoyData()    {}
func UpdateForecastedWeatherData() {}
func UpdateForecastedConditions()  {}
*/
