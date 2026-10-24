package auth

import (
	"time"
)

// Request
type GithubDeviceCompleteRequest struct {
	DeviceCode string    `json:"device_code"`
	ExpiresAt  time.Time `json:"expires_at"`
	Interval   int64     `json:"interval"`
}

// Response

type GithubDeviceStartResponse struct {
	DeviceCode              string    `json:"device_code"`
	UserCode                string    `json:"user_code"`
	VerificationURI         string    `json:"verification_uri"`
	VerificationURIComplete string    `json:"verification_uri_complete,omitempty"`
	ExpiresAt               time.Time `json:"expires_at"`
	Interval                int64     `json:"interval"`
}
