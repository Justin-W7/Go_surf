package menu

import (
	"database/sql"
	"fmt"
	"go_surf/backend/src/api"
	"go_surf/backend/src/config"
	"go_surf/backend/src/database"
	"time"
)

// DatabaseMenu is the cli for static database updates.
func DatabaseMenu(db *sql.DB) {
	i := ""
	for i != "q" {
		fmt.Println()
		fmt.Println("------------------------------------------------------------")
		fmt.Println()
		fmt.Println("a - Update STATIC buoy table from csv 'buoys.csv'.")
		fmt.Println("b - Update STATIC surfspot table from csv 'surfspots.csv'.")
		fmt.Println("c - Update STATIC cities table from csv 'cities.csv'.")
		fmt.Println("d - Fetch current buoy data.")
		fmt.Println("e - Update real time buoy table.")
		fmt.Println("f - Update real time weather table.")
		fmt.Println("g - CLEAR real time table data.")
		fmt.Println("i - MOVE current rt buoy data to cold folder.")
		fmt.Println("j - TEST UpdateCurrentSurfConditions().")
		fmt.Println("x - REBUILD ALL STATIC TABLES")
		fmt.Println()
		fmt.Println("q - BACK")
		fmt.Println("------------------------------------------------------------")
		fmt.Print("> ")
		fmt.Scan(&i)
		fmt.Println()

		switch i {
		case "a":
			database.UpdateBuoyTable(db)
		case "b":
			database.UpdateSurfSpotTable(db)
		case "c":
			database.UpdateCitiesTable(db)
		case "d":
			api.FetchNDBCBuoyDataFromStationList(config.Path(config.NDBCBouyDataURL), config.Path(config.STATION_ID_FILE))
		case "e":
			database.UpdateRTBuoyDataTable(db)
		case "f":
			database.UpdateRTWeatherTable(db)
		case "g":
			database.ClearRTData(db)
		case "i":
			database.MoveOldBuoyData()
		case "j":
			database.UpdateCurrentSurfConditions(db)
		case "x":
			fmt.Println("Running X case...")

			fmt.Println("Updating buoy...")
			database.UpdateBuoyTable(db)

			fmt.Println("Updating cities...")
			database.UpdateCitiesTable(db)

			fmt.Println("Updating surf spots...")
			database.UpdateSurfSpotTable(db)

			fmt.Println("All static tables rebuilt.")
			time.Sleep(2 * time.Second)
		default:
			if i != "q" {
				fmt.Println("Invalid selection, try again.")
			}
		}
	}
}
