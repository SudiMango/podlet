package handler

import (
	"errors"
	"net/http"

	"github.com/SudiMango/podlet/internal/api/auth"
	"github.com/SudiMango/podlet/internal/api/service"
	"github.com/SudiMango/podlet/internal/api/util"
	"github.com/SudiMango/podlet/internal/contract"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type ServerHandler struct {
	serverService *service.ServerService
}

func NewServerHandler(serverService *service.ServerService) *ServerHandler {
	return &ServerHandler{
		serverService: serverService,
	}
}

/*
	Pending server functions
*/

/*
Create a pending server connection for a user

API URL: POST /servers/add

Args:

	none

Returns:

	the created pending server (HTTP 200).

Raises:

	401: unauthorized user
	500: unable to create pending server
*/
func (h *ServerHandler) CreatePendingServerByUserId(w http.ResponseWriter, r *http.Request) {

	userId, ok := auth.CurrentUserID(r)
	if !ok {
		util.SendError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	pendingServer, err := h.serverService.CreatePendingServerByUserId(r.Context(), userId)
	if err != nil {
		util.SendError(w, http.StatusInternalServerError, err)
		return
	}

	util.SendJSON(w, http.StatusOK, *pendingServer)
}

/*
Get all valid pending servers owned by a user

API URL: GET /servers/pending

Args:

	none

Returns:

	the user's valid pending servers (HTTP 200).

Raises:

	401: unauthorized user
	500: unable to get pending servers
*/
func (h *ServerHandler) GetValidPendingServersByUserId(w http.ResponseWriter, r *http.Request) {
	userId, ok := auth.CurrentUserID(r)
	if !ok {
		util.SendError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	pendingServers, err := h.serverService.GetValidPendingServersByUserId(r.Context(), userId)
	if err != nil {
		util.SendError(w, http.StatusInternalServerError, err)
		return
	}

	util.SendJSON(w, http.StatusOK, pendingServers)
}

/*
Get a valid pending server owned by a user

API URL: GET /servers/pending/{pendingServerId}

Args:

	pendingServerId: the pending server ID from the URL path

Returns:

	the requested pending server (HTTP 200).

Raises:

	400: pendingServerId is missing
	401: unauthorized user
	500: unable to get the pending server
*/
func (h *ServerHandler) GetValidPendingServerByUserIdAndPendingServerId(w http.ResponseWriter, r *http.Request) {
	userId, ok := auth.CurrentUserID(r)
	if !ok {
		util.SendError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	pendingServerId := chi.URLParam(r, "pendingServerId")
	if pendingServerId == "" {
		util.SendError(w, http.StatusBadRequest, errors.New("pendingServerId is required"))
		return
	}

	pendingServer, err := h.serverService.GetValidPendingServerByUserIdAndPendingServerId(r.Context(), userId, pendingServerId)
	if err != nil {
		util.SendError(w, http.StatusInternalServerError, err)
		return
	}

	util.SendJSON(w, http.StatusOK, pendingServer)
}

/*
Delete a pending server owned by a user

API URL: DELETE /servers/pending/{pendingServerId}

Args:

	pendingServerId: the pending server ID from the URL path

Returns:

	an empty response (HTTP 204).

Raises:

	400: pendingServerId is missing
	401: unauthorized user
	404: pending server not found
	500: unable to delete the pending server
*/
func (h *ServerHandler) DeletePendingServerByUserIdAndPendingServerId(w http.ResponseWriter, r *http.Request) {
	userId, ok := auth.CurrentUserID(r)
	if !ok {
		util.SendError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	pendingServerId := chi.URLParam(r, "pendingServerId")
	if pendingServerId == "" {
		util.SendError(w, http.StatusBadRequest, errors.New("pendingServerId is required"))
		return
	}

	if err := h.serverService.DeletePendingServerByUserIdAndPendingServerId(r.Context(), userId, pendingServerId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			util.SendError(w, http.StatusNotFound, errors.New("pending server not found"))
			return
		}

		util.SendError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/*
	Server functions
*/

/*
Register a server using a pending registration token

API URL: POST /servers/register

Args:

	contract.CreateServerRequest

Returns:

	the registered server and its authentication token (HTTP 200).

Raises:

	400: invalid request body
	500: unable to register server
*/
func (h *ServerHandler) CreateServerByUserId(w http.ResponseWriter, r *http.Request) {
	var request contract.CreateServerRequest
	if err := util.DecodeJSON(r, &request); err != nil {
		util.SendError(w, http.StatusBadRequest, errors.New("Invalid request body"))
		return
	}

	server, err := h.serverService.CreateServerByUserId(r.Context(), request)
	if err != nil {
		util.SendError(w, http.StatusInternalServerError, err)
		return
	}

	util.SendJSON(w, http.StatusOK, *server)
}

/*
Get all registered servers owned by a user

API URL: GET /servers

Args:

	none

Returns:

	the user's registered servers (HTTP 200).

Raises:

	401: unauthorized user
	500: unable to get registered servers
*/
func (h *ServerHandler) GetServersByUserId(w http.ResponseWriter, r *http.Request) {
	userId, ok := auth.CurrentUserID(r)
	if !ok {
		util.SendError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	servers, err := h.serverService.GetServersByUserId(r.Context(), userId)
	if err != nil {
		util.SendError(w, http.StatusInternalServerError, err)
		return
	}

	util.SendJSON(w, http.StatusOK, servers)
}

/*
Get a registered server owned by a user

API URL: GET /servers/{serverId}

Args:

	serverId: the registered server ID from the URL path

Returns:

	the requested registered server (HTTP 200).

Raises:

	400: serverId is missing
	401: unauthorized user
	500: unable to get the registered server
*/
func (h *ServerHandler) GetServerByUserIdAndServerId(w http.ResponseWriter, r *http.Request) {
	userId, ok := auth.CurrentUserID(r)
	if !ok {
		util.SendError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	serverId := chi.URLParam(r, "serverId")
	if serverId == "" {
		util.SendError(w, http.StatusBadRequest, errors.New("serverId is required"))
		return
	}

	server, err := h.serverService.GetServerByUserIdAndServerId(r.Context(), userId, serverId)
	if err != nil {
		util.SendError(w, http.StatusInternalServerError, err)
		return
	}

	util.SendJSON(w, http.StatusOK, server)
}

/*
Update the authenticated server agent's last-seen timestamp

API URL: POST /agent/heartbeat

Args:

	none

Returns:

	a heartbeat action (HTTP 200) when the agent must perform an operation, or
	an empty response (HTTP 204) when no action is required.

Raises:

	401: invalid agent identity or authentication token
	500: unable to update the server's last-seen timestamp
*/
func (h *ServerHandler) UpdateServerLastSeenAtByUserIdAndServerId(w http.ResponseWriter, r *http.Request) {
	identity, ok := auth.AgentIdentityFromContext(r.Context())
	if !ok {
		util.SendError(w, http.StatusUnauthorized, errors.New("Invalid agent identity"))
		return
	}

	lifecycleState, err := h.serverService.UpdateServerLastSeenAtByUserIdAndServerId(r.Context(), identity.UserID, identity.ServerID)
	if err != nil {
		util.SendError(w, http.StatusInternalServerError, err)
		return
	}
	if lifecycleState != nil && *lifecycleState == contract.ServerLifecyclePendingDeletion {
		util.SendJSON(w, http.StatusOK, contract.HeartbeatResponse{
			Action: new(contract.ActionUninstall),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/*
Update a registered server owned by a user

API URL: PATCH /servers/{serverId}

Args:

	serverId: the registered server ID from the URL path
	contract.UpdateServerRequest

Returns:

	the updated registered server (HTTP 200).

Raises:

	400: serverId is missing or the request body is invalid
	401: unauthorized user
	500: unable to update the registered server
*/
func (h *ServerHandler) UpdateServerByUserIdAndServerId(w http.ResponseWriter, r *http.Request) {
	userId, ok := auth.CurrentUserID(r)
	if !ok {
		util.SendError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	serverId := chi.URLParam(r, "serverId")
	if serverId == "" {
		util.SendError(w, http.StatusBadRequest, errors.New("serverId is required"))
		return
	}

	var request contract.UpdateServerRequest
	if err := util.DecodeJSON(r, &request); err != nil {
		util.SendError(w, http.StatusBadRequest, errors.New("Invalid request body"))
		return
	}

	server, err := h.serverService.UpdateServerByUserIdAndServerId(r.Context(), userId, serverId, request)
	if err != nil {
		util.SendError(w, http.StatusInternalServerError, err)
		return
	}

	util.SendJSON(w, http.StatusOK, server)
}

/*
Delete the server represented by the authenticated server agent

API URL: DELETE /agent

Args:

	none

Returns:

	an empty response (HTTP 204).

Raises:

	401: invalid agent identity or authentication token
	404: registered server not found
	500: unable to delete the registered server
*/
func (h *ServerHandler) DeleteServerByUserIdAndServerId(w http.ResponseWriter, r *http.Request) {
	agentIdentity, ok := auth.AgentIdentityFromContext(r.Context())
	if !ok {
		util.SendError(w, http.StatusUnauthorized, errors.New("Invalid agent identity"))
		return
	}

	if err := h.serverService.DeleteServerByUserIdAndServerId(r.Context(), agentIdentity.UserID, agentIdentity.ServerID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			util.SendError(w, http.StatusNotFound, errors.New("registered server not found"))
			return
		}

		util.SendError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/*
Mark a registered server for deletion by its agent

API URL: DELETE /servers/{serverId}

Args:

	serverId: the registered server ID from the URL path

Returns:

	an empty response (HTTP 204). The server remains stored until its agent
	receives the uninstall action and deletes the registration.

Raises:

	400: serverId is missing
	401: unauthorized user
	404: registered server not found
	500: unable to mark the registered server for deletion
*/
func (h *ServerHandler) MarkServerForDeletionByUserIdAndServerId(w http.ResponseWriter, r *http.Request) {
	userId, ok := auth.CurrentUserID(r)
	if !ok {
		util.SendError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	serverId := chi.URLParam(r, "serverId")
	if serverId == "" {
		util.SendError(w, http.StatusBadRequest, errors.New("serverId is required"))
		return
	}

	if err := h.serverService.MarkServerForDeletionByUserIdAndServerId(r.Context(), userId, serverId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			util.SendError(w, http.StatusNotFound, errors.New("registered server not found"))
			return
		}

		util.SendError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
