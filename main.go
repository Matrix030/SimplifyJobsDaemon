package main

import "net/http"


URL := "https://raw.githubusercontent.com/SimplifyJobs/New-Grad-Positions/refs/heads/dev/.github/scripts/listings.json"
type jobStruct []struct {
		Source      string   `json:"source"`
		CompanyName string   `json:"company_name"`
		ID          string   `json:"id"`
		Title       string   `json:"title"`
		Active      bool     `json:"active"`
		DateUpdated int      `json:"date_updated"`
		DatePosted  int      `json:"date_posted"`
		URL         string   `json:"url"`
		Locations   []string `json:"locations"`
		CompanyURL  string   `json:"company_url"`
		IsVisible   bool     `json:"is_visible"`
		Sponsorship string   `json:"sponsorship"`
		Category    string   `json:"category,omitempty"`
		Degrees     []any    `json:"degrees,omitempty"`
	}

func main() {
	//TODO- get JSON from URL
	//PARSE the JSON in a struct


		 
}
