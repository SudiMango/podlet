package cloudflare

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	apiBaseURL      = "https://api.cloudflare.com/client/v4"
	maxResponseSize = 1 << 20
)

type Client struct {
	httpClient *http.Client
	apiToken   string
	accountId  string
	zoneId     string
}

func NewClient() (*Client, error) {
	apiToken := os.Getenv("CLOUDFLARE_API_TOKEN")
	if apiToken == "" {
		return nil, fmt.Errorf("CLOUDFLARE_API_TOKEN is required")
	}

	accountId := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	if accountId == "" {
		return nil, fmt.Errorf("CLOUDFLARE_ACCOUNT_ID is required")
	}

	zoneId := os.Getenv("CLOUDFLARE_ZONE_ID")
	if zoneId == "" {
		return nil, fmt.Errorf("CLOUDFLARE_ZONE_ID is required")
	}

	return &Client{
		apiToken:  apiToken,
		accountId: accountId,
		zoneId:    zoneId,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}
