package main

import "github.com/Matrix030/simplify_jobs_cli/internal/simplifyapi"

type clientConfig struct {
	simplifyapiClient simplifyapi.Client
}

type jobConfig struct {
	sponsorshipType bool
	locations       []string
	category        string
}
