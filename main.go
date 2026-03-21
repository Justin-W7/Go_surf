package main

import (
	"go_surf/database"
	"go_surf/menu"
)

func main() {

	db := database.ConnectDatabase()
	defer database.DisconnectDatabase(db)

	menu.StartMenuLoop(db)
}
