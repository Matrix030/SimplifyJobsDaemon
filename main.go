package main

import (
	api "github.com/Matrix030/simplify_jobs_cli/internal/simplifyapi"
	"time"
)

func main() {
	simplifyClient := api.NewClient(5 * time.Minute)
	cfg := &config{
		jobClient: simplifyClient,
	}

	startClient(cfg)

}
