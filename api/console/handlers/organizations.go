package handlers

import (
	"encoding/json"
	"errors"
	authmiddlewares "idiom-api-services/api/auth/middlewares"
	"idiom-api-services/domains/organizations"
	response "idiom-api-services/packages/responses"
	"log"
	"net/http"
)

func (h *Handler) CreateOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	r, err := authmiddlewares.VerifyUserToken(h.config.JWTSettings, w, r)
	if err != nil {
		log.Printf("console organization create failed: unauthorized: %v", err)
		return
	}

	user, ok := authmiddlewares.UserFromContext(r.Context())
	if !ok {
		log.Printf("console organization create failed: missing authenticated user")
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req organizations.CreateInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	organization, err := organizations.Create(
		r.Context(),
		h.orgRepo,
		user.IdentityID,
		req,
	)
	if err != nil {
		log.Printf(
			"console organization create failed identity=%q: %v",
			user.IdentityID,
			err,
		)

		if errors.Is(err, errors.New("organization name is required")) {
			response.Error(w, http.StatusBadRequest, "Organization name is required")
			return
		}

		response.Error(w, http.StatusInternalServerError, "Failed to create organization")
		return
	}

	log.Printf(
		"console organization create succeeded identity=%q organization=%q",
		user.IdentityID,
		organization.ID,
	)

	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"organization": organization,
	})
}

func (h *Handler) UpdateOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	r, err := authmiddlewares.VerifyUserToken(h.config.JWTSettings, w, r)
	if err != nil {
		log.Printf("console organization update failed: unauthorized: %v", err)
		return
	}

	user, ok := authmiddlewares.UserFromContext(r.Context())
	if !ok {
		log.Printf("console organization update failed: missing authenticated user")
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	organizationID := r.PathValue("oid")
	if organizationID == "" {
		response.Error(w, http.StatusBadRequest, "Organization ID is required")
		return
	}

	var req organizations.UpdateInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	organization, err := organizations.Update(
		r.Context(),
		h.orgRepo,
		user.IdentityID,
		organizationID,
		req,
	)
	if err != nil {
		log.Printf(
			"console organization update failed identity=%q organization=%q: %v",
			user.IdentityID,
			organizationID,
			err,
		)

		response.Error(w, http.StatusInternalServerError, "Failed to update organization")
		return
	}

	log.Printf(
		"console organization update succeeded identity=%q organization=%q",
		user.IdentityID,
		organizationID,
	)

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"organization": organization,
	})
}

func (h *Handler) UpdateOrganizationStatusHandler(w http.ResponseWriter, r *http.Request) {
	r, err := authmiddlewares.VerifyUserToken(h.config.JWTSettings, w, r)
	if err != nil {
		log.Printf("console organization status update failed: unauthorized: %v", err)
		return
	}

	user, ok := authmiddlewares.UserFromContext(r.Context())
	if !ok {
		log.Printf("console organization status update failed: missing authenticated user")
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	organizationID := r.PathValue("oid")
	if organizationID == "" {
		response.Error(w, http.StatusBadRequest, "Organization ID is required")
		return
	}

	var req struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	organization, err := organizations.UpdateStatus(
		r.Context(),
		h.orgRepo,
		user.IdentityID,
		organizationID,
		req.Status,
	)
	if err != nil {
		log.Printf(
			"console organization status update failed identity=%q organization=%q: %v",
			user.IdentityID,
			organizationID,
			err,
		)

		response.Error(w, http.StatusBadRequest, "Invalid organization status")
		return
	}

	log.Printf(
		"console organization status update succeeded identity=%q organization=%q status=%q",
		user.IdentityID,
		organizationID,
		organization.Status,
	)

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"organization": organization,
	})
}

func (h *Handler) DeleteOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	r, err := authmiddlewares.VerifyUserToken(h.config.JWTSettings, w, r)
	if err != nil {
		log.Printf("console organization delete failed: unauthorized: %v", err)
		return
	}

	user, ok := authmiddlewares.UserFromContext(r.Context())
	if !ok {
		log.Printf("console organization delete failed: missing authenticated user")
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	organizationID := r.PathValue("oid")
	if organizationID == "" {
		response.Error(w, http.StatusBadRequest, "Organization ID is required")
		return
	}

	if err := organizations.Delete(
		r.Context(),
		h.orgRepo,
		user.IdentityID,
		organizationID,
	); err != nil {
		log.Printf(
			"console organization delete failed identity=%q organization=%q: %v",
			user.IdentityID,
			organizationID,
			err,
		)

		response.Error(w, http.StatusNotFound, "Organization not found")
		return
	}

	log.Printf(
		"console organization delete succeeded identity=%q organization=%q",
		user.IdentityID,
		organizationID,
	)

	w.WriteHeader(http.StatusNoContent)
}
