package main

import (
	api "github.com/Matrix030/simplify_jobs_cli/internal/simplifyapi"
)

type config struct {
	jobClient api.Client
}

// func startClient(cfg *config) {
//
// 	// ticker := time.NewTicker(30 * time.Minute)
// 	// defer ticker.Stop()
// 	//
// 	//running once on startup
// 	// fetchAndProcessJobs(cfg)
// 	// for {
// 	// 	select {
// 	// 	case <-ticker.C:
// 	// 		fetchAndProcessJobs(cfg)
// 	// 	}
// 	// }
// }

// func fetchAndProcessJobs(cfg *config) {
// 	jobs, err := cfg.jobClient.GetJobData()
// 	if err != nil {
// 		fmt.Println("There was an error %v\n", err)
// 		return
// 	}
// }
