package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// SaveRawBuoyDataToFile saves the raw HTTP response body of buoy data to a timestamped text file.
// The file is named using the current time and the provided buoy ID, and is stored in
// "database/raw_data/NDBC_buoy_data". Any errors encountered during file creation or writing
// are logged or printed.
func SaveRawBuoyDataToFile(data *http.Response, id string) {
	t := time.Now().Format("2006-01-02 15:04:05")

	fileName := fmt.Sprintf("%s_buoydata_%s.txt", t, id)
	fullPath := filepath.Join("database/raw_data/NDBC_buoy_data", fileName)

	file, err := os.Create(fullPath)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	_, err = io.Copy(file, data.Body)
	if err != nil {
		log.Fatal(err)
	}
}
