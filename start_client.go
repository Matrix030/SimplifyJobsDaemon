package main

import (
	"fmt"
	"time"

	"github.com/Matrix030/SimplifyJobsDaemon/internal/ai"
	"github.com/Matrix030/SimplifyJobsDaemon/internal/scraper"
	api "github.com/Matrix030/SimplifyJobsDaemon/internal/simplifyapi"
	utils "github.com/Matrix030/SimplifyJobsDaemon/internal/simplifyutils"
)

type config struct {
	jobClient api.Client
	llmClient *ai.LLMClient
	projects  []ai.Project
}

func startClient(cfg *config) {
	var oldJobs api.Jobs
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	//Check  LLM server health on startup
	if cfg.llmClient != nil {
		fmt.Println("Checking LLM server Health...")
		if err := cfg.llmClient.HealthCheck(); err != nil {
			fmt.Printf("Warning: LLM server not available: %v\n", err)
			fmt.Println("Job descriptions will be scraped but not analyzed.")
		} else {
			fmt.Println("LLM server is healthy")
		}
	}

	//Load existing jobs
	oldJobsData, err := utils.LoadExistingJobs()
	if err != nil {
		fmt.Printf("File not found: %v\n", err)
		oldJobs, err = fetchAndProcessJobs(cfg)
		if err != nil {
			fmt.Printf("An error occurred while fetching the jobs: %v\n", err)
			return
		}
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

	//main monitoring loop
	for {
		<-ticker.C
		fmt.Println("\n=== Checking for new jobs ===")

		newJobs, err := fetchAndProcessJobs(cfg)
		if err != nil {
			fmt.Printf("Error fetching jobs: %v\n", err)
			continue
		}

		//Compare and find new jobs
		newJobSlice := compareJobs(oldJobs, newJobs)

		//Write new jobs file
		err = utils.WriteNewJobsOnly(newJobSlice)
		if err != nil {
			fmt.Printf("Error writing new jobs file: %v\n", err)
		}

		if len(newJobSlice) > 0 {
			fmt.Printf("Found %d new jobs!\n", len(newJobSlice))

			//Scrape  descriptions
			fmt.Println("\n--- Scraping job descriptions ---")
			descriptions := scrapeNewJobDescriptions(newJobSlice)

			//Save descriptions
			err = utils.SaveJobDescriptions(descriptions)
			if err != nil {
				fmt.Printf("Error saving descriptions: %v\n", err)
			}

			//Analyze with LLM if available
			if cfg.llmClient != nil && len(cfg.projects) > 0 {
				fmt.Println("\n--- Analyzing jobs with LLM ---")
				analyzeJobsWithLLM(cfg, descriptions)
			}

			//Send notification
			err = utils.SendNotification(newJobSlice)
			if err != nil {
				fmt.Printf("Error sending notification: %v\n", err)
			}

			//Update jobs file
			err = utils.JsonFileWriter(newJobs)
			if err != nil {
				fmt.Printf("Error updating jobs file: %v\n", err)
			}

		} else {
			fmt.Println("No new jobs found")
		}

		oldJobs = newJobs
	}
}

func analyzeJobsWithLLM(cfg *config, descriptions []scraper.JobDescription) {
	for _, desc := range descriptions {
		if !desc.ScrapeSuccess {
			fmt.Printf("Skipping %s - %s (scrape failed)\n", desc.CompanyName, desc.Title)
			continue
		}

		fmt.Printf("\nAnalyzing: %s - %s\n", desc.CompanyName, desc.Title)

		analysis, err := cfg.llmClient.AnalyzeJob(desc.Description, cfg.projects)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
			continue
		}

		fmt.Printf(" Selected Projects: %v\n", analysis.SelectedProjects)
		fmt.Printf(" Reasoning: %s\n", analysis.Reasoning)

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
