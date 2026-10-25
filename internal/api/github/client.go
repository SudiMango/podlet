package github

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	apiBaseURL      = "https://api.github.com"
	apiVersion      = "2026-03-10"
	maxResponseSize = 1 << 20
)

type Client struct {
	clientId   string
	privateKey *rsa.PrivateKey
	httpClient *http.Client
}

func NewClient() (*Client, error) {
	clientId := os.Getenv("GITHUB_CLIENT_ID")
	if clientId == "" {
		return nil, fmt.Errorf("GITHUB_CLIENT_ID is required")
	}

	privateKeyPath := os.Getenv("GITHUB_PRIVATE_KEY")
	if privateKeyPath == "" {
		return nil, fmt.Errorf("GITHUB_PRIVATE_KEY is required")
	}

	privateKeyPEM, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("read GitHub private key: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("parse GitHub private key: %w", err)
	}

	return &Client{
		clientId:   clientId,
		privateKey: privateKey,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}, nil
}
