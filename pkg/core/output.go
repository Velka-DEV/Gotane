package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// GetOutputPath returns the path to the output file for a given status and time
// If directories do not exist, they will be created
// Example: {currentDirectory}/output/free/2021-01-01_16-34-24.txt
func createOutputPath(status CheckStatus, t time.Time, basePath string) (string, error) {
	outputDirectory := filepath.Join(basePath, "output")
	statusDirectory := filepath.Join(outputDirectory, status.String())
	filePath := filepath.Join(statusDirectory, fmt.Sprintf("%s.txt", t.Format("2006-01-02-15-04-05")))

	// Use os.MkdirAll to create all necessary directories
	err := os.MkdirAll(statusDirectory, 0755)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func WriteResultToFile(result *CheckResult, info *CheckerInfo, basePath string) error {
	sb := strings.Builder{}

	sb.WriteString(result.Combo.String())

	// Create a slice to hold the keys of the map
	keys := make([]string, 0, len(result.Captures))

	// Add all keys to the slice
	for key := range result.Captures {
		keys = append(keys, key)
	}

	// Sort the keys
	sort.Strings(keys)

	// Iterate over the sorted keys
	for i, key := range keys {

		if i == 0 {
			sb.WriteString(fmt.Sprintf("|%s=%s", key, result.Captures[key]))
			continue
		}
		sb.WriteString(fmt.Sprintf(",%s=%s", key, result.Captures[key]))
	}

	sb.WriteString(fmt.Sprintf("|%s", info.StartTime.Format("2006-01-02 15:04:05")))

	outputPath, err := createOutputPath(result.Status, info.StartTime, basePath)

	if err != nil {
		return err
	}

	return writeLineToFile(outputPath, []byte(sb.String()))
}

func writeLineToFile(path string, data []byte) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("Could not open file: ", err)
		return err
	}

	defer file.Close()

	line := append(data, []byte("\n")...)

	_, err2 := file.WriteString(string(line))

	if err2 != nil {
		fmt.Println("Could not " + string(data) + " text to " + path + "(" + err2.Error() + ")")

	}

	return nil
}
