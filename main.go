package main

import (
	"go_surf/api"
)

func main() {
	api.FetchNDBCBouyData(api.NDBCBouyDataURL, 41002)
}
