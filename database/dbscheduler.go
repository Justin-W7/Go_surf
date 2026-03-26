package database

import (
	"database/sql"
	"fmt"
	"go_surf/api"
	"time"
)

func StartDataIngestion(db *sql.DB) {
	go tickerRunner(15*time.Minute, func() { updateBuoyData(db) })
	go tickerRunner(time.Hour, func() { updateWeatherData(db) })
	go tickerRunner(15*time.Minute, func() { updateSurfConditions(db) })
}

// tickerRunner manages job execution timing.
// Job runs once when tickerRunner is called,
// then runs at time implementation.
func tickerRunner(interval time.Duration, job func()) {
	job()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		job()
	}

}

func updateBuoyData(db *sql.DB) {
	fmt.Println("Updating real time buoy data.")
	api.FetchNDBCBuoyDataFromStationList(api.NDBCBouyDataURL, api.STATION_ID_FILE)
	if err := UpdateRTBuoyDataTable(db); err != nil {
		fmt.Println("Error updating buoy data: ", err)
	}
	MoveOldBuoyData()
}

func updateWeatherData(db *sql.DB) {
	fmt.Println("Updating real time weather data.")
	if err := UpdateRTWeatherTable(db); err != nil {
		fmt.Println("Error updating weather data: ", err)
	}
}

func updateSurfConditions(db *sql.DB) {
	fmt.Println("Updating current surf conditions.")
	if err := UpdateCurrentSurfConditions(db); err != nil {
		fmt.Println("Error updating current surf conditions: ", err)
	}
}
