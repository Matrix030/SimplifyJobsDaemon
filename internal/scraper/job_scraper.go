package scraper

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Scraper struct {
	httpClient *http.Client
}

func NewScraper(timeout time.Duration) *Scraper {
	return &Scraper{
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (s *Scraper) ScrapeJobDescription(url, jobID, companyName, title string) JobDescription {
	result := JobDescription{
		JobID:       jobID,
		CompanyName: companyName,
		Title:       title,
		URL:         url,
		ScrapedAt:   time.Now().Unix(),
	}

	resp, err := s.httpClient.Get(url)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("HTTP error: %v", err)
		return result
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		result.ErrorMessage = fmt.Sprintf("Status code %d", resp.StatusCode)
		return result
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("Parse error: %v", err)
		return result
	}

	//Try multiple types of description
	description := s.extractdescription(doc)

	if description == "" {
		result.ErrorMessage = "No description found"
		return result
	}

	result.Description = cleanText(description)
	result.ScrapeSuccess = true
	return result
}

func (s *Scraper) extractdescription(doc *goquery.Document) string {
	//strategy one
	if desc := doc.Find("meta[property='og:description']").AttrOr("content", ""); desc != "" {
		return desc
	}

	//strategy two
	if desc := doc.Find("meta[name='description']").AttrOr("content", ""); desc != "" {
		return desc
	}

	selectors := []string{
		".job-description",
		".description",
		"[class*='job-detail']",
		"[class*='description']",
		"article",
		"main",
	}

	for _, selector := range selectors {
		if text := doc.Find(selector).First().Text(); text != "" {
			return text
		}
	}

	return ""
}

func cleanText(text string) string {
	text = strings.TrimSpace(text)

	lines := strings.Fields(text)
	return strings.Join(lines, " ")
}
