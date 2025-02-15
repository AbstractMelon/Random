package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

type Config struct {
	Replacement string `json:"replacement"`
}

func loadConfig(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func processFile(inputFile, outputFile string, config Config) (int, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	output, err := os.Create(outputFile)
	if err != nil {
		return 0, err
	}
	defer output.Close()

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(output)
	nwordRegex := regexp.MustCompile(`(?i)\bnigg[aeu][rh]?\b`)
	count := 0

	for scanner.Scan() {
		line := scanner.Text()
		matches := nwordRegex.FindAllString(line, -1)
		count += len(matches)
		
		if config.Replacement != "" {
			line = nwordRegex.ReplaceAllString(line, config.Replacement)
		} else {
			line = nwordRegex.ReplaceAllString(line, "")
		}

		writer.WriteString(line + "\n")
	}

	writer.Flush()

	if err := scanner.Err(); err != nil {
		return count, err
	}

	return count, nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <input_file> <output_file> [config_file]")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]
	configFile := "config.json"
	if len(os.Args) > 3 {
		configFile = os.Args[3]
	}

	config, err := loadConfig(configFile)
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	count, err := processFile(inputFile, outputFile, config)
	if err != nil {
		fmt.Println("Error processing file:", err)
		return
	}

	fmt.Printf("Processing complete. %d occurrences removed/replaced.\n", count)
}
