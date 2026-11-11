package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/SudiMango/podlet/internal/contract"
)

type RemovePrevDeploy struct{}

func (RemovePrevDeploy) Execute(ctx context.Context, job *contract.ClaimJobResponse) error {
	var payload contract.CloneRepositoryJob
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return fmt.Errorf("decode payload: %w", err)
	}

	// Previous containers
	output, err := exec.CommandContext(
		ctx,
		"docker", "ps", "-a",
		"--filter", "label=podlet.repository-id="+strconv.FormatInt(payload.RepositoryId, 10),
		"--format", "{{.Names}}",
	).Output()
	if err != nil {
		return fmt.Errorf("list previous containers: %w", err)
	}

	if names := without(strings.Fields(string(output)), job.Id); len(names) > 0 {
		args := append([]string{"rm", "-f"}, names...)
		if output, err := exec.CommandContext(ctx, "docker", args...).CombinedOutput(); err != nil {
			return fmt.Errorf("remove previous containers: %w: %s", err, strings.TrimSpace(string(output)))
		}
	}

	// Previous images
	currentImage := imageRepo(payload.RepositoryId) + ":" + job.Id
	output, err = exec.CommandContext(
		ctx,
		"docker", "images",
		"--filter", "reference="+imageRepo(payload.RepositoryId)+":*",
		"--format", "{{.Repository}}:{{.Tag}}",
	).Output()
	if err != nil {
		return fmt.Errorf("list previous images: %w", err)
	}

	if images := without(strings.Fields(string(output)), currentImage); len(images) > 0 {
		args := append([]string{"rmi", "-f"}, images...)
		if output, err := exec.CommandContext(ctx, "docker", args...).CombinedOutput(); err != nil {
			return fmt.Errorf("remove previous images: %w: %s", err, strings.TrimSpace(string(output)))
		}
	}

	return nil
}

func without(items []string, skip string) []string {
	kept := items[:0]
	for _, item := range items {
		if item != skip {
			kept = append(kept, item)
		}
	}
	return kept
}
