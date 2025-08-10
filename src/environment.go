package src

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// LoadDotEnv loads environment variables from the specified .env file into the OS.
func LoadDotEnv(envFilePath string) error {
	file, err := os.Open(envFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lineCount int
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == '#' {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return errors.New("Line " + strconv.Itoa(lineCount) + " is invalid.")
		}

		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
		os.Setenv(key, value)

		lineCount++
	}

	return scanner.Err()
}

// UpdateEnvVar replaces or adds a key=value pair in the specified .env file.
func UpdateEnvVar(envFilePath string, key string, newValue string) error {
	file, err := os.Open(envFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	var found bool

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Update existing line if environment variable is found in file.
		if !found && strings.HasPrefix(line, key+"=") {
			lines = append(lines, fmt.Sprintf("%s='%s'", key, newValue))
			found = true
			continue
		}

		// Otherwise, add the line.
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Add to new line if environment variable is not found in file.
	if !found {
		lines = append(lines, fmt.Sprintf("%s=%s", key, newValue))
	}

	return os.WriteFile(envFilePath, []byte(strings.Join(lines, "\n")), 0644)
}
