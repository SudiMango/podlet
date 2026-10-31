package initialize

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/SudiMango/podlet/internal/contract"
)

type CreatePodletNetworkJob struct{}

func (CreatePodletNetworkJob) Execute(ctx context.Context, job *contract.ClaimJobResponse) error {
	output, err := exec.CommandContext(ctx, "docker", "network", "create", "podlet").CombinedOutput()
	if err != nil {
		return fmt.Errorf("create Podlet network: %w: %s", err, strings.TrimSpace(string(output)))
	}

	return nil
}
