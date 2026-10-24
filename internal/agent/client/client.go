package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func New(baseURL string, token string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   token,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *Client) DoJSON(
	ctx context.Context,
	method string,
	path string,
	input any,
	output any,
) error {
	var body bytes.Buffer

	if input != nil {
		if err := json.NewEncoder(&body).Encode(input); err != nil {
			return fmt.Errorf("encode request: %w", err)
		}
	}

	request, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, &body)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+c.token)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		var apiError struct {
			Message string `json:"message"`
		}

		if err := json.NewDecoder(response.Body).Decode(&apiError); err == nil && apiError.Message != "" {
			return errors.New(apiError.Message)
		}

		return fmt.Errorf("API returned %s", response.Status)
	}

	if output != nil && response.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(response.Body).Decode(output); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}
