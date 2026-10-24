package auth

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/SudiMango/podlet/internal/api/auth"
	"github.com/SudiMango/podlet/internal/keyring"
	"github.com/spf13/cobra"
)

type LoginCommand struct{}

func (LoginCommand) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Login to Podlet",
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := login(cmd.Context(), cmd)
			if err != nil {
				return err
			}

			if err = keyring.SaveToken(token); err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Failed to save token: %v\n", err)
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Logged in successfully")
			return nil
		},
	}
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

	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "First copy your one-time code: %s\n", authorization.UserCode)
	fmt.Fprintf(out, "Press Enter to open %s in your browser... ", authorization.VerificationURI)

	if _, err := bufio.NewReader(cmd.InOrStdin()).ReadString('\n'); err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("wait for confirmation: %w", err)
	}

	if err := openBrowser(authorization.VerificationURI); err != nil {
		fmt.Fprintf(out, "Failed opening a web browser at %s\n", authorization.VerificationURI)
		fmt.Fprintf(out, "Open %s in your browser and enter code %s\n", authorization.VerificationURI, authorization.UserCode)
	}

	headers, err := postJSON(
		ctx,
		baseURL+"/auth/device/github/complete",
		auth.GithubDeviceCompleteRequest{
			DeviceCode: authorization.DeviceCode,
			ExpiresAt:  authorization.ExpiresAt,
			Interval:   authorization.Interval,
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

func openBrowser(url string) error {
	var command *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		command = exec.Command("open", url)
	case "windows":
		command = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		command = exec.Command("xdg-open", url)
	}

	return command.Run()
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
