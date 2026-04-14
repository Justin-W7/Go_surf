package main

import (
	"go_surf/backend/src/api"
	"go_surf/backend/src/database"
)

func main() {

	db := database.ConnectDatabase()
	defer database.DisconnectDatabase(db)

	api.StartRouter(db)
}
