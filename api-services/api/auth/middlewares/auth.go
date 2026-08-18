package middlewares

import (
	"context"
	"errors"
	apikeys "idiom-api-services/domains/api_keys"
	"idiom-api-services/domains/projects"
	response "idiom-api-services/packages/responses"
	"log"
	"net/http"
	"strings"

	"idiom-api-services/packages/jwt"
)

type contextKey string

const authUserContextKey contextKey = "auth_user"
const projectContextKey contextKey = "project"
const apiKeyContextKey contextKey = "api_key"

type AuthUser struct {
	IdentityID string
	SessionID  string
	ProjectID  string
}

func VerifyProject(projectRepo *projects.Repository, w http.ResponseWriter, r *http.Request) (*http.Request, error) {
	projectIdentifier := r.PathValue("project_id")
	if projectIdentifier == "" {
		log.Printf("project verification failed: missing project path param path=%q", r.URL.Path)
		http.NotFound(w, r)
		return nil, errors.New("project path param missing")
	}

	project, active, err := projects.ResolveActive(r.Context(), projectRepo, projectIdentifier)
	if err != nil {
		log.Printf("project verification errored project=%q path=%q: %v", projectIdentifier, r.URL.Path, err)
		response.Error(w, http.StatusInternalServerError, "Failed to resolve project")
		return nil, err
	}
	if !active {
		log.Printf("project verification failed project=%q path=%q", projectIdentifier, r.URL.Path)
		http.NotFound(w, r)
		return nil, errors.New("project not found")
	}

	return r.WithContext(ContextWithProject(r.Context(), project)), nil
}

func VerifyAPIKey(apiKeyRepo *apikeys.Repository, projectID string, w http.ResponseWriter, r *http.Request) (*http.Request, error) {
	apiKey := r.Header.Get("X-API-Key")
	if apiKey == "" {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return nil, errors.New("api key missing")
	}

	key, err := apikeys.Verify(r.Context(), apiKeyRepo, projectID, apiKey)
	if err != nil {
		if errors.Is(err, apikeys.ErrInvalidKey) {
			response.Error(w, http.StatusUnauthorized, "Unauthorized")
			return nil, errors.New("invalid api key")
		}

		log.Printf("api key verification failed project=%q: %v", projectID, err)
		response.Error(w, http.StatusInternalServerError, "Failed to verify API key")
		return nil, err
	}

	return r.WithContext(ContextWithAPIKey(r.Context(), key)), nil
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

	if routeProject, ok := ProjectFromContext(r.Context()); ok && routeProject.ID != projectID {
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

func ContextWithProject(ctx context.Context, project *projects.Project) context.Context {
	return context.WithValue(ctx, projectContextKey, project)
}

func ProjectFromContext(ctx context.Context) (*projects.Project, bool) {
	project, ok := ctx.Value(projectContextKey).(*projects.Project)
	return project, ok && project != nil
}

func ContextWithAPIKey(ctx context.Context, apiKey *apikeys.APIKey) context.Context {
	return context.WithValue(ctx, apiKeyContextKey, apiKey)
}

func APIKeyFromContext(ctx context.Context) (*apikeys.APIKey, bool) {
	apiKey, ok := ctx.Value(apiKeyContextKey).(*apikeys.APIKey)
	return apiKey, ok && apiKey != nil
}
