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
	var oldJobs api.Jobs
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	//running once on startup
	//TODO: put the data in a file after every fetch and check the if the file is present before fetching so that I could just assign oldJobs value to the data present in the file
	condition, oldJobsData, err := jsonFileChecker()
	if condition {
		// file with data found
		oldJobs = oldJobsData
	} else {
		oldJobs, err = fetchAndProcessJobs(cfg)
	}

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

		newJobSlice := compareJobs(oldJobs, newJobs)
		//TODO:
		sendNotification(newJobSlice)
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
