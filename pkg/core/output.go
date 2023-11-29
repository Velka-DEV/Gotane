package core

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

// GetOutputPath returns the path to the output file for a given status and time
// If directories do not exist, they will be created
// Example: {currentDirectory}/output/free/2021-01-01_16-34-24.txt
func createOutputPath(status CheckStatus, time time.Time) string {
	currentDirectory, _ := os.Getwd()
	outputDirectory := fmt.Sprintf("%s/output", currentDirectory)
	statusDirectory := fmt.Sprintf("%s/%s", outputDirectory, status.String())
	filePath := fmt.Sprintf("%s/%s.txt", statusDirectory, time.Format("2006-01-02-15-04-05"))

	if _, err := os.Stat(outputDirectory); os.IsNotExist(err) {
		os.Mkdir(outputDirectory, 0755)
	}

	if _, err := os.Stat(statusDirectory); os.IsNotExist(err) {
		os.Mkdir(statusDirectory, 0755)
	}

	return filePath
}

func WriteResultToFile(result *CheckResult, info *CheckerInfo) error {
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
	for _, key := range keys {
		sb.WriteString(fmt.Sprintf("|%s=%s", key, result.Captures[key]))
	}

	sb.WriteString(fmt.Sprintf("|%s", info.StartTime.Format("2006-01-02 15:04:05")))

	return writeLineToFile(createOutputPath(result.Status, info.StartTime), []byte(sb.String()))
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
