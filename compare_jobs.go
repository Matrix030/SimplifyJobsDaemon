package main

import (
	"fmt"
	"time"

	"github.com/Matrix030/SimplifyJobsDaemon/internal/scraper"
	api "github.com/Matrix030/SimplifyJobsDaemon/internal/simplifyapi"
)

func compareJobs(oldJobs, newJobs api.Jobs) api.Jobs {
	oldJobIDs := make(map[string]bool)
	for _, job := range oldJobs {
		oldJobIDs[job.ID] = true
	}

	var newJobsOnly api.Jobs

	for i := len(newJobs) - 1; i >= 0; i-- {
		job := newJobs[i]

		if oldJobIDs[job.ID] {
			break
		}

		if !job.Active {
			continue
		}

		if !isEligibleSponsorship(job.Sponsorship) {
			continue
		}

		newJobsOnly = append(api.Jobs{job}, newJobsOnly...)
	}

	return newJobsOnly
}

// New Scrape descriptions for new jobs
func scrapeNewJobDescriptions(jobs api.Jobs) []scraper.JobDescription {

}
