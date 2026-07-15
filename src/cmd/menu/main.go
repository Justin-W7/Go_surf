package main

import (
	meteo "Go_surf_redesign/src/backend/api"
	dbLib "Go_surf_redesign/src/backend/db_lib"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {

	// initialize database client
	dc := dbLib.NewDBClient()
	if err := dc.PingDB(); err != nil {
		log.Fatalf("From Main() - could not connect to database: %v", err)
	}
	// instantiate api client
	api := meteo.NewClient()

	// instantiate context
	ctx := context.Background()

	// meteo.StartRouter(dc.DB)
	// dbLib.StartDataIngestion(ctx, dc, api)

	mainMenu(ctx, dc, api)
}

func mainMenu(ctx context.Context, dc *dbLib.DataClient, api *meteo.Client) {
	input := ""
	for {
		fmt.Println("MAIN MENU")
		fmt.Println()
		fmt.Println("	(a) Start application - (starts data ingestion and router)")
		fmt.Println("	(b) Start API server")
		fmt.Println("	(c) Enter options menu")
		fmt.Println()
		fmt.Println("[q] Quit")
		fmt.Print("> ")
		fmt.Scan(&input)

		input = strings.TrimSpace(input)
		switch input {
		case "a":
			dbLib.StartDataIngestion(ctx, dc, api)
			meteo.StartRouter(dc.DB)
		case "b":
			meteo.StartRouter(dc.DB)
		case "c":
			optionsMenu(ctx, dc, api)
		case "q":
			quit(dc)
		}
	}
}

func optionsMenu(ctx context.Context, dc *dbLib.DataClient, api *meteo.Client) {
	fmt.Println()
	input := ""

	for input != "q" {
		fmt.Println("OPTIONS MENU")
		fmt.Println()
		fmt.Println("	(a) Load static data sets into database.")
		fmt.Println("	(b) Update real-time buoy data.")
		fmt.Println("	(c) Update real-time weather data.")
		fmt.Println("	(d) Update current surf condition data.")
		fmt.Println(" 	(e) Update static tide data.")
		fmt.Println()
		fmt.Println("[q] Back")
		fmt.Println()
		fmt.Print("> ")
		fmt.Scan(&input)

		input = strings.TrimSpace(input)
		switch input {
		case "a":
			dc.LoadStaticData()
		case "b":
			dc.UpdateRTBuoyData(ctx, api)
		case "c":
			dc.UpdateRTWeatherData(ctx, api)
		case "d":
			dc.UpdateCurrentSurfConditions(api)
		case "e":
			dc.UpdateStaticTideData()
		}
	}
}

func quit(db *dbLib.DataClient) {
	db.Close()
	fmt.Println("Goodbye")
	os.Exit(0)
}
