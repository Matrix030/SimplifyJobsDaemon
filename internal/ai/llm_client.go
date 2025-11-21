package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LLMClient struct {
	baseURL    string
	httpClient *http.Client
}

func newLLMClient(baseURL string, timeout time.Duration) *LLMClient {
	return &LLMClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// HealthCheck verifies the LLM sever is running
func (c *LLMClient) HealthCheck() error {
	resp, err := c.httpClient.Get(c.baseURL + "/health")
	if err != nil {
		return &LLMError{
			Message: "Failed to connect to LLM server",
			Err:     err,
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return &LLMError{
			Message: fmt.Sprintf("LLM server unhealthy (status %d): %s", resp.StatusCode, string(body)),
		}
	}

	return nil
}

// AnalyzeJob sends job description and projects to LLMfor analysis
func (c *LLMClient) AnalyzeJob(jobDescription string, projects []Project) (*AnalysisResponse, error) {
	//Create request payload
	reqPayload := AnalysisRequest{
		JobDescription: jobDescription,
		Project:        projects,
	}

	jsonData, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, &LLMError{
			Message: "Failed to marshal request",
			Err:     err,
		}
	}

	// Send request to Flask Server
	resp, err := c.httpClient.Post(
		c.baseURL+"/analyze",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, &LLMError{
			Message: "Failed to send request to LLM server",
			Err:     err,
		}
	}
	defer resp.Body.Close()

	//Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &LLMError{
			Message: "Failed to read response",
			Err:     err,
		}
	}

	//Check for error status codes
	if resp.StatusCode != http.StatusOK {
		return nil, &LLMError{
			Message: fmt.Sprintf("LLM server returned error (status %d): %s", resp.StatusCode, string(body)),
		}
	}

	//Parse response
	var analysisResp AnalysisResponse
	err = json.Unmarshal(body, &analysisResp)
	if err != nil {
		return nil, &LLMError{
			Message: "Failed to parse LLM response",
			Err:     err,
		}
	}

	return &analysisResp, nil
}
