package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"idiom-api-services/packages/jwt"
)

type contextKey string

const authUserContextKey contextKey = "auth_user"

type AuthUser struct {
	Email string
}

func VerifyUserToken(jwtSettings *jwt.JWTSettings, w http.ResponseWriter, r *http.Request) (*http.Request, error) {
	authorization := r.Header.Get("Authorization")
	if !strings.HasPrefix(authorization, "Bearer ") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return nil, errors.New("Unauthorized")
	}

	token := authorization[len("Bearer "):]
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return nil, errors.New("Unauthorized")
	}

	claims, err := jwtSettings.VerifyToken(token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return nil, errors.New("Unauthorized")
	}

	user := &AuthUser{
		Email: claims.Email,
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
