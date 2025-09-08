package main

import (
	"fmt"
	api "github.com/Matrix030/simplify_jobs_cli/internal/simplifyapi"
	"time"
)

type config struct {
	jobClient api.Client
}

func startClient(cfg *config) {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	//running once on startup

	oldJobs, err := fetchAndProcessJobs(cfg)
	if err != nil {
		fmt.Printf("An error occurred while fetch the job %v\n", err)
		return
	}

	for {
		<-ticker.C
		newJobs, err := fetchAndProcessJobs(cfg)
		if err != nil {
			fmt.Printf("An Error Occurred while getting the new jobs\n")
			return
		}

		oldJobs = newJobs
		//TODO: comparison logic goes here
		compareJobs(oldJobs, newJobs)
		//TODO: sendNotification(newJobSlice)
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
