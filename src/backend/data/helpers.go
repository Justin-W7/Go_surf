package data

import (
	"fmt"
	"os"
	"path/filepath"
)

func FilePathBuilder(subFolders ...string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory: ", err)
	}

	basePath := filepath.Join(currentDir, "src", "backend", "data")

	if len(subFolders) == 0 {
		return basePath
	}

	// Below is a work around for the following:
	// return filepath.Join(basePath, subFolders...)
	//
	// Compiler complained about the slice, even though
	// it looks like filepath.Join() accepts (elem ...string) as an arg
	// per documentation.
	//
	// Compiler returned this error:
	// src/backend/data/helpers.go:21:33: too many arguments in call to filepath.Join
	// have (string, []string...)
	// want (...string)

	parts := []string{filepath.Join(currentDir, "src", "backend", "data")}
	parts = append(parts, subFolders...)

	return filepath.Join(parts...)
}
