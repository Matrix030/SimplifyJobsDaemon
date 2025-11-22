package main

import (
	"fmt"
	"time"

	"github.com/Matrix030/SimplifyJobsDaemon/internal/ai"
	"github.com/Matrix030/SimplifyJobsDaemon/internal/resume"
)

func main() {
	//Create LLM client
	client := ai.NewLLMClient("http://localhost:5000", 2*time.Minute)

	//Health check
	fmt.Println("Checking LLM server health...")
	if err := client.HealthCheck(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("LLM  server is healthy\n")

	//Load projects
	projects, err := ai.LoadProjectsFromFile("../../projects.json")
	if err != nil {
		fmt.Printf("Error loading projects: %v\n", err)
		return
	}
	fmt.Printf("Loaded %d projects\n\n", len(projects))

	//Test analysis
	jobDesc := "We're looking for a backend engineer with strong Go experience, REST API design, and experience building scalable microservices. Bonus points for system monitoring and notification systems."

	fmt.Println("Analyzing job...")
	result, err := client.AnalyzeJob(jobDesc, projects)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("\nSelected Projects: %v\n", result.SelectedProjects)
	fmt.Printf("Reasoning: %s\n", result.Reasoning)

	//Create resume editor
	fmt.Println("\n--- Testing Resume Generation ---")
	editor := resume.NewEditor(
		"../../scripts/edit_resume.py",
		"../../scripts/resume_template.odt",
		"../../projects.json",
		"../../tailored_resumes",
	)

	//Generate tailored resume
	outputPath, err := editor.TailorResume(result.SelectedProjects, "test_output.pdf")
	if err != nil {
		fmt.Printf("Error generating resume: %v\n", err)
		return
	}

	fmt.Printf("\nSuccess! Generated resume at: %s\n", outputPath)
}
