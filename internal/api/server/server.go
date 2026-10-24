package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"github.com/SudiMango/podlet/internal/api/auth"
	"github.com/SudiMango/podlet/internal/api/handler"
	"github.com/SudiMango/podlet/internal/api/repository"
	"github.com/go-chi/chi/v5"
	"github.com/thecodearcher/limen"
	"go.uber.org/fx"
)

var Module = fx.Module("server",
	fx.Provide(provideRouter),
	fx.Invoke(startServer),
)

type RouterParams struct {
	fx.In

	Auth             *limen.Limen
	ServerRepository *repository.ServerRepository
	HealthHandler    *handler.HealthHandler
	ProfileHandler   *handler.ProfileHandler
	ServerHandler    *handler.ServerHandler
}

func provideRouter(params RouterParams) *chi.Mux {
	r := chi.NewRouter()

	// All public routes
	r.Get("/health", params.HealthHandler.IsDatabaseHealthy)
	r.Handle("/auth", params.Auth.Handler())
	r.Handle("/auth/*", params.Auth.Handler())

	r.Post("/servers/register", params.ServerHandler.CreateServerByUserId)

	// All protected routes
	r.Group(func(r chi.Router) {
		r.Use(auth.RequireAuth(params.Auth))

		r.Get("/profile", params.ProfileHandler.GetPodletProfileByUserId)

		r.Route("/servers", func(r chi.Router) {
			r.Post("/add", params.ServerHandler.CreatePendingServerByUserId)
			r.Get("/pending", params.ServerHandler.GetValidPendingServersByUserId)
			r.Get("/pending/{pendingServerId}", params.ServerHandler.GetValidPendingServerByUserIdAndPendingServerId)
			r.Delete("/pending/{pendingServerId}", params.ServerHandler.DeletePendingServerByUserIdAndPendingServerId)

			r.Get("/", params.ServerHandler.GetServersByUserId)
			r.Get("/{serverId}", params.ServerHandler.GetServerByUserIdAndServerId)
			r.Patch("/{serverId}", params.ServerHandler.UpdateServerByUserIdAndServerId)
			r.Delete("/{serverId}", params.ServerHandler.MarkServerForDeletionByUserIdAndServerId)
		})
	})

	// All server agent routes
	r.Group(func(r chi.Router) {
		r.Use(auth.RequireAgentAuth(params.ServerRepository))

		r.Route("/agent", func(r chi.Router) {
			r.Post("/heartbeat", params.ServerHandler.UpdateServerLastSeenAtByUserIdAndServerId)
			r.Delete("/", params.ServerHandler.DeleteServerByUserIdAndServerId)
		})
	})

	return r
}

func startServer(lc fx.Lifecycle, r *chi.Mux) error {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		return errors.New("BASE_URL environment variable is required")
	}

	parsed, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("error parsing BASE_URL: %w", err)
	}

	srv := &http.Server{Addr: parsed.Host, Handler: r}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			slog.InfoContext(ctx, "server running", "base_url", baseURL)
			go func() {
				if err := srv.ListenAndServe(); err != http.ErrServerClosed {
					slog.Error("server stopped unexpectedly", "error", err)
					os.Exit(1)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			slog.InfoContext(ctx, "stopping server...")
			if err := srv.Shutdown(ctx); err != nil {
				return fmt.Errorf("error shutting down server: %w", err)
			}
			return nil
		},
	})
	return nil
}
