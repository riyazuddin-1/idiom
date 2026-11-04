package handlers

import (
	authmiddlewares "idiom-api-services/api/auth/middlewares"
	"idiom-api-services/api/console/config"
	response "idiom-api-services/packages/responses"
	"log"
	"net/http"
)

func ListOrganizationsHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := authmiddlewares.UserFromContext(r.Context())
	if !ok {
		log.Printf("console organizations list failed: missing authenticated user")
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	organizations, err := config.OrganizationRepo.ListByIdentityID(r.Context(), user.IdentityID)
	if err != nil {
		log.Printf("console organizations list failed identity=%q: %v", user.IdentityID, err)
		response.Error(w, http.StatusInternalServerError, "Failed to list organizations")
		return
	}

	log.Printf("console organizations list succeeded identity=%q count=%d", user.IdentityID, len(organizations))
	response.OK(w, "Organizations retrieved", organizations)
}
