package database

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"go_surf/api"
	"go_surf/models"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// DatabaseMenu is the cli for static database updates.
func DatabaseMenu(db *sql.DB) {
	i := ""
	fmt.Println()
	fmt.Println("a - Update static buoy table from csv 'buoys.csv'.")
	fmt.Println("b - Update static surfspot table from csv 'surfspots.csv'.")
	fmt.Println("c - Update static cities table from csv 'cities.csv'.")
	fmt.Println("d - Fetch current buoy data.")
	fmt.Println("e - Update real time buoy table.")
	fmt.Println("------------------------------------------------------------")
	fmt.Println()
	fmt.Print("> ")
	fmt.Scan(&i)

	switch i {
	case "a":
		updateBuoyTable(db)
	case "b":
		updateSurfSpotTable(db)
	case "c":
		updateCitiesTable(db)
	case "d":
		api.FetchNDBCBuoyDataFromStationList(api.NDBCBouyDataURL, api.STATION_ID_FILE)
	case "e":
		updateRealTimeBuoyDataTable(db)
	default:
		fmt.Println("Invalid selection, try again.")
	}
}

// ConnectDatabase connects the program to the local database on program start.
func ConnectDatabase() *sql.DB {
	connStr := "user=postgres password=password dbname=surftest sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to database successfuly.")

	// fmt.Println("Table update succesful!")
	fmt.Println()
	time.Sleep(1 * time.Second)

	return db
}

// DisconnectDatabase disconnects the program from the local database on program shutdown.
func DisconnectDatabase(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Println("Error closing database: ", err)
		return
	}
	fmt.Println("Database disconnect succesful.")
}

// udpateBuoyTable updates the static table Buoys via a csv in the api folder.
func updateBuoyTable(db *sql.DB) {
	file, err := os.Open(api.DATABASE_BUOYS_FILE)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lineNumber := 0

	sqlStmnt, err := db.Prepare(`
		INSERT INTO buoys (id, name, latitude, longitude)
		VALUES ($1, $2, $3, $4)
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`TRUNCATE TABLE buoys RESTART IDENTITY CASCADE`)
	if err != nil {
		log.Fatal(err)
	}

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		if lineNumber == 0 {
			lineNumber++
			continue
		}

		// Convert csv fields to proper types
		id, err := strconv.Atoi(record[0])
		lat, err := strconv.ParseFloat(record[2], 64)
		long, err := strconv.ParseFloat(record[3], 64)

		_, err = sqlStmnt.Exec(id, record[1], lat, long)
		lineNumber++
	}

	defer sqlStmnt.Close()
}

// updateSurfSpotTable updates the static table surfspot via surfspots.csv
func updateSurfSpotTable(db *sql.DB) {
	file, err := os.Open(api.DATABASE_SURFSPOTS_FILE)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lineNumber := 0

	sqlStmnt, err := db.Prepare(`
		INSERT INTO surfspot (id, name, latitude, longitude, city_id, break_type, orientation)
		VALUES($1, $2, $3, $4, $5, $6, $7)		
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`TRUNCATE TABLE surfspot RESTART IDENTITY CASCADE`)
	if err != nil {
		log.Fatal(err)
	}

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		if lineNumber == 0 {
			lineNumber++
			continue
		}

		// convert csv fields to proper types
		id, err := strconv.Atoi(record[0])
		lat, err := strconv.ParseFloat(record[2], 64)
		long, err := strconv.ParseFloat(record[3], 64)
		city_id, err := strconv.Atoi(record[4])
		orientation, err := strconv.ParseFloat(record[6], 64)

		_, err = sqlStmnt.Exec(id, record[1], lat, long, city_id, record[5], orientation)
		if err != nil {
			log.Fatalf("line %d: insert failed: %v", lineNumber, err)
		}
		lineNumber++
	}
}

// updateCitiesTable updates static table "cities" with city record information.
func updateCitiesTable(db *sql.DB) {
	file, err := os.Open(api.DATABASE_CITIES_FILE)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lineNumber := 0

	sqlStmnt, err := db.Prepare(`
		INSERT INTO cities (id, name, latitude, longitude, country, state, county)
		VALUES($1, $2, $3, $4, $5, $6, $7)		
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`TRUNCATE TABLE cities RESTART IDENTITY CASCADE`)
	if err != nil {
		log.Fatal(err)
	}

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		if lineNumber == 0 {
			lineNumber++
			continue
		}

		// convert csv fields to proper types
		id, err := strconv.Atoi(record[0])
		lat, err := strconv.ParseFloat(record[2], 64)
		long, err := strconv.ParseFloat(record[3], 64)

		_, err = sqlStmnt.Exec(id, record[1], lat, long, record[4], record[5], record[6])
		if err != nil {
			log.Fatalf("line %d: insert failed: %v", lineNumber, err)
		}
		lineNumber++
	}
}

func updateRealTimeBuoyDataTable(db *sql.DB) {
	// get file
	folder := api.DATABASE_BUOYS_RT_RAW_DATA

	files, err := os.ReadDir(folder)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		name := file.Name()
		start := len(name) - 9
		buoystr := name[start : start+5]

		// convert buoystr to int
		buoyID, err := strconv.Atoi(buoystr)
		if err != nil {
			log.Fatal(err)
		}

		filepath := filepath.Join(folder, name)

		processRTBuoyFile(db, filepath, buoyID)
	}
}

func processRTBuoyFile(db *sql.DB, filepath string, buoyID int) {
	// open file
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// scan line with data
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "2026") {
			continue
		}

		buoyDataPoint, err := parseRTBuoyLine(line, buoyID)
		if err != nil {
			log.Fatal(err)
		}

		err = insertBuoyData(db, buoyDataPoint)
		if err != nil {
			log.Fatal(err)
		}
		break
	}
}

// parseRTBuoyLine returns a pointer to a models.BuoyDataPoint and an error.
func parseRTBuoyLine(line string, buoyID int) (*models.BuoyDataPoint, error) {
	data := strings.Fields(line)

	// parse time
	timeLayout := "2006 01 02 15 04"
	t, err := time.Parse(timeLayout, strings.Join(data[:5], " "))
	if err != nil {
		log.Fatal(err)
	}

	// parse data types
	windDirection, _ := parseDataFloat(data[5])
	windSpeed, _ := parseDataFloat(data[6])
	windGust, _ := parseDataFloat(data[7])
	waveHeightM, _ := parseDataFloat(data[8])
	dominantWavePeriod, _ := parseDataFloat(data[9])
	avgWavePeriod, _ := parseDataFloat(data[10])
	meanWaveDirection, _ := parseDataFloat(data[11])
	airTemperature, _ := parseDataFloat(data[12])
	waterTemperature, _ := parseDataFloat(data[13])

	return &models.BuoyDataPoint{
		BuoyID:                buoyID,
		RecordedAt:            t,
		WindDirectionDegT:     windDirection,
		WindSpeedMetersPerSec: windSpeed,
		WindGustMetersPerSec:  windGust,
		WaveHeightM:           waveHeightM,
		DominantWavePeriodSec: dominantWavePeriod,
		AvgWavePeriodSec:      avgWavePeriod,
		MeanWaveDirectionDegT: meanWaveDirection,
		AirTempDegC:           airTemperature,
		WaterTempDegC:         waterTemperature,
	}, nil
}

// parseDataTypes returns a pointer so it can return multiple states.
// This enables it to return a number, a missing value ("MM") or nil to mean no value, instead of 0
// which may be a valid value for the data.
func parseDataFloat(value string) (*float64, error) {
	if value == "MM" {
		return nil, nil
	}

	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Fatal(err)
	}
	return &result, nil
}

func insertBuoyData(db *sql.DB, p *models.BuoyDataPoint) error {
	_, err := db.Exec(`
			INSERT INTO real_time_buoy_data_points (
				buoy_id,
				recorded_at,
				wind_direction_degt,
				wind_speed_in_meters_per_sec,
				wind_gust_in_meters_per_sec,
				wave_height_m,
				dominant_wave_period_sec,
				avg_wave_period_sec,
				mean_wave_direction_degt,
				air_temp_degc,
				water_temp_degc
				)
			Values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`,
		p.BuoyID,
		p.RecordedAt,
		p.WindDirectionDegT,
		p.WindSpeedMetersPerSec,
		p.WindGustMetersPerSec,
		p.WaveHeightM,
		p.DominantWavePeriodSec,
		p.AvgWavePeriodSec,
		p.MeanWaveDirectionDegT,
		p.AirTempDegC,
		p.WaterTempDegC,
	)
	if err != nil {
		return err
	}

	return nil
}
