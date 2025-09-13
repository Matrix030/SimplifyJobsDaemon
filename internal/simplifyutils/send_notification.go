package simplifyutils

import (
	"fmt"
	api "github.com/Matrix030/SimplifyJobsDaemon/internal/simplifyapi"
	"os/exec"
	"strings"
)

func SendNotification(jobSlice api.Jobs) error {
	if len(jobSlice) == 0 {
		return nil
	}

	var title, body string

	if len(jobSlice) == 1 {
		// Single job notification - show company and title
		job := jobSlice[0]
		title = "New Job Available!"
		body = fmt.Sprintf("%s - %s", job.CompanyName, job.Title)
	} else if len(jobSlice) <= 3 {
		// Few jobs - show count and list companies
		title = fmt.Sprintf("%d New Jobs Available!", len(jobSlice))
		var companies []string
		for _, job := range jobSlice {
			companies = append(companies, job.CompanyName)
		}
		body = strings.Join(companies, ", ")
	} else {
		// Many jobs - just show count
		title = fmt.Sprintf("%d New Jobs Available!", len(jobSlice))
		body = "Check your terminal for details"
	}

	// Use notify-send command (Linux desktop notification)
	cmd := exec.Command("notify-send", "-u", "critical", title, body)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	return nil
}
