package menu

import (
	"database/sql"
	"fmt"
	"go_surf/backend/src/database"
)

func StartMenuLoop(db *sql.DB) {
	var input string
	for {
		printMenu()
		fmt.Print("> ")
		fmt.Scan(&input)
		fmt.Println()

		if input == "q" {
			break
		}

		selectMenuItem(input, db)
		fmt.Println()
	}
}

func printMenu() {
	fmt.Println("--- Enter one of the options below ---")
	fmt.Println()
	fmt.Println("a - DATABASE MENU")
	fmt.Println("b - API MENU")
	fmt.Println("c - START data ingestion")
	fmt.Println("h - SHOW HELP")
	fmt.Println("q - QUIT PROGRAM")
	fmt.Println()
	fmt.Println("--------------------------------------")
}

func selectMenuItem(i string, db *sql.DB) {
	switch i {
	case "a":
		DatabaseMenu(db)
	case "b":
		ApiMenu(db)
	case "c":
		database.StartDataIngestion(db)
	case "h":
		fmt.Println("h - place holder.")
	default:
		fmt.Println("Invalid selection, try again.")
	}
}
