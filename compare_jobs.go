package main

import (
	api "github.com/Matrix030/simplify_jobs_cli/internal/simplifyapi"
)

func compareJobs(oldJobs, newJobs api.Jobs) api.Jobs {
	var newJobsOnly api.Jobs
	for len(oldJobs) > 0 && len(newJobs) > 0 && oldJobs[len(oldJobs)-1].ID != newJobs[len(newJobs)-1].ID {
		poppedJob := newJobs[len(newJobs)-1]
		newJobs = newJobs[:len(newJobs)-1]

		newJobsOnly = append(newJobsOnly, poppedJob)
	}
	return newJobsOnly
}
