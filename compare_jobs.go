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

// NEW: Scrape descriptions for new jobs
func scrapeNewJobDescriptions(jobs api.Jobs) []scraper.JobDescription {
	s := scraper.NewScraper(10 * time.Second)
	var descriptions []scraper.JobDescription

	for _, job := range jobs {
		fmt.Printf("Scraping: %s - %s\n", job.CompanyName, job.Title)

		desc := s.ScrapeJobDescription(job.URL, job.ID, job.CompanyName, job.Title)
		descriptions = append(descriptions, desc)

		// Rate limiting - be respectful to servers
		time.Sleep(2 * time.Second)
	}

	return descriptions
}

func isEligibleSponsorship(sponsorship string) bool {
	switch sponsorship {
	case "Other", "Offers Sponsorship":
		return true
	case "U.S. Citizenship is Required", "Does Not Offer Sponsorship":
		return false
	default:
		return true
	}
}
