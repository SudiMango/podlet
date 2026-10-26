package contract

import (
	"time"
)

// Request

type CreateServerRequest struct {
	Name              string `json:"name"`
	RegistrationToken string `json:"registration_token"`
}

type UpdateServerRequest struct {
	Name *string `json:"name"`
}

type CreateJobRequest struct {
	ServerId string `json:"server_id"`
	RepoName string `json:"repo_name"`
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
	Id             string               `json:"id"`
	Name           string               `json:"name"`
	UserId         string               `json:"user_id"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
	Token          *string              `json:"token"`
	LastSeenAt     *time.Time           `json:"last_seen_at"`
	LifecycleState ServerLifecycleState `json:"lifecycle_state"`
}

type ServersResponse []ServerResponse

type HeartbeatResponse struct {
	Action *HeartbeatAction `json:"action"`
}

type CreateJobResponse struct {
	JobId  string    `json:"job_id"`
	Status JobStatus `json:"status"`
}
