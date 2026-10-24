package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SudiMango/podlet/internal/cli/config"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func New() (*Client, error) {
	baseURL, err := config.LoadBaseURL()
	if err != nil {
		return nil, err
	}

	return &Client{
		baseURL:    baseURL,
		httpClient: http.DefaultClient,
	}, nil
}

func (c *Client) DoJSON(
	ctx context.Context,
	method string,
	path string,
	input any,
	output any,
	token string,
) (http.Header, error) {
	var body bytes.Buffer

	if input != nil {
		if err := json.NewEncoder(&body).Encode(input); err != nil {
			return nil, fmt.Errorf("encode request: %w", err)
		}
	}

	request, err := http.NewRequestWithContext(
		ctx,
		method,
		c.baseURL+path,
		&body,
	)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		var apiError struct {
			Message string `json:"message"`
		}

		if err := json.NewDecoder(response.Body).Decode(&apiError); err == nil &&
			apiError.Message != "" {
			return nil, fmt.Errorf("%s", apiError.Message)
		}

		return nil, fmt.Errorf("API returned %s", response.Status)
	}

	if output != nil {
		if err := json.NewDecoder(response.Body).Decode(output); err != nil {
			return nil, fmt.Errorf("decode response: %w", err)
		}
	}

	return response.Header, nil
}
