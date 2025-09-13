package simplifyutils

import (
	"encoding/json"
	"fmt"
	"os"

	api "github.com/Matrix030/SimplifyJobsDaemon/internal/simplifyapi"
)

var fileName = "jobs.json"
var newJobsFileName = "newJobsOnly.json"

func JsonFileWriter(jobsData api.Jobs) error {
	jsonData, err := json.MarshalIndent(jobsData, "", " ")
	if err != nil {
		fmt.Println("There was an error while marshalling data", err)
		return err
	}

	err = os.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		fmt.Println("There was an error which writing the json file", err)
		return err
	}

	fmt.Println("Successfully wrote the jobs to the jobs.json")

	return nil
}

// WriteNewJobsOnly writes only the new jobs to a separate file
// This file gets rewritten (not appended) each time new jobs are found
func WriteNewJobsOnly(newJobs api.Jobs) error {
	if len(newJobs) == 0 {
		// If no new jobs, don't modify the file
		fmt.Println("No new jobs found - keeping existing newJobsOnly.json unchanged")
		return nil
	}

	var existingJobs api.Jobs

	// Try to read existing jobs from the file
	if _, err := os.Stat(newJobsFileName); err == nil {
		// File exists, read it
		existingData, err := os.ReadFile(newJobsFileName)
		if err != nil {
			return fmt.Errorf("error reading existing new jobs file: %w", err)
		}

		// Only try to unmarshal if file is not empty
		if len(existingData) > 0 {
			err = json.Unmarshal(existingData, &existingJobs)
			if err != nil {
				return fmt.Errorf("error unmarshalling existing jobs data: %w", err)
			}
		}
	}
	// If file doesn't exist or is empty, existingJobs remains empty slice

	// Prepend new jobs to existing jobs (new jobs at the top)
	combinedJobs := append(newJobs, existingJobs...)

	// Marshal the combined jobs to JSON with proper indentation
	jsonData, err := json.MarshalIndent(combinedJobs, "", " ")
	if err != nil {
		return fmt.Errorf("error marshalling combined jobs data: %w", err)
	}

	// Write to file (this overwrites with the combined data)
	err = os.WriteFile(newJobsFileName, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing combined jobs file: %w", err)
	}

	fmt.Printf("Successfully prepended %d new jobs to %s (total: %d jobs)\n",
		len(newJobs), newJobsFileName, len(combinedJobs))
	return nil
}
func LoadExistingJobs() (api.Jobs, error) {
	fileData := &api.Jobs{}
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		fmt.Println("File does not exists")
		return api.Jobs{}, err
	}

	//read file content
	data, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("There was an error while reading the file")
		return api.Jobs{}, err
	}

	err = json.Unmarshal(data, &fileData)

	return *fileData, nil
}
