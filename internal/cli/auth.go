package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/SudiMango/podlet/internal/api/auth"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Podlet",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := login(cmd.Context(), cmd)
		if err != nil {
			return err
		}

		// TODO: Save this in the OS keyring.
		_ = token

		fmt.Fprintln(cmd.OutOrStdout(), "Logged in successfully")
		return nil
	},
}

func login(ctx context.Context, cmd *cobra.Command) (string, error) {
	baseURL := strings.TrimRight(os.Getenv("BASE_URL"), "/")
	if baseURL == "" {
		return "", errors.New("BASE_URL is required")
	}

	var authorization auth.GithubDeviceStartResponse
	if _, err := postJSON(
		ctx,
		baseURL+"/auth/device/github/start",
		nil,
		&authorization,
	); err != nil {
		return "", fmt.Errorf("start login: %w", err)
	}

	fmt.Fprintf(
		cmd.OutOrStdout(),
		"Open %s and enter code %s\n",
		authorization.VerificationURI,
		authorization.UserCode,
	)

	// This request blocks while the API polls GitHub for authorization.
	headers, err := postJSON(
		ctx,
		baseURL+"/auth/device/github/complete",
		map[string]any{
			"device_code": authorization.DeviceCode,
			"expires_at":  authorization.ExpiresAt,
			"interval":    authorization.Interval,
		},
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("complete login: %w", err)
	}

	token := headers.Get("Set-Auth-Token")
	if token == "" {
		return "", errors.New("API did not return a Podlet session token")
	}

	return token, nil
}

func postJSON(
	ctx context.Context,
	url string,
	input any,
	output any,
) (http.Header, error) {
	var body bytes.Buffer

	if input != nil {
		if err := json.NewEncoder(&body).Encode(input); err != nil {
			return nil, err
		}
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		&body,
	)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		var apiError struct {
			Message string `json:"message"`
		}
		if err := json.NewDecoder(response.Body).Decode(&apiError); err == nil &&
			apiError.Message != "" {
			return nil, errors.New(apiError.Message)
		}
		return nil, fmt.Errorf("API returned %s", response.Status)
	}

	if output != nil {
		if err := json.NewDecoder(response.Body).Decode(output); err != nil {
			return nil, err
		}
	}

	return response.Header, nil
}

func init() {
	authCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(authCmd)
}
