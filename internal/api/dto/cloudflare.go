package dto

import "time"

// Request

type CreateDNSRecordRequest struct {
	TunnelId    string    `json:"tunnel_id"`
	Service     string    `json:"service"`
	DnsRecordId string    `json:"-"`
	Hostname    string    `json:"-"`
	CreatedAt   time.Time `json:"-"`
}

// Response
