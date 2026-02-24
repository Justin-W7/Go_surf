package utils

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

func DatabaseMenu(db *sql.DB) {
	i := ""
	fmt.Println()
	fmt.Println("a - Update static Buoy table from csv 'buoys.csv'.")
	fmt.Print("> ")
	fmt.Scan(&i)

	if i == "a" {
		updateBuoyTable(db)
	}
}

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

func DisconnectDatabase(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Println("Error closing database: ", err)
		return
	}
	fmt.Println("Database disconnect succesful.")
}

func updateBuoyTable(db *sql.DB) {
	file, err := os.Open("/home/waffles/personal/projects/Go_surf/api/station_lists/buoys.csv")
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
