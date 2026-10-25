package contract

import "time"

// Response

type CreateGithubInstallationStateResponse struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type GithubInstallationStartResponse struct {
	InstallationUrl string    `json:"installation_url"`
	ExpiresAt       time.Time `json:"expires_at"`
}
