package cloudflare

import (
	"github.com/SudiMango/podlet/internal/contract"
)

// Cloudflare specific

type CloudflareResponse[T any] struct {
	Success  bool              `json:"success"`
	Errors   []CloudflareError `json:"errors"`
	Messages []CloudflareError `json:"messages"`
	Result   T                 `json:"result"`
}

type CloudflareError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Request

type CreateTunnelRequest struct {
	Name      string `json:"name"`
	ConfigSrc string `json:"config_src"`
}

type PutIngressRequest struct {
	Config contract.IngressConfig `json:"config"`
}

type CreateDNSRecordRequest struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Proxied bool   `json:"proxied"`
	TTL     int    `json:"ttl"`
	Comment string `json:"comment"`
}
