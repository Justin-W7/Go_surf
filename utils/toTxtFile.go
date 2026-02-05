package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func BouyDataToTextFile(data *http.Response, stationID int) {
	timestamp := time.Now().Format("2006-01-02")

	fileName := fmt.Sprintf("%s_bouydata_%d.txt", timestamp, stationID)
	fullPath := filepath.Join("db/bouy_data", fileName)

	file, err := os.Create(fullPath)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	_, _ = io.Copy(file, data.Body)
}
