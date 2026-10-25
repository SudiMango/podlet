package contract

import "time"

// Response

type GithubInstallationStartResponse struct {
	InstallationUrl string    `json:"installation_url"`
	ExpiresAt       time.Time `json:"expires_at"`
}
