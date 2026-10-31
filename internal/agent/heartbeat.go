package agent

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SudiMango/podlet/internal/contract"
)

func (a *Agent) Heartbeat(ctx context.Context) (*contract.HeartbeatResponse, error) {
	var response contract.HeartbeatResponse
	if err := a.client.DoJSON(
		ctx,
		http.MethodPost,
		"/agent/heartbeat",
		nil,
		&response,
	); err != nil {
		return nil, err
	}

	return &response, nil
}

func (a *Agent) reportHeartbeat(ctx context.Context) (bool, error) {
	response, err := a.Heartbeat(ctx)
	if err != nil {
		return false, fmt.Errorf("heartbeat: %w", err)
	}
	if response.Action == nil {
		return false, nil
	}

	switch *response.Action {
	case contract.ActionUninstall:
		a.uninstallRequested.Store(true)

		if a.busy.Load() {
			return false, nil
		}

		if err := a.Uninstall(ctx); err != nil {
			return false, err
		}
		return true, nil
	default:
		return false, fmt.Errorf("unsupported heartbeat action %q", *response.Action)
	}
}
