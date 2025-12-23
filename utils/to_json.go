// Package utils contains program helper functions.
package utils

import (
	"encoding/json"
	"os"
)

func ToJSONFile[T any](data T, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}
