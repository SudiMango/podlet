package contract

import "time"

// Request

type ReportJobRequest struct {
	JobId     string    `json:"job_id"`
	JobStatus JobStatus `json:"job_status"`
	Error     string    `json:"error"`
}

// Response

type ClaimJobResponse struct {
	JobId       string  `json:"id"`
	JobType     JobType `json:"type"`
	CloneUrl    string  `json:"clone_url"`
	GithubToken string  `json:"github_token"`
}

type JobResponse struct {
	Id           string    `json:"id"`
	ServerId     string    `json:"server_id"`
	JobStatus    JobStatus `json:"job_status"`
	ErrorMessage *string   `json:"error,omitempty"`
	JobType      JobType   `json:"job_type"`
	RepoName     string    `json:"repo_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type JobsResponse []JobResponse
