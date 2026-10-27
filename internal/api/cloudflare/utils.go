package cloudflare

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
)

func (c *Client) doJSON(
	ctx context.Context,
	method string,
	path string,
	input any,
	expectedStatus int,
	output any,
) error {
	var body bytes.Buffer

	if input != nil {
		if err := json.NewEncoder(&body).Encode(input); err != nil {
			return fmt.Errorf("encode Cloudflare request: %w", err)
		}
	}

	request, err := http.NewRequestWithContext(
		ctx,
		method,
		apiBaseURL+path,
		&body,
	)
	if err != nil {
		return fmt.Errorf("create Cloudflare request: %w", err)
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+c.apiToken)
	if input != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("send Cloudflare request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != expectedStatus {
		message, _ := io.ReadAll(io.LimitReader(response.Body, 8<<10))
		return fmt.Errorf("Cloudflare returned %s: %s", response.Status, message)
	}

	if output == nil {
		return nil
	}

	if err := json.NewDecoder(io.LimitReader(response.Body, maxResponseSize)).Decode(output); err != nil {
		return fmt.Errorf("decode Cloudflare response: %w", err)
	}

	return nil
}

func generateCode(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

	code := make([]byte, length)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[n.Int64()]
	}

	return string(code), nil
}

func validateResponse[T any](response CloudflareResponse[T]) error {
	if response.Success {
		return nil
	}

	if len(response.Errors) > 0 {
		return fmt.Errorf("Cloudflare returned error %d: %s", response.Errors[0].Code, response.Errors[0].Message)
	}

	return fmt.Errorf("Cloudflare request failed")
}
