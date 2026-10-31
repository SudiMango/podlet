package jobs

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/SudiMango/podlet/internal/contract"
)

type EnsurePodletNetworkJob struct{}

func (EnsurePodletNetworkJob) Execute(ctx context.Context, job *contract.ClaimJobResponse) error {
	fmt.Println("Checking for Podlet network...")

	output, err := exec.CommandContext(ctx, "docker", "network", "inspect", "podlet").CombinedOutput()
	if err != nil {
		return fmt.Errorf("check for Podlet network: %w: %s", err, strings.TrimSpace(string(output)))
	}

	fmt.Println("Podlet network found!")
	return nil
}
