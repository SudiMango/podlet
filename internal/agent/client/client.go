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

type APIError struct {
	StatusCode int
	Message    string
}

func NewClient(baseURL string, token string) *Client {
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
		var body struct {
			Message string `json:"message"`
		}

		_ = json.NewDecoder(response.Body).Decode(&body)

		return &APIError{
			StatusCode: response.StatusCode,
			Message:    body.Message,
		}
	}

	if output != nil && response.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(response.Body).Decode(output); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return fmt.Sprintf(
		"API returned HTTP %d",
		e.StatusCode,
	)
}

func IsStatus(err error, status int) bool {
	var apiError *APIError

	return errors.As(err, &apiError) &&
		apiError.StatusCode == status
}
