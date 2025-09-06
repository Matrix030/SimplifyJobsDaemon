package main

import (
	"time"

	"github.com/Matrix030/simplify_jobs_cli/internal/simplifyapi"
)

func main() {
	simplifyClient := simplifyapi.NewClient(5 * time.Minute)
	cfg := &clientConfig{
		simplifyapiClient: simplifyClient,
	}
}
