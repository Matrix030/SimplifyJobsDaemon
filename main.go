package main

import (
	"fmt"
	"reflect"
	"time"

	time_format "github.com/Matrix030/simplify_jobs_cli/internal"
	api "github.com/Matrix030/simplify_jobs_cli/internal/simplifyapi"
)

func main() {
	simplifyClient := api.NewClient(5 * time.Minute)
	cfg := &clientConfig{
		simplifyapiClient: simplifyClient,
	}

	var jobs api.Jobs
	jobs, err := cfg.simplifyapiClient.GetJobData()
	if err != nil {
		fmt.Printf("There was an error %v\n", err)

	}

	v := reflect.ValueOf(jobs[len(jobs)-1])
	t := reflect.TypeOf(jobs[len(jobs)-1])

	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Name
		fieldValue := v.Field(i).Interface()
		if fieldName == "DateUpdated" || fieldName == "DatePosted" {
			num := int64(fieldValue.(int))
			fieldValue = time_format.FormatUnixTime(num)
		}
		fmt.Printf("%s: %v\n", fieldName, fieldValue)
	}
}
