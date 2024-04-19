package src

import (
	"encoding/json"
	"io"
	"os"
)

// ReadExercises reads the exercises.json file and returns a slice of Exercise structs
func ReadExercises(path string) ([]Exercise, error) {
	// Open the exercises.json file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the contents of the file
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into a slice of Exercise structs
	var exercises []Exercise
	if err := json.Unmarshal(data, &exercises); err != nil {
		return nil, err
	}

	return exercises, nil
}

func WriteExercises(path string, exercises []Exercise) error {
	// Marshal the exercises slice into JSON
	data, err := json.Marshal(exercises)
	if err != nil {
		return err
	}

	// Write the JSON data to the exercises.json file
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
