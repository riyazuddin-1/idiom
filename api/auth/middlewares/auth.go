package middlewares

import (
	"context"
	"errors"
	response "idiom-api-services/packages/responses"
	"net/http"
	"strings"

	"idiom-api-services/packages/jwt"
)

type contextKey string

const authUserContextKey contextKey = "auth_user"

type AuthUser struct {
	IdentityID string
	SessionID  string
}

func VerifyUserToken(jwtSettings *jwt.JWTSettings, w http.ResponseWriter, r *http.Request) (*http.Request, error) {
	authorization := r.Header.Get("Authorization")
	if !strings.HasPrefix(authorization, "Bearer ") {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return nil, errors.New("unauthorized")
	}

	token := authorization[len("Bearer "):]
	if token == "" {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return nil, errors.New("unauthorized")
	}

	claims, err := jwtSettings.VerifyToken(token)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return nil, errors.New("unauthorized")
	}

	identityID := claims.String("sub")
	sessionID := claims.String("sid")
	if identityID == "" || sessionID == "" {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return nil, errors.New("unauthorized")
	}

	user := &AuthUser{
		IdentityID: identityID,
		SessionID:  sessionID,
	}

	return r.WithContext(ContextWithUser(r.Context(), user)), nil
}

func ContextWithUser(ctx context.Context, user *AuthUser) context.Context {
	return context.WithValue(ctx, authUserContextKey, user)
}

func UserFromContext(ctx context.Context) (*AuthUser, bool) {
	user, ok := ctx.Value(authUserContextKey).(*AuthUser)
	return user, ok
}
