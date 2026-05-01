package database

import (
	"database/sql"
	"fmt"
	"go_surf/backend/src/api"
	"go_surf/backend/src/config"
	"sync"
	"time"
)

// StartDataIngestion starts realtime data pipeline.
// surfConditionStart is not started until updateWehaterData has run once.
// After that, updateSurfConditions will run independently.
func StartDataIngestion(db *sql.DB) {
	var surfConditionStart sync.Once

	go tickerRunner(15*time.Minute, func() { updateBuoyData(db) })

	go tickerRunner(time.Hour, func() {
		updateWeatherData(db)

		surfConditionStart.Do(func() {
			go tickerRunner(15*time.Minute, func() { updateSurfConditions(db) })
		})

	})
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

func updateBuoyData(db *sql.DB) error {
	fmt.Println("Updating real time buoy data.")

	err := api.FetchNDBCBuoyDataFromStationList(config.NDBCBouyDataURL, config.STATION_ID_FILE)
	if err != nil {
		return fmt.Errorf("%v:%w", "updateBuoyData", err)
	}

	if err := UpdateRTBuoyDataTable(db); err != nil {
		fmt.Println("Error updating buoy data: ", err)
	}
	MoveOldBuoyData()

	return nil
}

func updateWeatherData(db *sql.DB) error {
	fmt.Println("Updating real time weather data.")
	if err := UpdateRTWeatherData(db); err != nil {
		fmt.Println("Error updating weather data: ", err)
	}
	return nil
}

func updateSurfConditions(db *sql.DB) error {
	fmt.Println("Updating current surf conditions.")

	if err := UpdateCurrentSurfConditions(db); err != nil {
		fmt.Println("Error updating current surf conditions: ", err)
	}
	return nil
}
