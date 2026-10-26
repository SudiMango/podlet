package contract

// Request

type CompleteJobRequest struct {
	JobId   string `json:"job_id"`
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// Response

type ClaimJobResponse struct {
	JobId       string  `json:"id"`
	JobType     JobType `json:"type"`
	Repository  string  `json:"repository"`
	CloneUrl    string  `json:"clone_url"`
	Branch      string  `json:"branch"`
	GithubToken string  `json:"github_token"`
}
