package contract

import "time"

// Request

// Response

type CreateTunnelResponse struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateDNSRecordResponse struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_on"`
}

type TunnelResponse struct {
	Id             string    `json:"id"`
	Name           string    `json:"name"`
	UserId         string    `json:"user_id"`
	ServerId       string    `json:"server_id"`
	CreatedAt      time.Time `json:"created_at"`
	ConnectorToken *string   `json:"connector_token,omitempty"`
}

type TunnelsResponse []TunnelResponse

type DnsRecordResponse struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Service   string    `json:"service"`
	UserId    string    `json:"user_id"`
	ServerId  string    `json:"server_id"`
	TunnelId  string    `json:"tunnel_id"`
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
