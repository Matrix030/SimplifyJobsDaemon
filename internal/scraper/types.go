package scraper

type JobDescription struct {
	JobID         string
	CompanyName   string
	Title         string
	URL           string
	Description   string
	ScrapedAt     int64
	ScrapeSuccess bool
	ErrorMessage  string
}
