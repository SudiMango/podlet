package job

import "github.com/spf13/cobra"

type JobCommand struct{}

func (JobCommand) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jobs",
		Short: "Manage jobs",
	}

	cmd.AddCommand(
		DeployCommand{}.Command(),
		GetJobsCommand{}.Command(),
	)

	return cmd
}
