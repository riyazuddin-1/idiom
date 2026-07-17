package handlers

import (
	"idiom-api-services/api/console/config"
	response "idiom-api-services/packages/responses"
	"net/http"
)

type Handler struct {
	config config.AppConfig
}

func NewHandler(config config.AppConfig) *Handler {
	return &Handler{
		config: config,
	}
}

func (h *Handler) ListOrganizationsHandler(w http.ResponseWriter, r *http.Request) {
	response.Error(w, http.StatusNotImplemented, "Not implemented")
}

func (h *Handler) CreateOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	response.Error(w, http.StatusNotImplemented, "Not implemented")
}

func (h *Handler) UpdateOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	response.Error(w, http.StatusNotImplemented, "Not implemented")
}

func (h *Handler) UpdateOrganizationStatusHandler(w http.ResponseWriter, r *http.Request) {
	response.Error(w, http.StatusNotImplemented, "Not implemented")
}

func (h *Handler) DeleteOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	response.Error(w, http.StatusNotImplemented, "Not implemented")
}

func (h *Handler) ListOrganizationMembersHandler(w http.ResponseWriter, r *http.Request) {
	response.Error(w, http.StatusNotImplemented, "Not implemented")
}

func (h *Handler) UpdateOrganizationMemberStatusHandler(w http.ResponseWriter, r *http.Request) {
	response.Error(w, http.StatusNotImplemented, "Not implemented")
}

func (h *Handler) DeleteOrganizationMemberHandler(w http.ResponseWriter, r *http.Request) {
	response.Error(w, http.StatusNotImplemented, "Not implemented")
}
