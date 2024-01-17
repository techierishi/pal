package util

import (
	"os"
	"strings"
)

func AppendToFile(filename, content string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return err
	}

	return nil
}

func FileContains(filename, search string) (bool, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return false, err
	}

	return strings.Contains(string(data), search), nil
}
