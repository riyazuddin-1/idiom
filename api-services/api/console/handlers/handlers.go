package handlers

import (
	"encoding/json"
	authmiddlewares "idiom-api-services/api/auth/middlewares"
	"idiom-api-services/api/console/config"
	apikeys "idiom-api-services/domains/api_keys"
	"idiom-api-services/domains/loadtest"
	"idiom-api-services/domains/org_members"
	"idiom-api-services/domains/organizations"
	"idiom-api-services/domains/projects"
	response "idiom-api-services/packages/responses"
	"log"
	"net/http"
)

type Handler struct {
	config       config.AppConfig
	orgRepo      *organizations.Repository
	memberRepo   *org_members.Repository
	projectRepo  *projects.Repository
	apiKeyRepo   *apikeys.Repository
	loadtestRepo *loadtest.Repository
}

func NewHandler(config config.AppConfig) *Handler {
	return &Handler{
		config:       config,
		orgRepo:      organizations.NewRepository(config.PostgresDB),
		memberRepo:   org_members.NewRepository(config.PostgresDB),
		projectRepo:  projects.NewRepository(config.PostgresDB),
		apiKeyRepo:   apikeys.NewRepository(config.PostgresDB),
		loadtestRepo: loadtest.NewRepository(config.PostgresDB),
	}
}

type loadTestRequest struct {
	ID      string          `json:"id"`
	Payload json.RawMessage `json:"payload"`
}

func (h *Handler) StressTestHandler(w http.ResponseWriter, r *http.Request) {
	r, err := authmiddlewares.VerifyUserToken(h.config.JWTSettings, w, r)
	if err != nil {
		log.Printf("console stress test failed: unauthorized: %v", err)
		return
	}

	var req loadTestRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil || req.ID == "" || len(req.Payload) == 0 || string(req.Payload) == "null" {
		response.Error(w, http.StatusBadRequest, `Request must contain "id" and "payload"`)
		return
	}

	var payload any
	if err := json.Unmarshal(req.Payload, &payload); err != nil {
		response.Error(w, http.StatusBadRequest, "Payload must be valid JSON")
		return
	}

	if err := h.loadtestRepo.EnsureTable(r.Context()); err != nil {
		log.Printf("console stress test failed: create table: %v", err)
		response.Error(w, http.StatusInternalServerError, "Failed to prepare load test table")
		return
	}

	action, err := h.loadtestRepo.Toggle(r.Context(), req.ID, req.Payload)
	if err != nil {
		log.Printf("console stress test failed id=%q: %v", req.ID, err)
		response.Error(w, http.StatusInternalServerError, "Failed to toggle load test item")
		return
	}

	response.OK(w, "Load test item "+action, map[string]string{
		"id":     req.ID,
		"action": action,
	})
}

func (h *Handler) ListOrganizationsHandler(w http.ResponseWriter, r *http.Request) {
	r, err := authmiddlewares.VerifyUserToken(h.config.JWTSettings, w, r)
	if err != nil {
		log.Printf("console organizations list failed: unauthorized: %v", err)
		return
	}

	user, ok := authmiddlewares.UserFromContext(r.Context())
	if !ok {
		log.Printf("console organizations list failed: missing authenticated user")
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	organizations, err := h.orgRepo.ListByIdentityID(r.Context(), user.IdentityID)
	if err != nil {
		log.Printf("console organizations list failed identity=%q: %v", user.IdentityID, err)
		response.Error(w, http.StatusInternalServerError, "Failed to list organizations")
		return
	}

	log.Printf("console organizations list succeeded identity=%q count=%d", user.IdentityID, len(organizations))
	response.OK(w, "Organizations retrieved", organizations)
}
