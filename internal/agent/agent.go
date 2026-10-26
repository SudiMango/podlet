package agent

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/SudiMango/podlet/internal/agent/config"
)

const (
	agentPollingInterval time.Duration = 30 * time.Second
)

type Agent struct {
	config *config.Config
	client *Client
	busy   bool
}

func New() (*Agent, error) {
	agentConfig, err := config.LoadAgentConfig()
	if err != nil {
		return nil, err
	}

	return &Agent{
		config: agentConfig,
		client: NewClient(agentConfig.APIURL, agentConfig.Token),
		busy:   false,
	}, nil
}

func (a *Agent) Run(ctx context.Context) error {
	// Report heartbeat
	uninstalled, err := a.reportHeartbeat(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	} else if uninstalled {
		return nil
	}

	// Attempt to claim job
	if !a.busy {
		claimed, err := a.claimJob(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		if claimed != "" {
			a.reportJobClaimed(ctx, claimed)
		}
	}

	ticker := time.NewTicker(agentPollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-ticker.C:
			// Report heartbeat
			uninstalled, err := a.reportHeartbeat(ctx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				continue
			}
			if uninstalled {
				return nil
			}

			// Attempt to claim job
			if !a.busy {
				claimed, err := a.claimJob(ctx)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%v\n", err)
				}
				if claimed != "" {
					a.reportJobClaimed(ctx, claimed)
				}
			}
		}
	}
}
