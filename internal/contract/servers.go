package contract

import "time"

// Request

type CreatePendingServerRequest struct {
	Id                    string    `json:"-"`
	CreatedAt             time.Time `json:"-"`
	ExpiresAt             time.Time `json:"-"`
	RegistrationTokenHash string    `json:"-"`
}

type CreateServerRequest struct {
	Id        string    `json:"-"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	TokenHash string    `json:"-"`
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
	TokenHash *string   `json:"-"`
}

type ServersResponse []ServerResponse
