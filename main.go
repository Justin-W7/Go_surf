package main

import (
	"go_surf/backend/src/api"
	"go_surf/backend/src/database"
	"go_surf/backend/src/menu"
)

func main() {

	db := database.ConnectDatabase()
	defer database.DisconnectDatabase(db)

	go menu.StartMenuLoop(db)

	api.StartRouter(db)

}
