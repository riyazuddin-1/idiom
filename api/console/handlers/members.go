package handlers

import (
	"encoding/json"
	"errors"
	authmiddlewares "idiom-api-services/api/auth/middlewares"
	"idiom-api-services/domains/org_members"
	response "idiom-api-services/packages/responses"
	"log"
	"net/http"
)

func (h *Handler) ListOrganizationMembersHandler(w http.ResponseWriter, r *http.Request) {
	r, err := authmiddlewares.VerifyUserToken(h.config.JWTSettings, w, r)
	if err != nil {
		log.Printf("console organization members list failed: unauthorized: %v", err)
		return
	}

	user, ok := authmiddlewares.UserFromContext(r.Context())
	if !ok {
		log.Printf("console organization members list failed: missing authenticated user")
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	organizationID := r.PathValue("oid")
	if organizationID == "" {
		response.Error(w, http.StatusBadRequest, "Organization ID is required")
		return
	}

	members, err := org_members.ListByOrganizationID(
		r.Context(),
		h.memberRepo,
		organizationID,
	)
	if err != nil {
		log.Printf(
			"console organization members list failed identity=%q organization=%q: %v",
			user.IdentityID,
			organizationID,
			err,
		)
		response.Error(w, http.StatusInternalServerError, "Failed to list organization members")
		return
	}

	log.Printf(
		"console organization members list succeeded identity=%q organization=%q count=%d",
		user.IdentityID,
		organizationID,
		len(members),
	)

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"members": members,
	})
}

func (h *Handler) UpdateOrganizationMemberStatusHandler(w http.ResponseWriter, r *http.Request) {
	r, err := authmiddlewares.VerifyUserToken(h.config.JWTSettings, w, r)
	if err != nil {
		log.Printf("console organization member status update failed: unauthorized: %v", err)
		return
	}

	user, ok := authmiddlewares.UserFromContext(r.Context())
	if !ok {
		log.Printf("console organization member status update failed: missing authenticated user")
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	organizationID := r.PathValue("oid")
	identityID := r.PathValue("uid")

	if organizationID == "" || identityID == "" {
		response.Error(w, http.StatusBadRequest, "Organization ID and identity ID are required")
		return
	}

	var req struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	member, err := org_members.UpdateStatus(
		r.Context(),
		h.memberRepo,
		organizationID,
		identityID,
		req.Status,
	)
	if err != nil {
		log.Printf(
			"console organization member status update failed identity=%q organization=%q member=%q: %v",
			user.IdentityID,
			organizationID,
			identityID,
			err,
		)

		switch {
		case errors.Is(err, org_members.ErrInvalidMemberStatus):
			response.Error(w, http.StatusBadRequest, "Invalid member status")
		case errors.Is(err, org_members.ErrMemberNotFound):
			response.Error(w, http.StatusNotFound, "Member not found")
		default:
			response.Error(w, http.StatusInternalServerError, "Failed to update member status")
		}

		return
	}

	log.Printf(
		"console organization member status update succeeded identity=%q organization=%q member=%q status=%q",
		user.IdentityID,
		organizationID,
		identityID,
		member.Status,
	)

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"member": member,
	})
}

func (h *Handler) DeleteOrganizationMemberHandler(w http.ResponseWriter, r *http.Request) {
	r, err := authmiddlewares.VerifyUserToken(h.config.JWTSettings, w, r)
	if err != nil {
		log.Printf("console organization member delete failed: unauthorized: %v", err)
		return
	}

	user, ok := authmiddlewares.UserFromContext(r.Context())
	if !ok {
		log.Printf("console organization member delete failed: missing authenticated user")
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	organizationID := r.PathValue("oid")
	identityID := r.PathValue("uid")

	if organizationID == "" || identityID == "" {
		response.Error(w, http.StatusBadRequest, "Organization ID and identity ID are required")
		return
	}

	if err := org_members.Delete(
		r.Context(),
		h.memberRepo,
		organizationID,
		identityID,
	); err != nil {
		log.Printf(
			"console organization member delete failed identity=%q organization=%q member=%q: %v",
			user.IdentityID,
			organizationID,
			identityID,
			err,
		)

		if errors.Is(err, org_members.ErrMemberNotFound) {
			response.Error(w, http.StatusNotFound, "Member not found")
			return
		}

		response.Error(w, http.StatusInternalServerError, "Failed to delete organization member")
		return
	}

	log.Printf(
		"console organization member delete succeeded identity=%q organization=%q member=%q",
		user.IdentityID,
		organizationID,
		identityID,
	)

	w.WriteHeader(http.StatusNoContent)
}
