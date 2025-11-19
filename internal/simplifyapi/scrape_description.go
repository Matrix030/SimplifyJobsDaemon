package simplifyapi

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

func (c *Client) ScrapeJobDescription(url string) (string, error) {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	description := doc.Find("meta[property='og:description']").AttrOr("content", "")

	if description == "" {
		description = doc.Find("meta[name='description']").AttrOr("content", "")
	}

	return description, nil
}
