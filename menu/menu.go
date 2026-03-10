package menu

import (
	"database/sql"
	"fmt"
)

func StartMenuLoop(db *sql.DB) {
	var input string
	for {
		printMenu()
		fmt.Print("> ")
		fmt.Scan(&input)

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
	fmt.Println("b - Option B (placeholder)")
	fmt.Println("c - Option C (placeholder)")
	fmt.Println("d - Option D (placeholder)")
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
		fmt.Println("b - place holder.")
	case "c":
		fmt.Println("c - place holder.")
	case "d":
		fmt.Println("d - place holder.")
	case "h":
		fmt.Println("h - place holder.")
	default:
		fmt.Println("Invalid selection, try again.")

	}
}
