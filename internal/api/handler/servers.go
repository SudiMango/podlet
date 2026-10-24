package handler

import (
	"errors"
	"net/http"

	"github.com/SudiMango/podlet/internal/api/auth"
	"github.com/SudiMango/podlet/internal/api/service"
	"github.com/SudiMango/podlet/internal/api/util"
	"github.com/SudiMango/podlet/internal/contract"
	"github.com/go-chi/chi/v5"
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

API URL: GET /servers/get

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

API URL: GET /servers/get/{serverId}

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
