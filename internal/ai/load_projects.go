package ai

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadProjectsFromFile reads projects from a JSON file
func LoadProjectsFromFile(filepath string) ([]Project, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read projects file: %w", err)
	}

	var projects []Project
	err = json.Unmarshal(data, &projects)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse projects JSON: %w", err)
	}

	if len(projects) == 0 {
		return nil, fmt.Errorf("projects file is empty")
	}

	return projects, nil
}
