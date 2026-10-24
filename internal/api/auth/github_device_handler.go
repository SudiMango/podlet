package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/SudiMango/podlet/internal/api/util"
	limenoauth "github.com/thecodearcher/limen/plugins/oauth"
	"golang.org/x/oauth2"
)

/*
	flow functions
*/

func (p *githubDevicePlugin) start(w http.ResponseWriter, r *http.Request) {
	device, err := p.config.DeviceAuth(r.Context())
	if err != nil {
		util.SendError(w, http.StatusInternalServerError, fmt.Errorf("start GitHub device authorization: %w", err))
		return
	}

	util.SendJSON(w, http.StatusOK, GithubDeviceStartResponse{
		DeviceCode:              device.DeviceCode,
		UserCode:                device.UserCode,
		VerificationURI:         device.VerificationURI,
		VerificationURIComplete: device.VerificationURIComplete,
		ExpiresAt:               device.Expiry,
		Interval:                device.Interval,
	})
}

func (p *githubDevicePlugin) complete(w http.ResponseWriter, r *http.Request) {
	var request GithubDeviceCompleteRequest
	if err := util.DecodeJSON(r, &request); err != nil {
		util.SendError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}

	if err := validateDeviceCompleteRequest(request); err != nil {
		util.SendError(w, http.StatusBadRequest, err)
		return
	}

	device := &oauth2.DeviceAuthResponse{
		DeviceCode: request.DeviceCode,
		Expiry:     request.ExpiresAt,
		Interval:   request.Interval,
	}
	token, err := p.config.DeviceAccessToken(r.Context(), device)
	if err != nil {
		util.SendError(w, http.StatusUnauthorized, errors.New("GitHub device authorization failed"))
		return
	}

	profile, err := p.oauthAPI.GetUserInfoWithTokens(r.Context(), p.provider, oauthTokenResponse(token))
	if err != nil {
		util.SendError(w, http.StatusInternalServerError, err)
		return
	}

	authResult, err := p.oauthAPI.CreateOrLinkAccount(r.Context(), profile)
	if err != nil {
		util.SendError(w, http.StatusInternalServerError, err)
		return
	}

	session, err := p.core.CreateSession(r.Context(), r, w, authResult)
	if err != nil {
		util.SendError(w, http.StatusInternalServerError, err)
		return
	}

	if err := p.responder.SessionResponse(w, r, p.core, authResult, session); err != nil {
		util.SendError(w, http.StatusInternalServerError, err)
	}
}

/*
	private helpers
*/

func validateDeviceCompleteRequest(request GithubDeviceCompleteRequest) error {
	now := time.Now()
	if strings.TrimSpace(request.DeviceCode) == "" {
		return errors.New("device_code is required")
	}
	if request.ExpiresAt.IsZero() || !request.ExpiresAt.After(now) {
		return errors.New("device authorization has expired")
	}
	if request.ExpiresAt.After(now.Add(maxDeviceFlowDuration)) {
		return errors.New("invalid device authorization expiry")
	}
	if request.Interval < 1 || request.Interval > 60 {
		return errors.New("invalid device authorization interval")
	}
	return nil
}

func oauthTokenResponse(token *oauth2.Token) *limenoauth.TokenResponse {
	scope, _ := token.Extra("scope").(string)
	idToken, _ := token.Extra("id_token").(string)
	return &limenoauth.TokenResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresAt:    token.Expiry,
		TokenType:    token.TokenType,
		Scope:        scope,
		IDToken:      idToken,
	}
}
