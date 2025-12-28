// Package utils contains program helper functions.
package utils

import (
	"encoding/json"
	"os"
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
