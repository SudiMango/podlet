package contract

import "time"

// Request

// Response

type CreateTunnelResponse struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type GetIngressConfigsResponse struct {
	Version int           `json:"version"`
	Config  IngressConfig `json:"config"`
}

type IngressConfig struct {
	Ingress []Ingress `json:"ingress"`
}

type Ingress struct {
	Hostname *string `json:"hostname,omitempty"`
	Service  string  `json:"service"`
}

type CreateDNSRecordResponse struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_on"`
}
