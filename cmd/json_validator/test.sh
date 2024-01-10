#!/bin/bash

# Specify the directory containing test files
TEST_DIR="test/"

# Iterate through files in the test directory
for file in "$TEST_DIR"/*; do
    if [ -e "$file" ] ;then
        echo "Testing file: $file"
        
        # Run the json_validator program on the file
        ./json_validator "$(cat "$file")"
        
        # Capture the exit code of the json_validator program
        exit_code=$?

        # Check the exit code and print the result
        if [ $exit_code -eq 0 ]; then
            echo "Test Passed: Valid JSON"
        else
            echo "Test Failed: Invalid JSON"
        fi

        echo "---------------------------"
    fi
done
