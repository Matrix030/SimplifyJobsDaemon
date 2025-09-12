package simplifyutils

import (
	"encoding/json"
	"fmt"
	"os"

	api "github.com/Matrix030/simplify_jobs_cli/internal/simplifyapi"
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
		// If no new jobs, write an empty array to the file
		emptyJobs := api.Jobs{}
		jsonData, err := json.MarshalIndent(emptyJobs, "", " ")
		if err != nil {
			return fmt.Errorf("error marshalling empty jobs data: %w", err)
		}

		err = os.WriteFile(newJobsFileName, jsonData, 0644)
		if err != nil {
			return fmt.Errorf("error writing empty new jobs file: %w", err)
		}

		fmt.Println("No new jobs found - cleared newJobsOnly.json")
		return nil
	}

	// Marshal the new jobs to JSON with proper indentation
	jsonData, err := json.MarshalIndent(newJobs, "", " ")
	if err != nil {
		return fmt.Errorf("error marshalling new jobs data: %w", err)
	}

	// Write to file (this overwrites the existing file)
	err = os.WriteFile(newJobsFileName, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing new jobs file: %w", err)
	}

	fmt.Printf("Successfully wrote %d new jobs to %s\n", len(newJobs), newJobsFileName)
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

