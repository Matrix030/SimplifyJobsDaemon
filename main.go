package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	api "github.com/Matrix030/SimplifyJobsDaemon/internal/simplifyapi"
)

func main() {
	fmt.Println("=== SimplifyJobs CLI Monitor ===")
	fmt.Println("Monitoring new grad positions...")
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println()

	// Create HTTP client with reasonable timeout
	simplifyClient := api.NewClient(5 * time.Minute)
	cfg := &config{
		jobClient: simplifyClient,
	}

	// Set up graceful shutdown handling
	// This catches Ctrl+C (SIGINT) and SIGTERM signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the monitoring in a goroutine so we can handle shutdown signals
	done := make(chan struct{})
	go func() {
		defer close(done)
		startClient(cfg)
	}()

	// Wait for either the client to finish or a shutdown signal
	select {
	case <-sigChan:
		fmt.Println("\n🛑 Shutdown signal received. Stopping job monitor...")
		fmt.Println("Thank you for using SimplifyJobs CLI!")
	case <-done:
		fmt.Println("Job monitor stopped.")
	}
}
