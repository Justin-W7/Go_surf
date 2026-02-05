package main

import (
	"go_surf/api"
)

func main() {

	inputfile := "api/station_lists/ndbcstations.txt"
	api.FetchNDBCBouyData(api.NDBCBouyDataURL, inputfile)
}
