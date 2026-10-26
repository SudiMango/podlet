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
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	ServerID       string    `json:"server_id"`
	JobStatus      JobStatus `json:"job_status"`
	JobType        JobType   `json:"job_type"`
	RepoName       string    `json:"repo_name"`
	RepositoryID   int64     `json:"repository_id"`
	InstallationID int64     `json:"installation_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
