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
