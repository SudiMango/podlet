package client

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
		return nil, apiError(response)
	}

	if output != nil {
		if err := json.NewDecoder(response.Body).Decode(output); err != nil {
			return nil, fmt.Errorf("decode response: %w", err)
		}
	}

	return response.Header, nil
}

func (c *Client) Stream(
	ctx context.Context,
	path string,
	token string,
	onData func(data []byte) (done bool, err error),
) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	request.Header.Set("Accept", "text/event-stream")
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return apiError(response)
	}

	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		data, ok := strings.CutPrefix(scanner.Text(), "data: ")
		if !ok {
			continue // blank separators and ": ping" keepalives
		}

		if done, err := onData([]byte(data)); err != nil || done {
			return err
		}
	}

	return scanner.Err()
}

func apiError(response *http.Response) error {
	var body struct {
		Message string `json:"message"`
	}

	if err := json.NewDecoder(response.Body).Decode(&body); err == nil && body.Message != "" {
		return fmt.Errorf("%s", body.Message)
	}

	return fmt.Errorf("API returned %s", response.Status)
}
