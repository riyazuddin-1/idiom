package handlers

import (
	"encoding/json"
	"errors"
	authmiddlewares "idiom-api-services/api/auth/middlewares"
	apikeys "idiom-api-services/domains/api_keys"
	response "idiom-api-services/packages/responses"
	"log"
	"net/http"
)

func (h *Handler) ListProjectAPIKeysHandler(w http.ResponseWriter, r *http.Request) {
	r, err := authmiddlewares.VerifyUserToken(h.config.JWTSettings, w, r)
	if err != nil {
		log.Printf("console api keys list failed: unauthorized: %v", err)
		return
	}

	user, ok := authmiddlewares.UserFromContext(r.Context())
	if !ok {
		log.Printf("console api keys list failed: missing authenticated user")
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	projectID := r.PathValue("pid")
	if projectID == "" {
		response.Error(w, http.StatusBadRequest, "Project ID is required")
		return
	}

	apiKeys, err := apikeys.ListByProjectID(
		r.Context(),
		h.apiKeyRepo,
		projectID,
		user.IdentityID,
	)
	if err != nil {
		log.Printf(
			"console api keys list failed identity=%q project=%q: %v",
			user.IdentityID,
			projectID,
			err,
		)
		response.Error(w, http.StatusInternalServerError, "Failed to list API keys")
		return
	}

	log.Printf(
		"console api keys list succeeded identity=%q project=%q count=%d",
		user.IdentityID,
		projectID,
		len(apiKeys),
	)

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"api_keys": apiKeys,
	})
}

func (h *Handler) CreateProjectAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	r, err := authmiddlewares.VerifyUserToken(h.config.JWTSettings, w, r)
	if err != nil {
		log.Printf("console api key create failed: unauthorized: %v", err)
		return
	}

	user, ok := authmiddlewares.UserFromContext(r.Context())
	if !ok {
		log.Printf("console api key create failed: missing authenticated user")
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	projectID := r.PathValue("pid")
	if projectID == "" {
		response.Error(w, http.StatusBadRequest, "Project ID is required")
		return
	}

	var req apikeys.CreateInput

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := apikeys.Create(
		r.Context(),
		h.apiKeyRepo,
		projectID,
		user.IdentityID,
		req,
	)
	if err != nil {
		log.Printf(
			"console api key create failed identity=%q project=%q: %v",
			user.IdentityID,
			projectID,
			err,
		)

		if errors.Is(err, apikeys.ErrNameRequired) {
			response.Error(w, http.StatusBadRequest, "API key name is required")
			return
		}

		if errors.Is(err, apikeys.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "Project not found")
			return
		}

		response.Error(w, http.StatusInternalServerError, "Failed to create API key")
		return
	}

	log.Printf(
		"console api key create succeeded identity=%q project=%q key=%q",
		user.IdentityID,
		projectID,
		result.APIKey.ID,
	)

	response.JSON(w, http.StatusCreated, result)
}

func (h *Handler) RotateProjectAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	r, err := authmiddlewares.VerifyUserToken(h.config.JWTSettings, w, r)
	if err != nil {
		log.Printf("console api key rotation failed: unauthorized: %v", err)
		return
	}

	user, ok := authmiddlewares.UserFromContext(r.Context())
	if !ok {
		log.Printf("console api key rotation failed: missing authenticated user")
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	projectID := r.PathValue("pid")
	apiKeyID := r.PathValue("kid")

	if projectID == "" || apiKeyID == "" {
		response.Error(w, http.StatusBadRequest, "Project ID and API key ID are required")
		return
	}

	result, err := apikeys.Rotate(
		r.Context(),
		h.apiKeyRepo,
		projectID,
		user.IdentityID,
		apiKeyID,
	)
	if err != nil {
		log.Printf(
			"console api key rotation failed identity=%q project=%q key=%q: %v",
			user.IdentityID,
			projectID,
			apiKeyID,
			err,
		)

		if errors.Is(err, apikeys.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "API key not found")
			return
		}

		response.Error(w, http.StatusInternalServerError, "Failed to rotate API key")
		return
	}

	log.Printf(
		"console api key rotation succeeded identity=%q project=%q key=%q",
		user.IdentityID,
		projectID,
		apiKeyID,
	)

	response.JSON(w, http.StatusOK, result)
}

func (h *Handler) RevokeProjectAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	r, err := authmiddlewares.VerifyUserToken(h.config.JWTSettings, w, r)
	if err != nil {
		log.Printf("console api key revoke failed: unauthorized: %v", err)
		return
	}

	user, ok := authmiddlewares.UserFromContext(r.Context())
	if !ok {
		log.Printf("console api key revoke failed: missing authenticated user")
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	projectID := r.PathValue("pid")
	apiKeyID := r.PathValue("kid")

	if projectID == "" || apiKeyID == "" {
		response.Error(w, http.StatusBadRequest, "Project ID and API key ID are required")
		return
	}

	err = apikeys.Revoke(
		r.Context(),
		h.apiKeyRepo,
		projectID,
		user.IdentityID,
		apiKeyID,
	)
	if err != nil {
		log.Printf(
			"console api key revoke failed identity=%q project=%q key=%q: %v",
			user.IdentityID,
			projectID,
			apiKeyID,
			err,
		)

		if errors.Is(err, apikeys.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "API key not found")
			return
		}

		response.Error(w, http.StatusInternalServerError, "Failed to revoke API key")
		return
	}

	log.Printf(
		"console api key revoke succeeded identity=%q project=%q key=%q",
		user.IdentityID,
		projectID,
		apiKeyID,
	)

	w.WriteHeader(http.StatusNoContent)
}
