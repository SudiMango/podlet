package dto

import "time"

// Request

type CreateGithubInstallationStateRequest struct {
	Id        string    `json:"-"`
	CreatedAt time.Time `json:"-"`
	ExpiresAt time.Time `json:"-"`
}

type UpsertGithubInstallationRequest struct {
	Id                  int64  `json:"-"`
	AccountId           int64  `json:"-"`
	AccountLogin        string `json:"-"`
	AccountType         string `json:"-"`
	RepositorySelection string `json:"-"`
}

// Response

type GetStateResponse struct {
	Id        string    `json:"-"`
	UserId    string    `json:"-"`
	CreatedAt time.Time `json:"-"`
	ExpiresAt time.Time `json:"-"`
}

type GithubInstallation struct {
	Id                  int64     `json:"-"`
	UserId              string    `json:"-"`
	AccountId           int64     `json:"-"`
	AccountLogin        string    `json:"-"`
	AccountType         string    `json:"-"`
	RepositorySelection string    `json:"-"`
	CreatedAt           time.Time `json:"-"`
	UpdatedAt           time.Time `json:"-"`
}

type GithubInstallationResponse struct {
	ID                  int64      `json:"id"`
	AppID               int64      `json:"app_id"`
	AppSlug             string     `json:"app_slug"`
	RepositorySelection string     `json:"repository_selection"`
	SuspendedAt         *time.Time `json:"suspended_at"`

	Account struct {
		ID    int64  `json:"id"`
		Login string `json:"login"`
		Type  string `json:"type"`
	} `json:"account"`

	Permissions struct {
		Contents string `json:"contents"`
		Metadata string `json:"metadata"`
	} `json:"permissions"`
}
