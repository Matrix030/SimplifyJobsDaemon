package main

import (
	api "github.com/Matrix030/simplify_jobs_cli/internal/simplifyapi"
)

func compareJobs(oldJobs, newJobs api.Jobs) api.Jobs {
	oldJobMap := make(map[string]bool)
	for _, job := range oldJobs {
		oldJobMap[job.ID] = true
	}

	var newJobsOnly api.Jobs
	for _, job := range newJobs {
		if job.Active && !oldJobMap[job.ID] {
			newJobsOnly = append(newJobsOnly, job)
		}
	}
	return newJobsOnly
}
