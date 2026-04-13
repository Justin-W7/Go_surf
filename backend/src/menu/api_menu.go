package menu

import (
	"database/sql"
	"fmt"
)

func ApiMenu(db *sql.DB) {
	i := ""
	for i != "q" {
		fmt.Println()
		fmt.Println("------------------------------------------------------------")
		fmt.Println()
		fmt.Println("a - TEST API getCities()")
		fmt.Println()
		fmt.Println("q - BACK")
		fmt.Println("------------------------------------------------------------")
		fmt.Print("> ")
		fmt.Scan(&i)
		fmt.Println()

		switch i {
		case "a":
			fmt.Println("api_endpoint.go")
		default:
			if i != "q" {
				fmt.Println("Invalid selection, try again.")
			}
		}
	}
}
