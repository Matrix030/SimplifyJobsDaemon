package simplifyutils

import (
	"encoding/json"
	"fmt"
	"os"

	api "github.com/Matrix030/simplify_jobs_cli/internal/simplifyapi"
)

var fileName = "jobs.json"

func JsonFileWriter(jobsData api.Jobs) error {
	jsonData, err := json.MarshalIndent(jobsData, "", " ")
	if err != nil {
		fmt.Println("There was an error while marshalling data", err)
		return err
	}

	err = os.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		fmt.Println("There was an error which writing the json file", err)
		return err
	}

	fmt.Println("Successfully wrote the jobs to the jobs.json")

	return nil

}

func LoadExistingJobs() (api.Jobs, error) {
	fileData := &api.Jobs{}
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		fmt.Println("File does not exists")
		return api.Jobs{}, err
	}

	//read file content
	data, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("There was an error while reading the file")
		return api.Jobs{}, err
	}

	err = json.Unmarshal(data, &fileData)

	return *fileData, nil
}
