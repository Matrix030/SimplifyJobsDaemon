package main

import (
	"fmt"
	time_format "github.com/Matrix030/simplify_jobs_cli/internal"
	api "github.com/Matrix030/simplify_jobs_cli/internal/simplifyapi"
	"reflect"
)

type config struct {
	jobClient api.Client
}

func startClient(cfg *config) {

	var jobs api.Jobs
	jobs, err := cfg.jobClient.GetJobData()
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
