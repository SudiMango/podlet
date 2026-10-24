package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/SudiMango/podlet/internal/api/dto"
	"github.com/SudiMango/podlet/internal/keyring"
	"github.com/spf13/cobra"
)

type MeCommand struct{}

func (MeCommand) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "me",
		Short: "Show the currently authenticated user",
		RunE: func(cmd *cobra.Command, args []string) error {
			profile, err := getCurrentProfile(cmd)
			if err != nil {
				return err
			}

			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetIndent("", "  ")
			return encoder.Encode(profile)
		},
	}
}

func getCurrentProfile(cmd *cobra.Command) (*dto.PodletProfileResponse, error) {
	token, err := keyring.LoadToken()
	if err != nil {
		return nil, fmt.Errorf(
			"load session token: %w; run `podlet auth login`",
			err,
		)
	}

	baseURL := strings.TrimRight(os.Getenv("BASE_URL"), "/")
	if baseURL == "" {
		return nil, errors.New("BASE_URL is required")
	}

	request, err := http.NewRequestWithContext(
		cmd.Context(),
		http.MethodGet,
		baseURL+"/profile",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("create profile request: %w", err)
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("get profile: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		return nil, errors.New(
			"session is invalid or expired; run `podlet auth login`",
		)
	}

	if response.StatusCode != http.StatusOK {
		var apiError struct {
			Message string `json:"message"`
		}

		if err := json.NewDecoder(response.Body).Decode(&apiError); err == nil &&
			apiError.Message != "" {
			return nil, errors.New(apiError.Message)
		}

		return nil, fmt.Errorf("API returned %s", response.Status)
	}

	var profile dto.PodletProfileResponse
	if err := json.NewDecoder(response.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("decode profile response: %w", err)
	}

	return &profile, nil
}
