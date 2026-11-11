package runtime

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/SudiMango/podlet/internal/contract"
)

const buildkitHost = "docker-container://podlet-buildkit"

type BuildRepositoryJob struct{}

func imageRepo(repositoryId int64) string {
	return "podlet/repo-" + strconv.FormatInt(repositoryId, 10)
}

func (BuildRepositoryJob) Execute(ctx context.Context, job *contract.ClaimJobResponse) error {
	var payload contract.CloneRepositoryJob
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return fmt.Errorf("decode payload: %w", err)
	}

	repositoryPath := filepath.Join(repositoryRoot, job.Id)
	image := imageRepo(payload.RepositoryId) + ":" + job.Id
	fmt.Printf("Building repository %s as image %s...\n", repositoryPath, image)

	cmd := exec.CommandContext(
		ctx,
		"railpack", "build",
		"--name", image,
		"--progress", "plain",
		repositoryPath,
	)

	cmd.Env = append(os.Environ(), "BUILDKIT_HOST="+buildkitHost)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)

	if err := cmd.Run(); err != nil {
		output := strings.TrimSpace(stdout.String() + "\n" + stderr.String())
		return fmt.Errorf("build repository with Railpack: %w: %s",
			err, output)
	}

	fmt.Printf("Repository build completed successfully: %s\n", job.Id)
	fmt.Printf("Starting container %s...\n", job.Id)

	runOutput, err := exec.CommandContext(
		ctx,
		"docker", "run",
		"--detach",
		"--restart", "unless-stopped",
		"--publish", "127.0.0.1::3000",
		"--env", "PORT=3000",
		"--network", "podlet",
		"--name", job.Id,

		// Metadata
		"--label", "podlet.repository="+payload.RepoName,
		"--label", "podlet.repository-id="+strconv.FormatInt(payload.RepositoryId, 10),
		"--label", "podlet.deployment="+job.Id,

		image,
	).CombinedOutput()
	if err != nil {
		return fmt.Errorf("start repository container: %w: %s",
			err, strings.TrimSpace(string(runOutput)))
	}

	portOutput, err := exec.CommandContext(ctx, "docker", "port", job.Id).CombinedOutput()
	if err != nil {
		_, _ = exec.Command("docker", "rm", "--force", job.Id).CombinedOutput()
		return fmt.Errorf("get repository container port: %w: %s",
			err, strings.TrimSpace(string(portOutput)))
	}

	ports := strings.TrimSpace(string(portOutput))
	if ports == "" {
		_, _ = exec.Command("docker", "rm", "--force", job.Id).CombinedOutput()
		return fmt.Errorf("repository image %s does not expose a port", image)
	}

	fmt.Printf("Repository container started successfully:\n%s\n", ports)

	return nil
}
