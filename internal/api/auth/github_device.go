package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/thecodearcher/limen"
	limenoauth "github.com/thecodearcher/limen/plugins/oauth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

const (
	githubDevicePluginName limen.PluginName = "github-device"
	maxDeviceFlowDuration                   = 20 * time.Minute
)

type githubDevicePlugin struct {
	clientID  string
	provider  limenoauth.Provider
	oauthAPI  limenoauth.API
	core      *limen.LimenCore
	responder *limen.Responder
	config    oauth2.Config
}

func newGitHubDevicePlugin(clientID string, provider limenoauth.Provider) *githubDevicePlugin {
	return &githubDevicePlugin{
		clientID: clientID,
		provider: provider,
		config: oauth2.Config{
			ClientID: clientID,
			Endpoint: endpoints.GitHub,
		},
	}
}

func (p *githubDevicePlugin) Name() limen.PluginName {
	return githubDevicePluginName
}

func (p *githubDevicePlugin) Initialize(core *limen.LimenCore) error {
	if strings.TrimSpace(p.clientID) == "" {
		return errors.New("github device flow: client ID is required")
	}

	registered, ok := core.GetPlugin(limen.PluginOAuth)
	if !ok {
		return errors.New("github device flow: OAuth plugin must be registered first")
	}

	oauthAPI, ok := registered.(limenoauth.API)
	if !ok {
		return errors.New("github device flow: registered OAuth plugin has an incompatible API")
	}

	p.core = core
	p.oauthAPI = oauthAPI
	return nil
}

func (p *githubDevicePlugin) PluginHTTPConfig() limen.PluginHTTPConfig {
	return limen.PluginHTTPConfig{
		BasePath: "/device/github",
		RateLimitRules: []*limen.RateLimitRule{
			limen.NewRateLimitRule("/start", 10, time.Minute),
			limen.NewRateLimitRule("/complete", 10, time.Minute),
		},
	}
}

func (p *githubDevicePlugin) RegisterRoutes(httpCore *limen.LimenHTTPCore, routes *limen.RouteBuilder) {
	p.responder = httpCore.Responder
	routes.POST("/start", "github-device-start", p.start)
	routes.POST("/complete", "github-device-complete", p.complete)
}

func (p *githubDevicePlugin) start(w http.ResponseWriter, r *http.Request) {
	device, err := p.config.DeviceAuth(r.Context())
	if err != nil {
		p.respondError(w, r, fmt.Errorf("start GitHub device authorization: %w", err))
		return
	}

	_ = p.responder.JSON(w, r, http.StatusOK, GithubDeviceStartResponse{
		DeviceCode:              device.DeviceCode,
		UserCode:                device.UserCode,
		VerificationURI:         device.VerificationURI,
		VerificationURIComplete: device.VerificationURIComplete,
		ExpiresAt:               device.Expiry,
		Interval:                device.Interval,
	})
}

// complete waits for the user to authorize the device, then creates the same
// Limen account and session that the browser-based OAuth callback would create.
func (p *githubDevicePlugin) complete(w http.ResponseWriter, r *http.Request) {
	var request githubDeviceCompleteRequest
	if err := decodeJSON(r, &request); err != nil {
		p.respondError(w, r, limen.NewLimenError("invalid request body", http.StatusBadRequest, nil))
		return
	}

	if err := validateDeviceCompleteRequest(request); err != nil {
		p.respondError(w, r, limen.NewLimenError(err.Error(), http.StatusBadRequest, nil))
		return
	}

	device := &oauth2.DeviceAuthResponse{
		DeviceCode: request.DeviceCode,
		Expiry:     request.ExpiresAt,
		Interval:   request.Interval,
	}
	token, err := p.config.DeviceAccessToken(r.Context(), device)
	if err != nil {
		p.respondError(w, r, limen.NewLimenError("GitHub device authorization failed", http.StatusUnauthorized, nil))
		return
	}

	profile, err := p.oauthAPI.GetUserInfoWithTokens(r.Context(), p.provider, oauthTokenResponse(token))
	if err != nil {
		p.respondError(w, r, err)
		return
	}

	authResult, err := p.oauthAPI.CreateOrLinkAccount(r.Context(), profile)
	if err != nil {
		p.respondError(w, r, err)
		return
	}

	session, err := p.core.CreateSession(r.Context(), r, w, authResult)
	if err != nil {
		p.respondError(w, r, err)
		return
	}

	if err := p.responder.SessionResponse(w, r, p.core, authResult, session); err != nil {
		p.respondError(w, r, err)
	}
}

func (p *githubDevicePlugin) respondError(w http.ResponseWriter, r *http.Request, err error) {
	_ = p.responder.Error(w, r, err)
}

func validateDeviceCompleteRequest(request githubDeviceCompleteRequest) error {
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

func decodeJSON(r *http.Request, target any) error {
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain one JSON object")
	}
	return nil
}
