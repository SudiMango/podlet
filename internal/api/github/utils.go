package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func (c *Client) doJSON(
	ctx context.Context,
	method string,
	path string,
	bearerToken string,
	input any,
	expectedStatus int,
	output any,
) error {
	var body bytes.Buffer

	if input != nil {
		if err := json.NewEncoder(&body).Encode(input); err != nil {
			return fmt.Errorf("encode GitHub request: %w", err)
		}
	}

	request, err := http.NewRequestWithContext(
		ctx,
		method,
		apiBaseURL+path,
		&body,
	)
	if err != nil {
		return fmt.Errorf("create GitHub request: %w", err)
	}

	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Authorization", "Bearer "+bearerToken)
	request.Header.Set("X-GitHub-Api-Version", apiVersion)

	if input != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("send GitHub request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != expectedStatus {
		message, _ := io.ReadAll(
			io.LimitReader(response.Body, 8<<10),
		)

		return fmt.Errorf(
			"GitHub returned %s: %s",
			response.Status,
			message,
		)
	}

	if output == nil {
		return nil
	}

	if err := json.NewDecoder(
		io.LimitReader(response.Body, maxResponseSize),
	).Decode(output); err != nil {
		return fmt.Errorf("decode GitHub response: %w", err)
	}

	return nil
}

func (c *Client) generateJWT() (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    c.clientId,
		IssuedAt:  jwt.NewNumericDate(now.Add(-60 * time.Second)),
		ExpiresAt: jwt.NewNumericDate(now.Add(9 * time.Minute)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(c.privateKey)
	if err != nil {
		return "", fmt.Errorf("sign GitHub App JWT: %w", err)
	}

	return signed, nil
}
