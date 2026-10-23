package auth

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/thecodearcher/limen"
	sqladapter "github.com/thecodearcher/limen/adapters/sql"
	"github.com/thecodearcher/limen/plugins/oauth"
	oauthgithub "github.com/thecodearcher/limen/plugins/oauth-github"
)

type uuidGenerator struct{}

func (uuidGenerator) GetColumnType() limen.ColumnType {
	return limen.ColumnTypeUUID
}

func (uuidGenerator) Generate(context.Context) (any, error) {
	return uuid.NewString(), nil
}

func NewLimen(pool *pgxpool.Pool) (*limen.Limen, error) {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		return nil, errors.New("BASE_URL environment variable is required")
	}
	if os.Getenv("GITHUB_CLIENT_ID") == "" {
		return nil, errors.New("GITHUB_CLIENT_ID environment variable is required")
	}
	if os.Getenv("GITHUB_CLIENT_SECRET") == "" {
		return nil, errors.New("GITHUB_CLIENT_SECRET environment variable is required")
	}

	db := stdlib.OpenDBFromPool(pool)
	githubProvider := oauthgithub.New(oauthgithub.WithScopes())

	return limen.New(&limen.Config{
		BaseURL:  baseURL,
		Database: sqladapter.NewPostgreSQL(db),
		Session: limen.NewDefaultSessionConfig(
			limen.WithBearerEnabled(),
		),
		Schema: limen.NewDefaultSchemaConfig(
			limen.WithSchemaIDGenerator(uuidGenerator{}),
		),
		Plugins: []limen.Plugin{
			oauth.New(
				oauth.WithProviders(githubProvider),
				oauth.WithMapProfileToUser(githubProfileFields),
			),
			newGitHubDevicePlugin(os.Getenv("GITHUB_CLIENT_ID"), githubProvider),
		},
	})
}

func githubProfileFields(profile *limen.OAuthAccountProfile) map[string]any {
	name := strings.Fields(profile.Name)
	if len(name) == 0 {
		if login, ok := profile.Raw["login"].(string); ok && login != "" {
			return map[string]any{"first_name": login}
		}
		return map[string]any{"first_name": "GitHub user"}
	}

	fields := map[string]any{"first_name": name[0]}
	if len(name) > 1 {
		fields["last_name"] = strings.Join(name[1:], " ")
	}
	return fields
}
