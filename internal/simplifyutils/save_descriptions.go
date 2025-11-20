package simplifyutils

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Matrix030/SimplifyJobsDaemon/internal/scraper"
)

var descriptionsFileName = "job_descriptions.json"

func SaveJobDescriptions(descriptions []scraper.JobDescription) error {
	if len(descriptions) == 0 {
		return nil
	}

	var exisiting []scraper.JobDescription

	if data, err := os.ReadFile(descriptionsFileName); err == nil && len(data) > 0 {
		json.Unmarshal(data, &exisiting)
	}

	combined := append(descriptions, exisiting...)
	jsonData, err := json.MarshalIndent(combined, "", " ")
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	err = os.WriteFile(descriptionsFileName, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("write error: %w", err)
	}

	fmt.Printf("Saved %d job descriptions\n", len(descriptions))
	return nil

}
