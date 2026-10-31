package handler

import (
	"errors"
	"net/http"

	"github.com/SudiMango/podlet/internal/api/auth"
	"github.com/SudiMango/podlet/internal/api/dto"
	"github.com/SudiMango/podlet/internal/api/service"
	"github.com/SudiMango/podlet/internal/api/util"
)

type CloudflareHandler struct {
	cloudflareService *service.CloudflareService
}

func NewCloudflareHandler(cloudflareService *service.CloudflareService) *CloudflareHandler {
	return &CloudflareHandler{
		cloudflareService: cloudflareService,
	}
}

/*
	Tunnels
*/

/*
Create a Cloudflare tunnel for the authenticated server agent

API URL: POST /agent/cftunnel

Args:

	none

Returns:

	the created Cloudflare tunnel (HTTP 200).

Raises:

	401: invalid agent identity or authentication token
	500: unable to create or persist the Cloudflare tunnel
*/
func (h *CloudflareHandler) CreateTunnel(w http.ResponseWriter, r *http.Request) {
	agentIdentity, ok := auth.AgentIdentityFromContext(r.Context())
	if !ok {
		util.SendError(w, http.StatusUnauthorized, errors.New("Invalid agent identity"))
		return
	}

	tunnel, err := h.cloudflareService.CreateTunnel(r.Context(), agentIdentity.UserID, agentIdentity.ServerID)
	if err != nil {
		util.SendError(w, http.StatusInternalServerError, err)
		return
	}

	util.SendJSON(w, http.StatusOK, tunnel)
}

/*
Get the Cloudflare tunnel assigned to the authenticated server agent

API URL: GET /agent/cftunnel

Args:

	none

Returns:

	the Cloudflare tunnel assigned to the server (HTTP 200).

Raises:

	401: invalid agent identity or authentication token
	404: no Cloudflare tunnel is assigned to the server
	500: unable to retrieve the Cloudflare tunnel
*/
func (h *CloudflareHandler) GetTunnelByUserIdAndServerId(w http.ResponseWriter, r *http.Request) {
	agentIdentity, ok := auth.AgentIdentityFromContext(r.Context())
	if !ok {
		util.SendError(w, http.StatusUnauthorized, errors.New("Invalid agent identity"))
		return
	}

	tunnel, err := h.cloudflareService.GetTunnelByUserIdAndServerId(r.Context(), agentIdentity.UserID, agentIdentity.ServerID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCfTunnelNotFound):
			util.SendError(w, http.StatusNotFound, err)
		default:
			util.SendError(w, http.StatusInternalServerError, err)
		}
		return
	}

	util.SendJSON(w, http.StatusOK, tunnel)
}

/*
	DNS Records
*/

/*
Create a Cloudflare DNS record for the authenticated server agent

API URL: POST /agent/cfdns

Args:

	tunnel_id: the Cloudflare tunnel the DNS record should target
	hostname: the public hostname for the ingress rule
	service: the internal service to which the hostname should route

Returns:

	the created Cloudflare DNS record (HTTP 200).

Raises:

	400: invalid request body
	401: invalid agent identity or authentication token
	500: unable to create or persist the DNS record or ingress configuration
*/
func (h *CloudflareHandler) CreateDNSRecord(w http.ResponseWriter, r *http.Request) {
	agentIdentity, ok := auth.AgentIdentityFromContext(r.Context())
	if !ok {
		util.SendError(w, http.StatusUnauthorized, errors.New("Invalid agent identity"))
		return
	}

	var request dto.CreateDNSRecordRequest
	if err := util.DecodeJSON(r, &request); err != nil {
		util.SendError(w, http.StatusBadRequest, errors.New("Invalid request body"))
		return
	}

	dnsRecord, err := h.cloudflareService.CreateDNSRecord(r.Context(), agentIdentity.UserID, agentIdentity.ServerID, request)
	if err != nil {
		util.SendError(w, http.StatusInternalServerError, err)
		return
	}

	util.SendJSON(w, http.StatusOK, dnsRecord)
}
