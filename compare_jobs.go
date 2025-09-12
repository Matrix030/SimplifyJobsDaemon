package main

import (
	api "github.com/Matrix030/simplify_jobs_cli/internal/simplifyapi"
)

func compareJobs(oldJobs, newJobs api.Jobs) api.Jobs {
	// Create a map to store the latest date (posted or updated) for each job ID from old jobs
	oldJobDates := make(map[string]int)
	for _, job := range oldJobs {
		// Use the more recent date between posted and updated
		latestDate := job.DatePosted
		if job.DateUpdated > latestDate {
			latestDate = job.DateUpdated
		}
		oldJobDates[job.ID] = latestDate
	}

	var newJobsOnly api.Jobs
	for _, job := range newJobs {
		// Only consider active jobs
		if !job.Active {
			continue
		}

		// Filter by sponsorship - only include jobs that offer sponsorship or are unclear
		if !isEligibleSponsorship(job.Sponsorship) {
			continue
		}

		// Get the latest date for this new job
		newJobLatestDate := job.DatePosted
		if job.DateUpdated > newJobLatestDate {
			newJobLatestDate = job.DateUpdated
		}

		// Check if this job is new or has been updated
		oldDate, exists := oldJobDates[job.ID]

		if !exists {
			// This is a completely new job (ID not in old jobs)
			newJobsOnly = append(newJobsOnly, job)
		} else if newJobLatestDate > oldDate {
			// This job exists but has been updated since last check
			newJobsOnly = append(newJobsOnly, job)
		}
		// If newJobLatestDate <= oldDate, the job hasn't been updated, so skip it
	}

	return newJobsOnly
}

// isEligibleSponsorship checks if a job's sponsorship status is acceptable
// Returns true for jobs that are eligible (Other or Offers Sponsorship)
func isEligibleSponsorship(sponsorship string) bool {
	switch sponsorship {
	case "Other", "Offers Sponsorship":
		return true
	case "U.S. Citizenship is Required", "Does Not Offer Sponsorship":
		return false
	default:
		// For any unexpected sponsorship values, default to true to be safe
		// You can change this to false if you want to be more restrictive
		return true
	}
}
