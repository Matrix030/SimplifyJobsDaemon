package simplifyapi

import (
	"encoding/json"
	"fmt"
	"io"
)

func (c *Client) GetJobData() (Jobs, error) {
	resp, err := c.httpClient.Get(URL)
	if err != nil {
		return Jobs{}, fmt.Errorf("Could not get the data from the URL %s", err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Jobs{}, err
	}

	jobResp := Jobs{}

	err = json.Unmarshal(data, &jobResp)
	if err != nil {
		return Jobs{}, err
	}

	return jobResp, nil

}
