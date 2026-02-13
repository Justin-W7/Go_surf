package main

import (
	"go_surf/api"
)

var STATION_ID_FILE string = "api/station_lists/ndbcstations_CA.txt"

func main() {
	api.FetchNDBCBuoyDataFromStationList(api.NDBCBouyDataURL, STATION_ID_FILE)
}
