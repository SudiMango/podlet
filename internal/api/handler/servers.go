package handler

import (
	"errors"
	"net/http"

	"github.com/SudiMango/podlet/internal/api/auth"
	"github.com/SudiMango/podlet/internal/api/service"
	"github.com/SudiMango/podlet/internal/api/util"
	"github.com/SudiMango/podlet/internal/contract"
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

	contract.PendingServerRequest

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
