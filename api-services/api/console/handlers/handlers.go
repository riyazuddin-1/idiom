package handlers

import (
	authmiddlewares "idiom-api-services/api/auth/middlewares"
	"idiom-api-services/api/console/config"
	apikeys "idiom-api-services/domains/api_keys"
	"idiom-api-services/domains/org_members"
	"idiom-api-services/domains/organizations"
	"idiom-api-services/domains/projects"
	response "idiom-api-services/packages/responses"
	"log"
	"net/http"
)

type Handler struct {
	config      config.AppConfig
	orgRepo     *organizations.Repository
	memberRepo  *org_members.Repository
	projectRepo *projects.Repository
	apiKeyRepo  *apikeys.Repository
}

func NewHandler(config config.AppConfig) *Handler {
	return &Handler{
		config:      config,
		orgRepo:     organizations.NewRepository(config.PostgresDB),
		memberRepo:  org_members.NewRepository(config.PostgresDB),
		projectRepo: projects.NewRepository(config.PostgresDB),
		apiKeyRepo:  apikeys.NewRepository(config.PostgresDB),
	}
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
