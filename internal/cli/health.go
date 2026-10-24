package cli

import (
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

type HealthCommand struct{}

func (HealthCommand) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "health",
		Short: "Check the health of podlet",
		Long:  "Check the health of podlet",
		Run: func(cmd *cobra.Command, args []string) {
			response, err := http.Get(os.Getenv("BASE_URL") + "/health")
			if err != nil {
				cmd.PrintErrln("Error checking health: ", err)
				os.Exit(1)
			}
			defer response.Body.Close()

			if response.StatusCode == http.StatusOK {
				cmd.Println("Podlet is healthy")
			} else {
				cmd.Println("Podlet is not healthy. Status code: ", response.StatusCode)
				os.Exit(1)
			}
		},
	}
}
