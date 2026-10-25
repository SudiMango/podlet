package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/SudiMango/podlet/internal/api/repository"
	"github.com/SudiMango/podlet/internal/api/util"
	"github.com/thecodearcher/limen"
)

type contextKey string

const (
	userSessionKey   contextKey = "user_session"
	agentIdentityKey contextKey = "agent_identity"
)

type AgentIdentity struct {
	ServerID string
	UserID   string
}

// User authentication
func RequireAuth(auth *limen.Limen) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := auth.GetSession(r)
			if err != nil || session == nil {
				util.SendError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
				return
			}

			ctx := context.WithValue(r.Context(), userSessionKey, session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Agent authentication
func RequireAgentAuth(servers *repository.ServerRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			token, found := strings.CutPrefix(header, "Bearer ")
			if !found || token == "" {
				util.SendError(w, http.StatusUnauthorized, errors.New("Invalid agent token"))
				return
			}

			server, err := servers.GetServerByTokenHash(r.Context(), util.HashToken(token))
			if err != nil {
				util.SendError(w, http.StatusUnauthorized, errors.New("Invalid agent token"))
				return
			}

			identity := AgentIdentity{
				ServerID: server.Id,
				UserID:   server.UserId,
			}

			ctx := context.WithValue(r.Context(), agentIdentityKey, identity)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Retrieve user session
func UserSessionFromContext(ctx context.Context) *limen.ValidatedSession {
	session, _ := ctx.Value(userSessionKey).(*limen.ValidatedSession)
	return session
}

// Retrieve agent identity
func AgentIdentityFromContext(ctx context.Context) (AgentIdentity, bool) {
	identity, ok := ctx.Value(agentIdentityKey).(AgentIdentity)
	return identity, ok
}

// Retrieve current authenticated user ID
func CurrentUserID(r *http.Request) (string, bool) {
	session := UserSessionFromContext(r.Context())
	if session == nil || session.User == nil {
		return "", false
	}

	userID, ok := session.User.ID.(string)
	if !ok || userID == "" {
		return "", false
	}

	return userID, true
}
