package handlers

import (
	"encoding/json"
	"errors"
	authmiddlewares "idiom-api-services/api/auth/middlewares"
	"idiom-api-services/domains/projects"
	response "idiom-api-services/packages/responses"
	"log"
	"net/http"
)

func (h *Handler) ListOrganizationProjectsHandler(w http.ResponseWriter, r *http.Request) {
	r, err := authmiddlewares.VerifyUserToken(h.config.JWTSettings, w, r)
	if err != nil {
		log.Printf("console organization projects list failed: unauthorized: %v", err)
		return
	}

	user, ok := authmiddlewares.UserFromContext(r.Context())
	if !ok {
		log.Printf("console organization projects list failed: missing authenticated user")
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	organizationID := r.PathValue("oid")
	if organizationID == "" {
		response.Error(w, http.StatusBadRequest, "Organization ID is required")
		return
	}

	projectList, err := projects.ListByOrganizationID(
		r.Context(),
		h.projectRepo,
		organizationID,
		user.IdentityID,
	)
	if err != nil {
		log.Printf(
			"console organization projects list failed identity=%q organization=%q: %v",
			user.IdentityID,
			organizationID,
			err,
		)
		response.Error(w, http.StatusInternalServerError, "Failed to list projects")
		return
	}

	log.Printf(
		"console organization projects list succeeded identity=%q organization=%q count=%d",
		user.IdentityID,
		organizationID,
		len(projectList),
	)

	response.OK(w, "Projects retrieved", projectList)
}

func (h *Handler) CreateOrganizationProjectHandler(w http.ResponseWriter, r *http.Request) {
	r, err := authmiddlewares.VerifyUserToken(h.config.JWTSettings, w, r)
	if err != nil {
		log.Printf("console project create failed: unauthorized: %v", err)
		return
	}

	user, ok := authmiddlewares.UserFromContext(r.Context())
	if !ok {
		log.Printf("console project create failed: missing authenticated user")
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	organizationID := r.PathValue("oid")
	if organizationID == "" {
		response.Error(w, http.StatusBadRequest, "Organization ID is required")
		return
	}

	var req projects.CreateInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	project, err := projects.Create(
		r.Context(),
		h.projectRepo,
		organizationID,
		user.IdentityID,
		req,
	)
	if err != nil {
		log.Printf(
			"console project create failed identity=%q organization=%q: %v",
			user.IdentityID,
			organizationID,
			err,
		)

		switch {
		case errors.Is(err, projects.ErrProjectNameRequired):
			response.Error(w, http.StatusBadRequest, "Project name is required")
		case errors.Is(err, projects.ErrProjectNotFound):
			response.Error(w, http.StatusNotFound, "Organization not found")
		default:
			response.Error(w, http.StatusInternalServerError, "Failed to create project")
		}

		return
	}

	log.Printf(
		"console project create succeeded identity=%q organization=%q project=%q",
		user.IdentityID,
		organizationID,
		project.ID,
	)

	response.Created(w, "Project created", project)
}

func (h *Handler) GetOrganizationProjectHandler(w http.ResponseWriter, r *http.Request) {
	r, err := authmiddlewares.VerifyUserToken(h.config.JWTSettings, w, r)
	if err != nil {
		log.Printf("console project get failed: unauthorized: %v", err)
		return
	}

	user, ok := authmiddlewares.UserFromContext(r.Context())
	if !ok {
		log.Printf("console project get failed: missing authenticated user")
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	organizationID := r.PathValue("oid")
	projectID := r.PathValue("pid")

	if organizationID == "" || projectID == "" {
		response.Error(w, http.StatusBadRequest, "Organization ID and project ID are required")
		return
	}

	project, err := projects.Get(
		r.Context(),
		h.projectRepo,
		organizationID,
		projectID,
		user.IdentityID,
	)
	if err != nil {
		if errors.Is(err, projects.ErrProjectNotFound) {
			response.Error(w, http.StatusNotFound, "Project not found")
			return
		}

		log.Printf(
			"console project get failed identity=%q organization=%q project=%q: %v",
			user.IdentityID,
			organizationID,
			projectID,
			err,
		)
		response.Error(w, http.StatusInternalServerError, "Failed to get project")
		return
	}

	response.OK(w, "Project retrieved", project)
}

func (h *Handler) UpdateOrganizationProjectHandler(w http.ResponseWriter, r *http.Request) {
	r, err := authmiddlewares.VerifyUserToken(h.config.JWTSettings, w, r)
	if err != nil {
		log.Printf("console project update failed: unauthorized: %v", err)
		return
	}

	user, ok := authmiddlewares.UserFromContext(r.Context())
	if !ok {
		log.Printf("console project update failed: missing authenticated user")
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	organizationID := r.PathValue("oid")
	projectID := r.PathValue("pid")

	if organizationID == "" || projectID == "" {
		response.Error(w, http.StatusBadRequest, "Organization ID and project ID are required")
		return
	}

	var req projects.UpdateInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	project, err := projects.Update(
		r.Context(),
		h.projectRepo,
		organizationID,
		projectID,
		user.IdentityID,
		req,
	)
	if err != nil {
		log.Printf(
			"console project update failed identity=%q organization=%q project=%q: %v",
			user.IdentityID,
			organizationID,
			projectID,
			err,
		)

		switch {
		case errors.Is(err, projects.ErrProjectNameRequired):
			response.Error(w, http.StatusBadRequest, "Project name is required")
		case errors.Is(err, projects.ErrProjectNotFound):
			response.Error(w, http.StatusNotFound, "Project not found")
		default:
			response.Error(w, http.StatusInternalServerError, "Failed to update project")
		}

		return
	}

	log.Printf(
		"console project update succeeded identity=%q organization=%q project=%q",
		user.IdentityID,
		organizationID,
		projectID,
	)

	response.OK(w, "Project updated", project)
}

func (h *Handler) DeleteOrganizationProjectHandler(w http.ResponseWriter, r *http.Request) {
	r, err := authmiddlewares.VerifyUserToken(h.config.JWTSettings, w, r)
	if err != nil {
		log.Printf("console project delete failed: unauthorized: %v", err)
		return
	}

	user, ok := authmiddlewares.UserFromContext(r.Context())
	if !ok {
		log.Printf("console project delete failed: missing authenticated user")
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	organizationID := r.PathValue("oid")
	projectID := r.PathValue("pid")

	if organizationID == "" || projectID == "" {
		response.Error(w, http.StatusBadRequest, "Organization ID and project ID are required")
		return
	}

	if err := projects.Delete(
		r.Context(),
		h.projectRepo,
		organizationID,
		projectID,
		user.IdentityID,
	); err != nil {
		log.Printf(
			"console project delete failed identity=%q organization=%q project=%q: %v",
			user.IdentityID,
			organizationID,
			projectID,
			err,
		)

		if errors.Is(err, projects.ErrProjectNotFound) {
			response.Error(w, http.StatusNotFound, "Project not found")
			return
		}

		response.Error(w, http.StatusInternalServerError, "Failed to delete project")
		return
	}

	log.Printf(
		"console project delete succeeded identity=%q organization=%q project=%q",
		user.IdentityID,
		organizationID,
		projectID,
	)

	response.NoContent(w)
}
