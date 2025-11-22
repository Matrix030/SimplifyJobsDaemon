package main

import (
	"fmt"
	ai "github.com/Matrix030/SimplifyJobsDaemon/internal/ai"
	"github.com/Matrix030/SimplifyJobsDaemon/internal/resume"
	api "github.com/Matrix030/SimplifyJobsDaemon/internal/simplifyapi"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Printf("=== SimplifyJobs CLI Monitor ===")
	fmt.Println("Monitoring new grad positions...")
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println()

	//Create HTTP client
	simplifyClient := api.NewClient(5 * time.Minute)

	//Create LLM client
	llmClient := ai.NewLLMClient("http://localhost:5000", 2*time.Minute)

	//Load Projects
	projects, err := ai.LoadProjectsFromFile("projects.json")
	if err != nil {
		fmt.Printf("Warning: Could not load  projects: %v\n", err)
		fmt.Printf("LLM analysis will be disabled,")
		projects = nil
	} else {
		fmt.Printf("Loaded %d projects for LLM analysis\n", len(projects))
	}

	resumeEditor := resume.NewEditor(
		"scripts/edit_resume.py",
		"scripts/resume_template.odt",
		"projects.json",
		"tailored_resumes",
	)

	cfg := &config{
		jobClient:    simplifyClient,
		llmClient:    llmClient,
		projects:     projects,
		resumeEditor: resumeEditor,
	}

	//Signal Handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})
	go func() {
		defer close(done)
		startClient(cfg)
	}()

	select {
	case <-sigChan:
		fmt.Println("\nShutdown signal received. Stopping job monitor...")
		fmt.Println("Thank you for using SimplifyJobs CLI!")

	case <-done:
		fmt.Println("Job monitor stopped.")
	}
}
