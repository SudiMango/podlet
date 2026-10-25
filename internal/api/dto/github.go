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

type GithubInstallationState struct {
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
