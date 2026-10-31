package jobs

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/SudiMango/podlet/internal/contract"
)

type InstallRailpackJob struct{}

func (InstallRailpackJob) Execute(ctx context.Context, job *contract.ClaimJobResponse) error {
	fmt.Println("Installing Railpack...")

	output, err := exec.CommandContext(
		ctx,
		"bash",
		"-o",
		"pipefail",
		"-c",
		"curl -fsSL https://railpack.com/install.sh | sh -s -- --yes",
	).CombinedOutput()
	if err != nil {
		return fmt.Errorf("install Railpack: %w: %s", err, strings.TrimSpace(string(output)))
	}

	fmt.Println("Successfully installed Railpack!")
	return nil
}
