package main

import (
	"fmt"
	api "github.com/Matrix030/simplify_jobs_cli/internal/simplifyapi"
	"time"
)

// util function to be taken to the internal/utils
func GetUniqueSponsorshipValues(jobs api.Jobs) []string {
	unique := make(map[string]bool)
	for _, job := range jobs {
		unique[job.Sponsorship] = true
	}

	// convert keys to slice
	values := make([]string, 0, len(unique))
	for k := range unique {
		values = append(values, k)
	}
	return values
}

func main() {
	simplifyClient := api.NewClient(5 * time.Minute)
	cfg := &config{
		jobClient: simplifyClient,
	}

	jobs, err := cfg.jobClient.GetJobData()
	if err != nil {
		return
	}

	//utils part of the program
	unique := GetUniqueSponsorshipValues(jobs)
	for _, value := range unique {
		fmt.Println("Unique Sponsorship values:", value)

	}
	// startClient(cfg)
}
