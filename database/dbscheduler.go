package database

import (
	"database/sql"
	"fmt"
	"go_surf/api"
	"time"
)

func StartRTBuoyDataIngestion(db *sql.DB) {
	fmt.Println("Updating Real Time Buoy Data")
	UpdateRTBuoyDataTable(db, api.NDBCBouyDataURL, api.STATION_ID_FILE)

	for {
		fmt.Println("Updating Real Time Buoy Data")
		UpdateRTBuoyDataTable(db, api.NDBCBouyDataURL, api.STATION_ID_FILE)
		time.Sleep(15 * time.Minute)
	}
}

func StartRTWeatherDataIngestion(db *sql.DB) {
	fmt.Println("Updating Real Time Weather Data")
	UpdateRTWeatherTable(db)

	for {
		fmt.Println("Updating Real Time Weather Data")
		UpdateRTWeatherTable(db)
		time.Sleep(time.Hour)
	}
}

func StartDataIngestion(db *sql.DB) {

	go StartRTBuoyDataIngestion(db)
	go StartRTWeatherDataIngestion(db)

}
