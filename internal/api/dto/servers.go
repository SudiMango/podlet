package dto

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
