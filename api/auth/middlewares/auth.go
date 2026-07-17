package middlewares

import (
	"context"
	"errors"
	"idiom-api-services/domains/projects"
	response "idiom-api-services/packages/responses"
	"net/http"
	"strings"

	"idiom-api-services/packages/jwt"
)

type contextKey string

const authUserContextKey contextKey = "auth_user"
const projectIDContextKey contextKey = "project_id"

type AuthUser struct {
	IdentityID string
	SessionID  string
	ProjectID  string
}

func VerifyProject(projectRepo *projects.Repository, w http.ResponseWriter, r *http.Request) (*http.Request, error) {
	projectID := r.PathValue("project_id")
	active, err := projects.IsActive(r.Context(), projectRepo, projectID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to resolve project")
		return nil, err
	}
	if !active {
		http.NotFound(w, r)
		return nil, errors.New("project not found")
	}

	return r.WithContext(ContextWithProjectID(r.Context(), projectID)), nil
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
	projectID := claims.String("pid")
	if identityID == "" || sessionID == "" || projectID == "" {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return nil, errors.New("unauthorized")
	}

	if routeProjectID, ok := ProjectIDFromContext(r.Context()); ok && routeProjectID != projectID {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return nil, errors.New("unauthorized")
	}

	user := &AuthUser{
		IdentityID: identityID,
		SessionID:  sessionID,
		ProjectID:  projectID,
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

func ContextWithProjectID(ctx context.Context, projectID string) context.Context {
	return context.WithValue(ctx, projectIDContextKey, projectID)
}

func ProjectIDFromContext(ctx context.Context) (string, bool) {
	projectID, ok := ctx.Value(projectIDContextKey).(string)
	return projectID, ok && projectID != ""
}
