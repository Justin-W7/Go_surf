// Package utils contains program helper functions.
package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// ToJSONFile writes any data structure to a JSON file.
//
// Parameters:
//   - data: The data to be unmarshalled.
//   - filename: The name of the file to write the output to.
//
// Returns:
//   - error: An error if the file creation or JSON encoding fails.
func ToJSONFile[T any](data T, filename string) error {
	path := filepath.Join("test", filename)

	file, err := os.Create(path)
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
