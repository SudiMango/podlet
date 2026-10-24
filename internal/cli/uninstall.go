package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

type UninstallCommand struct{}

func (UninstallCommand) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall Podlet",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return uninstall(cmd)
		},
	}
}

func uninstall(cmd *cobra.Command) error {
	dataDir := os.Getenv("XDG_DATA_HOME")
	if dataDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("find home directory: %w", err)
		}
		dataDir = filepath.Join(homeDir, ".local", "share")
	}

	uninstallerPath := filepath.Join(dataDir, "podlet", "uninstall-cli.sh")
	if _, err := os.Stat(uninstallerPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("uninstaller not found at %s", uninstallerPath)
		}
		return fmt.Errorf("access uninstaller: %w", err)
	}

	uninstaller := exec.CommandContext(cmd.Context(), uninstallerPath)
	uninstaller.Stdin = cmd.InOrStdin()
	uninstaller.Stdout = cmd.OutOrStdout()
	uninstaller.Stderr = cmd.ErrOrStderr()

	if err := uninstaller.Run(); err != nil {
		return fmt.Errorf("uninstall Podlet: %w", err)
	}

	return nil
}
