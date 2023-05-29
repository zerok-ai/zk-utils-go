package test

import (
	"fmt"
	"io"
	"os"
)

const (
// UnsortedWorkloadJs = “
// SortedWorkloadJs   = “
//
// EmptyScenarioJsonString = “
// NonScenarioJsonString   = `{"abc", 123}`
// ValidScenarioJsonString = “
)

func GetBytes(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	// Read the file content
	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	return content

}
