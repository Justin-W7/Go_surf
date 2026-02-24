package main

import (
	"go_surf/menu"
	"go_surf/utils"
)

func main() {

	db := utils.ConnectDatabase()
	defer utils.DisconnectDatabase(db)

	menu.StartMenuLoop(db)
}
