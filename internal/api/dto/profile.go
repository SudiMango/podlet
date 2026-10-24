package dto

import "time"

// Request

// Response

type PodletProfileResponse struct {
	UserId    string    `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  *string   `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
