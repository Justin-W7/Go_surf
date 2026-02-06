package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func BouyDataToTextFile(data *http.Response, id string) {
	t := time.Now().Format("2006-01-02 15:04:05")

	fileName := fmt.Sprintf("%s_bouydata_%s.txt", t, id)
	fullPath := filepath.Join("db/raw_data/NDBC_bouy_data", fileName)

	file, err := os.Create(fullPath)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	_, _ = io.Copy(file, data.Body)
}
