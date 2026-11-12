package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SudiMango/podlet/internal/cli/util"
	"github.com/SudiMango/podlet/internal/contract"
	"github.com/spf13/cobra"
)

type MeCommand struct{}

func (MeCommand) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "me",
		Short: "Show the currently authenticated user",
		RunE: func(cmd *cobra.Command, args []string) error {
			profile, err := getCurrentProfile(cmd.Context())
			if err != nil {
				return err
			}

			name := profile.FirstName
			if profile.LastName != nil {
				name += " " + *profile.LastName
			}

			return util.PrintFields(cmd.OutOrStdout(),
				"User ID", profile.UserId,
				"Name", name,
				"Email", profile.Email,
				"Created At", util.FormatTime(profile.CreatedAt),
				"Updated At", util.FormatTime(profile.UpdatedAt),
			)
		},
	}
}

func getCurrentProfile(ctx context.Context) (*contract.PodletProfileResponse, error) {
	// Fetch current user
	api, token, err := util.NewAuthedClient()
	if err != nil {
		return nil, err
	}

	var profile contract.PodletProfileResponse
	if _, err := api.DoJSON(
		ctx,
		http.MethodGet,
		"/profile",
		nil,
		&profile,
		token,
	); err != nil {
		return nil, fmt.Errorf("get profile: %w", err)
	}

	return &profile, nil
}
