package contract

import "time"

// Request

type CreateServerRequest struct {
	Name              string `json:"name"`
	RegistrationToken string `json:"registration_token"`
}

type UpdateServerRequest struct {
	Name string `json:"name"`
}

// Response

type PendingServerResponse struct {
	Id                string    `json:"id"`
	UserId            string    `json:"user_id"`
	CreatedAt         time.Time `json:"created_at"`
	ExpiresAt         time.Time `json:"expires_at"`
	RegistrationToken *string   `json:"registration_token"`
}

type PendingServersResponse []PendingServerResponse

type ServerResponse struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	UserId    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Token     *string   `json:"token"`
}

type ServersResponse []ServerResponse
