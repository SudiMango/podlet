package jobs

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/SudiMango/podlet/internal/contract"
)

type EnsureRailpackJob struct{}

func (EnsureRailpackJob) Execute(ctx context.Context, job *contract.ClaimJobResponse) error {
	output, err := exec.CommandContext(ctx, "railpack", "--version").CombinedOutput()
	if err != nil {
		return fmt.Errorf("check for Railpack: %w: %s", err, strings.TrimSpace(string(output)))
	}

	return nil
}
