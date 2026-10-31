package initialize

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/SudiMango/podlet/internal/contract"
)

type InstallRailpackJob struct{}

func (InstallRailpackJob) Execute(ctx context.Context, job *contract.ClaimJobResponse) error {
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

	return nil
}
