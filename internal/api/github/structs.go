package github

import "time"

// Request

type createInstallationTokenRequest struct {
	RepositoryIDs []int64           `json:"repository_ids,omitempty"`
	Permissions   map[string]string `json:"permissions,omitempty"`
}

// Response

type InstallationToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Repository struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	Private       bool   `json:"private"`
	HtmlUrl       string `json:"html_url"`
	CloneUrl      string `json:"clone_url"`
	DefaultBranch string `json:"default_branch"`
}

type RepositoriesResponse struct {
	TotalCount   int          `json:"total_count"`
	Repositories []Repository `json:"repositories"`
}

type Installation struct {
	Id                  int64      `json:"id"`
	AppId               int64      `json:"app_id"`
	AppSlug             string     `json:"app_slug"`
	RepositorySelection string     `json:"repository_selection"`
	SuspendedAt         *time.Time `json:"suspended_at"`

	Account Account `json:"account"`

	Permissions Permissions `json:"permissions"`
}

type Account struct {
	Id    int64  `json:"id"`
	Login string `json:"login"`
	Type  string `json:"type"`
}

type Permissions struct {
	Contents string `json:"contents"`
	Metadata string `json:"metadata"`
}
