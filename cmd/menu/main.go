package main

import (
	"go_surf/backend/src/database"
	"go_surf/backend/src/menu"
)

func main() {

	db := database.ConnectDatabase()
	defer database.DisconnectDatabase(db)

	menu.StartMenuLoop(db)
}
