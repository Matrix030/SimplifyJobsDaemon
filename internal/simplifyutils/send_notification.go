package simplifyutils

import (
	"fmt"
	api "github.com/Matrix030/simplify_jobs_cli/internal/simplifyapi"
	"os/exec"
)

func SendNotification(jobSlice api.Jobs) error {

	if len(jobSlice) == 0 {
		return nil
	}
	if len(jobSlice) > 3 {
		numJobs := fmt.Sprintf("%v New Job Notification", len(jobSlice))
		cmd := exec.Command("notify-send", "New Job Notifications", numJobs)
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("An Error occurred while sending the a notification: %v\n", err)
		}
		return nil
	}

	cmd := exec.Command("notify-send", "New Job Notification", jobSlice[0].CompanyName)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("An Error occurred while sending a notification: %v \n", err)
	}
	return nil
}
