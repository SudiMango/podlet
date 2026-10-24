package auth

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/SudiMango/podlet/internal/cli/browser"
	"github.com/SudiMango/podlet/internal/cli/client"
	"github.com/SudiMango/podlet/internal/contract"
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

			if err := keyring.SaveToken(token); err != nil {
				return fmt.Errorf("save session token: %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Logged in successfully")
			return nil
		},
	}
}

func login(ctx context.Context, cmd *cobra.Command) (string, error) {
	api, err := client.New()
	if err != nil {
		return "", err
	}

	var authorization contract.GithubDeviceStartResponse
	if _, err := api.DoJSON(
		ctx,
		http.MethodPost,
		"/auth/device/github/start",
		nil,
		&authorization,
		"",
	); err != nil {
		return "", fmt.Errorf("start login: %w", err)
	}

	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "First copy your one-time code: %s\n", authorization.UserCode)
	fmt.Fprintf(out, "Press Enter to open %s in your browser... ", authorization.VerificationURI)

	if _, err := bufio.NewReader(cmd.InOrStdin()).ReadString('\n'); err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("wait for confirmation: %w", err)
	}

	if err := browser.Open(authorization.VerificationURI); err != nil {
		fmt.Fprintf(out, "Failed opening a web browser at %s\n", authorization.VerificationURI)
		fmt.Fprintf(out, "Open %s in your browser and enter code %s\n", authorization.VerificationURI, authorization.UserCode)
	}

	headers, err := api.DoJSON(
		ctx,
		http.MethodPost,
		"/auth/device/github/complete",
		contract.GithubDeviceCompleteRequest{
			DeviceCode: authorization.DeviceCode,
			ExpiresAt:  authorization.ExpiresAt,
			Interval:   authorization.Interval,
		},
		nil,
		"",
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
