package main

import (
	"fmt"
	api "github.com/Matrix030/SimplifyJobsDaemon/internal/simplifyapi"
	utils "github.com/Matrix030/SimplifyJobsDaemon/internal/simplifyutils"
	"time"
)

type config struct {
	jobClient api.Client
}

func startClient(cfg *config) {
	var oldJobs api.Jobs
	ticker := time.NewTicker(30 * time.Minute) //Interval
	defer ticker.Stop()

	// Running once on startup
	oldJobsData, err := utils.LoadExistingJobs()
	if err != nil {
		fmt.Printf("File not found: %v\n", err)
		// Fetch new jobs since file doesn't exist
		oldJobs, err = fetchAndProcessJobs(cfg)
		if err != nil {
			fmt.Printf("An error occurred while fetching the jobs: %v\n", err)
			return
		}
		// Write the initial jobs to file
		err = utils.JsonFileWriter(oldJobs)
		if err != nil {
			fmt.Printf("Error writing initial jobs to file: %v\n", err)
			return
		}
	} else {
		oldJobs = oldJobsData
		fmt.Printf("Loaded %d existing jobs from file\n", len(oldJobs))
	}

	fmt.Println("Starting job monitoring... Press Ctrl+C to stop")

	// Main monitoring loop
	for {
		<-ticker.C
		fmt.Println("Checking for new jobs...")

		newJobs, err := fetchAndProcessJobs(cfg)
		if err != nil {
			fmt.Printf("An error occurred while getting new jobs: %v\n", err)
			continue // Don't exit, just skip this iteration
		}

		// Compare old jobs with new jobs to find newly added ones
		newJobSlice := compareJobs(oldJobs, newJobs)

		// Always write the new jobs file (even if empty)
		err = utils.WriteNewJobsOnly(newJobSlice)
		if err != nil {
			fmt.Printf("Error writing new jobs file: %v\n", err)
		}

		if len(newJobSlice) > 0 {
			fmt.Printf("Found %d new jobs!\n", len(newJobSlice))

			fmt.Println("Scraping job Descriptions....")
			descriptions := scrapeNewJobDescriptions(newJobSlice)

			err = utils.SaveJobDescriptions(descriptions)
			if err != nil {
				fmt.Printf("Error saving descriptions: %v\n", err)
			}

			// Send notification for new jobs
			err = utils.SendNotification(newJobSlice)
			if err != nil {
				fmt.Printf("Error sending notification: %v\n", err)
			}

			// Update the stored jobs file with all current jobs
			err = utils.JsonFileWriter(newJobs)
			if err != nil {
				fmt.Printf("Error updating jobs file: %v\n", err)
			}
		} else {
			fmt.Println("No new jobs found")
		}

		// Update oldJobs for next comparison
		oldJobs = newJobs
	}
}

func fetchAndProcessJobs(cfg *config) (api.Jobs, error) {
	jobs, err := cfg.jobClient.GetJobData()
	if err != nil {
		return api.Jobs{}, fmt.Errorf("failed to fetch job data: %w", err)
	}

	fmt.Printf("Fetched %d total jobs\n", len(jobs))
	return jobs, nil
}
