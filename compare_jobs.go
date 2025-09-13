package main

import (
	api "github.com/Matrix030/SimplifyJobsDaemon/internal/simplifyapi"
)

func compareJobs(oldJobs, newJobs api.Jobs) api.Jobs {
	// Create a set of old job IDs for fast lookup
	oldJobIDs := make(map[string]bool)
	for _, job := range oldJobs {
		oldJobIDs[job.ID] = true
	}

	var newJobsOnly api.Jobs

	// Start from the end (newest jobs) and work backwards
	// Stop when we find a job that exists in the old jobs
	for i := len(newJobs) - 1; i >= 0; i-- {
		job := newJobs[i]

		// If we find a job that already exists in old jobs,
		// we can stop because all jobs before this one should also exist
		if oldJobIDs[job.ID] {
			break
		}

		// Only consider active jobs
		if !job.Active {
			continue
		}

		// Filter by sponsorship - only include jobs that offer sponsorship or are unclear
		if !isEligibleSponsorship(job.Sponsorship) {
			continue
		}

		// This is a new job that passes our filters
		// Prepend to maintain chronological order (oldest new job first)
		newJobsOnly = append(api.Jobs{job}, newJobsOnly...)
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
