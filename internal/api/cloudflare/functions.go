package cloudflare

import (
	"context"
	"net/http"

	"github.com/SudiMango/podlet/internal/contract"
)

const codeLength = 7

func (c *Client) CreateTunnel(ctx context.Context) (*contract.CreateTunnelResponse, error) {
	// Generate random code
	code, err := generateCode(codeLength)
	if err != nil {
		return nil, err
	}

	// Create tunnel
	request := CreateTunnelRequest{
		Name:      "podlet-" + code,
		ConfigSrc: "cloudflare",
	}

	var response CloudflareResponse[contract.CreateTunnelResponse]
	err = c.doJSON(
		ctx,
		http.MethodPost,
		"/accounts/"+c.accountId+"/cfd_tunnel",
		request,
		http.StatusOK,
		&response,
	)
	if err != nil {
		return nil, err
	}
	if err := validateResponse(response); err != nil {
		return nil, err
	}

	return &response.Result, nil
}

func (c *Client) GetTunnelConnectorToken(ctx context.Context, tunnelId string) (*string, error) {
	var response CloudflareResponse[string]
	err := c.doJSON(
		ctx,
		http.MethodGet,
		"/accounts/"+c.accountId+"/cfd_tunnel/"+tunnelId+"/token",
		nil,
		http.StatusOK,
		&response,
	)
	if err != nil {
		return nil, err
	}
	if err := validateResponse(response); err != nil {
		return nil, err
	}

	return &response.Result, nil
}

func (c *Client) GetIngressConfigs(ctx context.Context, tunnelId string) (*contract.GetIngressConfigsResponse, error) {
	var response CloudflareResponse[contract.GetIngressConfigsResponse]
	err := c.doJSON(
		ctx,
		http.MethodGet,
		"/accounts/"+c.accountId+"/cfd_tunnel/"+tunnelId+"/configurations",
		nil,
		http.StatusOK,
		&response,
	)
	if err != nil {
		return nil, err
	}
	if err := validateResponse(response); err != nil {
		return nil, err
	}

	ingress := make([]contract.Ingress, 0, len(response.Result.Config.Ingress))
	for _, rule := range response.Result.Config.Ingress {
		if rule.Service != "http_status:404" {
			ingress = append(ingress, rule)
		}
	}
	response.Result.Config.Ingress = ingress

	return &response.Result, nil
}

func (c *Client) UpdateIngressConfig(ctx context.Context, tunnelId string, request PutIngressRequest) (*contract.GetIngressConfigsResponse, error) {

	// Remove 404s
	ingress := make([]contract.Ingress, 0, len(request.Config.Ingress)+1)
	for _, rule := range request.Config.Ingress {
		if rule.Service != "http_status:404" {
			ingress = append(ingress, rule)
		}
	}

	// Add 1 404
	ingress = append(ingress, contract.Ingress{
		Service: "http_status:404",
	})
	request.Config.Ingress = ingress

	var response CloudflareResponse[contract.GetIngressConfigsResponse]
	err := c.doJSON(
		ctx,
		http.MethodPut,
		"/accounts/"+c.accountId+"/cfd_tunnel/"+tunnelId+"/configurations",
		request,
		http.StatusOK,
		&response,
	)
	if err != nil {
		return nil, err
	}
	if err := validateResponse(response); err != nil {
		return nil, err
	}

	return &response.Result, nil
}

func (c *Client) CreateDNSRecord(ctx context.Context, tunnelId string) (*contract.CreateDNSRecordResponse, error) {
	// Generate random code
	code, err := generateCode(codeLength)
	if err != nil {
		return nil, err
	}

	// Create record
	request := CreateDNSRecordRequest{
		Type:    "CNAME",
		Name:    code + ".podlet.dev",
		Content: tunnelId + ".cfargotunnel.com",
		Proxied: true,
		TTL:     1,
		Comment: "Managed by Podlet",
	}

	var response CloudflareResponse[contract.CreateDNSRecordResponse]
	err = c.doJSON(
		ctx,
		http.MethodPost,
		"/zones/"+c.zoneId+"/dns_records",
		request,
		http.StatusOK,
		&response,
	)
	if err != nil {
		return nil, err
	}
	if err := validateResponse(response); err != nil {
		return nil, err
	}

	return &response.Result, nil
}
