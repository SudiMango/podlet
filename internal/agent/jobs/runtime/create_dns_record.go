package runtime

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SudiMango/podlet/internal/agent/client"
	"github.com/SudiMango/podlet/internal/contract"
)

type CreateDNSRecordJob struct {
	client *client.Client
}

func NewCreateDNSRecordJob(client *client.Client) *CreateDNSRecordJob {
	return &CreateDNSRecordJob{client: client}
}

func (j *CreateDNSRecordJob) Execute(ctx context.Context, job *contract.ClaimJobResponse) error {
	var tunnel contract.TunnelResponse
	if err := j.client.DoJSON(ctx, http.MethodGet, "/agent/cftunnel", nil, &tunnel); err != nil {
		return fmt.Errorf("create DNS record: get Cloudflare tunnel: %w", err)
	}

	request := struct {
		TunnelID string `json:"tunnel_id"`
		Service  string `json:"service"`
	}{
		TunnelID: tunnel.Id,
		Service:  "http://" + job.Id + ":3000",
	}

	var record contract.DnsRecordResponse
	if err := j.client.DoJSON(ctx, http.MethodPost, "/agent/cfdns", request, &record); err != nil {
		return fmt.Errorf("create DNS record: %w", err)
	}

	fmt.Printf("Deployment is available at https://%s\n", record.Name)
	return nil
}
