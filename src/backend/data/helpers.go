package data

import (
	"fmt"
	"os"
	"path/filepath"
)

func FilePathBuilder() string {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory: ", err)
	}

	path := filepath.Join(currentDir, "src", "backend", "data")
	return path
}
