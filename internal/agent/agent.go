package agent

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/SudiMango/podlet/internal/agent/client"
	"github.com/SudiMango/podlet/internal/agent/config"
	"github.com/SudiMango/podlet/internal/contract"
)

const (
	agentPollingInterval time.Duration = 30 * time.Second
)

type Agent struct {
	config    *config.Config
	client    *client.Client
	busy      atomic.Bool
	executors map[contract.JobType]JobExecutor
}

func New() (*Agent, error) {
	agentConfig, err := config.LoadAgentConfig()
	if err != nil {
		return nil, err
	}

	client := client.NewClient(agentConfig.APIURL, agentConfig.Token)

	return &Agent{
		config:    agentConfig,
		client:    client,
		executors: JobRegistry(client, agentConfig),
	}, nil
}

func (a *Agent) Run(ctx context.Context) error {
	if err := a.Initialize(ctx, false); err != nil {
		return fmt.Errorf("initialize agent: %w", err)
	}

	// Report heartbeat
	uninstalled, err := a.reportHeartbeat(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	} else if uninstalled {
		return nil
	}

	// Attempt to claim job
	if !a.busy.Load() {
		_, err := a.claimJob(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
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
			if !a.busy.Load() {
				_, err := a.claimJob(ctx)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%v\n", err)
				}
			}
		}
	}
}
