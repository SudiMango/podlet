package agent

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/SudiMango/podlet/internal/agent/client"
	"github.com/SudiMango/podlet/internal/agent/config"
)

type Agent struct {
	config *config.Config
	client *client.Client
}

func New() (*Agent, error) {
	agentConfig, err := config.LoadAgentConfig()
	if err != nil {
		return nil, err
	}

	return &Agent{
		config: agentConfig,
		client: client.New(agentConfig.APIURL, agentConfig.Token),
	}, nil
}

func (a *Agent) Heartbeat(ctx context.Context) error {
	return a.client.DoJSON(
		ctx,
		http.MethodPost,
		"/agent/heartbeat",
		nil,
		nil,
	)
}

func (a *Agent) reportHeartbeat(ctx context.Context) {
	if err := a.Heartbeat(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "heartbeat failed: %v\n", err)
	}
}

func (a *Agent) Run(ctx context.Context) error {
	a.reportHeartbeat(ctx)

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-ticker.C:
			a.reportHeartbeat(ctx)
		}
	}
}
