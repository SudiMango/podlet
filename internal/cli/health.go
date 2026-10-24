package cli

import (
	"fmt"
	"net/http"

	"github.com/SudiMango/podlet/internal/cli/config"
	"github.com/spf13/cobra"
)

type HealthCommand struct{}

func (HealthCommand) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "health",
		Short: "Check the health of podlet",
		Long:  "Check the health of podlet",
		RunE: func(cmd *cobra.Command, args []string) error {
			baseURL, err := config.LoadBaseURL()
			if err != nil {
				return err
			}

			response, err := http.Get(baseURL + "/health")
			if err != nil {
				return fmt.Errorf("check health: %w", err)
			}
			defer response.Body.Close()

			if response.StatusCode == http.StatusOK {
				cmd.Println("Podlet is healthy")
				return nil
			}

			return fmt.Errorf("Podlet is not healthy: %s", response.Status)
		},
	}
}
