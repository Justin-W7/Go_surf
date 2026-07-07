package dbLib

import (
	meteo "Go_surf_redesign/src/backend/api"
	"context"
	"fmt"
	"time"
)

func StartDataIngestion(ctx context.Context, db *DataClient, api *meteo.Client) error {
	fmt.Println("Starting data ingestion.")

	go func() {
		weatherReady := false

		nextBuoy := time.Now()
		nextWeather := time.Now()
		nextSurf := time.Now()

		for {
			now := time.Now()

			// 1. Buoy data
			if now.After(nextBuoy) {
				db.UpdateRTBuoyData(ctx, api)
				fmt.Println("DataBase: updating current buoy data.")
				nextBuoy = now.Add(15 * time.Minute)
			}

			// 2. Weather data
			if now.After(nextWeather) {
				db.UpdateRTWeatherData(ctx, api)
				nextWeather = now.Add(time.Hour)

				weatherReady = true
			}

			// 3. Surf conditions
			if weatherReady && now.After(nextSurf) {
				db.UpdateCurrentSurfConditions(api)
				nextSurf = now.Add(15 * time.Minute)
			}

			time.Sleep(30 * time.Second)
		}
	}()

	fmt.Println("Data ingestion started.")
	return nil
}

// tickerRunner manages job execution timing.
// Job runs once when tickerRunner is called,
// then runs again at time implementation.
// func tickerRunner(interval time.Duration, job func()) {
// 	job()
// 	ticker := time.NewTicker(interval)
// 	defer ticker.Stop()

// 	for range ticker.C {
// 		job()
// 	}
// }
