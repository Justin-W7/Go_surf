package database

import (
	"database/sql"
	"fmt"
	"go_surf/api"
	"time"
)

func StartDataIngestion(db *sql.DB) {

	go StartRTBuoyDataIngestion(db)
	go StartRTWeatherDataIngestion(db)

}

func StartRTBuoyDataIngestion(db *sql.DB) {
	for {
		api.FetchNDBCBuoyDataFromStationList(api.NDBCBouyDataURL, api.STATION_ID_FILE)
		fmt.Println("Updating Real Time Buoy Data")
		UpdateRTBuoyDataTable(db)
		time.Sleep(15 * time.Minute)
	}
}

func StartRTWeatherDataIngestion(db *sql.DB) {
	for {
		fmt.Println("Updating Real Time Weather Data")
		UpdateRTWeatherTable(db)
		time.Sleep(time.Hour)
	}
}
