package contract

import "time"

// Response

type GithubInstallationStartResponse struct {
	InstallationUrl string    `json:"installation_url"`
	ExpiresAt       time.Time `json:"expires_at"`
}

type GithubRepositoryResponse struct {
	Id             int64  `json:"id"`
	InstallationId int64  `json:"installation_id"`
	Name           string `json:"name"`
	FullName       string `json:"full_name"`
	Private        bool   `json:"private"`
	HtmlUrl        string `json:"html_url"`
	CloneUrl       string `json:"clone_url"`
	DefaultBranch  string `json:"default_branch"`
}

type GithubRepositoriesResponse []GithubRepositoryResponse
