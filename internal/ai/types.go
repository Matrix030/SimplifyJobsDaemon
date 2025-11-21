package ai

type Project struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AnalysisRequest struct {
	JobDescription string    `json:"job_description"`
	Project        []Project `json:"projects"`
}

type AnalysisResponse struct {
	SelectedProjects []string `json:"selected_projects"`
	Reasoning        string   `json:"reasoning"`
}

type LLMError struct {
	Message string
	Err     error
}

func (e *LLMError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}

	return e.Message
}
