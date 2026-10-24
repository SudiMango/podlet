package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/SudiMango/podlet/internal/api/util"
	"github.com/thecodearcher/limen"
)

type contextKey string
type requestMetadataKey struct{}
type RequestMetadata struct {
	UserID string
}

const sessionKey contextKey = "session"

func RequireAuth(auth *limen.Limen) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := auth.GetSession(r)
			if err != nil || session == nil {
				util.SendError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
				return
			}

			if session.User != nil {
				if userID, ok := session.User.ID.(string); ok {
					if metadata := RequestMetadataFromContext(r.Context()); metadata != nil {
						metadata.UserID = userID
					}
				}
			}

			ctx := context.WithValue(r.Context(), sessionKey, session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func SessionFromContext(ctx context.Context) *limen.ValidatedSession {
	session, _ := ctx.Value(sessionKey).(*limen.ValidatedSession)
	return session
}

func RequestMetadataFromContext(ctx context.Context) *RequestMetadata {
	metadata, _ := ctx.Value(requestMetadataKey{}).(*RequestMetadata)
	return metadata
}

func CurrentUserID(r *http.Request) (string, bool) {
	session := SessionFromContext(r.Context())
	if session == nil || session.User == nil {
		return "", false
	}
	value, ok := session.User.ID.(string)
	if !ok || value == "" {
		return "", false
	}
	return value, true
}
