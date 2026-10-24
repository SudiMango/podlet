package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/thecodearcher/limen"
	limenoauth "github.com/thecodearcher/limen/plugins/oauth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

/*
	Constants
*/

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

/*
	github device plugin struct
*/

func newGithubDevicePlugin(clientID string, provider limenoauth.Provider) *githubDevicePlugin {
	return &githubDevicePlugin{
		clientID: clientID,
		provider: provider,
		config: oauth2.Config{
			ClientID: clientID,
			Endpoint: endpoints.GitHub,
		},
	}
}

/*
	Required plugin functions
*/

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
