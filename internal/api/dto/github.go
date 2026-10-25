package dto

import "time"

// Request

type CreateGithubInstallationStateRequest struct {
	Id        string    `json:"-"`
	CreatedAt time.Time `json:"-"`
	ExpiresAt time.Time `json:"-"`
}
