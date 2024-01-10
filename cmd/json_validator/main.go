package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stdout, "Usage: json_validator <json_string>")
		os.Exit(1)
	}
	jsonString := os.Args[1]

	err := isValidJSON(jsonString)
	if err != nil {
		fmt.Fprintln(os.Stdout, err.Error())
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, "valid JSON")
	os.Exit(0)
}

func isValidJSON(jsonString string) error {
	var js interface{}
	err := json.Unmarshal([]byte(jsonString), &js)
	if err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	return nil
}
