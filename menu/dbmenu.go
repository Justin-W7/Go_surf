package menu

import (
	"database/sql"
	"fmt"
	"go_surf/api"
	"go_surf/database"
)

// DatabaseMenu is the cli for static database updates.
func DatabaseMenu(db *sql.DB) {
	i := ""
	fmt.Println()
	fmt.Println("a - Update static buoy table from csv 'buoys.csv'.")
	fmt.Println("b - Update static surfspot table from csv 'surfspots.csv'.")
	fmt.Println("c - Update static cities table from csv 'cities.csv'.")
	fmt.Println("d - Fetch current buoy data.")
	fmt.Println("e - Update real time buoy table.")
	fmt.Println("f - Update real time weather table.")
	fmt.Println("------------------------------------------------------------")
	fmt.Println()
	fmt.Print("> ")
	fmt.Scan(&i)

	switch i {
	case "a":
		database.UpdateBuoyTable(db)
	case "b":
		database.UpdateSurfSpotTable(db)
	case "c":
		database.UpdateCitiesTable(db)
	case "d":
		api.FetchNDBCBuoyDataFromStationList(api.NDBCBouyDataURL, api.STATION_ID_FILE)
	case "e":
		database.UpdateRTBuoyDataTable(db, api.NDBCBouyDataURL, api.STATION_ID_FILE)
	case "f":
		database.UpdateRTWeatherTable(db)
	default:
		fmt.Println("Invalid selection, try again.")
	}
}
