package core

import (
	"bufio"
	"os"
)

func ReadLinesFromFile(path string) ([]string, error) {
	var lines []string

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scan := bufio.NewScanner(file)
	for scan.Scan() {
		lines = append(lines, scan.Text())
	}

	return lines, nil
}

func LoadCombosFromFile(path string) ([]*Combo, error) {
	var lines, err = ReadLinesFromFile(path)

	if err != nil {
		return nil, err
	}

	var combos []*Combo

	for _, line := range lines {
		combo, err := ParseCombo(line)
		if err != nil {
			continue
		}

		combos = append(combos, combo)
	}

	return combos, nil
}
