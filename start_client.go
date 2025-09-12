package main

import (
	"fmt"
	api "github.com/Matrix030/simplify_jobs_cli/internal/simplifyapi"
	utils "github.com/Matrix030/simplify_jobs_cli/internal/simplifyutils"
	"time"
)

type config struct {
	jobClient api.Client
}

func startClient(cfg *config) {
	var oldJobs api.Jobs
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	//running once on startup
	oldJobsData, err := utils.LoadExistingJobs()
	if err != nil {
		fmt.Printf("File not found: %v\n", err)
		// Fetch new jobs since file doesn't exist
		oldJobs, err = fetchAndProcessJobs(cfg)
		if err != nil {
			fmt.Printf("An error occurred while fetching the jobs %v\n", err)
			return // or handle error appropriately
		}
		// Only write if fetch succeeded
		utils.JsonFileWriter(oldJobs)
	} else {
		oldJobs = oldJobsData
	}

	for {
		<-ticker.C
		newJobs, err := fetchAndProcessJobs(cfg)
		if err != nil {
			fmt.Printf("An Error Occurred while getting the new jobs\n")
			return
		}

		oldJobs = newJobs

		newJobSlice := compareJobs(oldJobs, newJobs)
		//TODO:
		utils.SendNotification(newJobSlice)
	}
}

func fetchAndProcessJobs(cfg *config) (api.Jobs, error) {
	jobs, err := cfg.jobClient.GetJobData()
	if err != nil {
		fmt.Printf("There was an error %v\n", err)
		return api.Jobs{}, err
	}
	return jobs, nil
}
