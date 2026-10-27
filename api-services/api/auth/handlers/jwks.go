package handlers

import (
	"encoding/json"
	"idiom-api-services/api/auth/config"
	response "idiom-api-services/packages/responses"
	"net/http"
)

func NewJWKSHandler(config config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if config.JWTSettings == nil {
			http.Error(w, "Signing key unavailable", http.StatusInternalServerError)
			return
		}

		Keys, err := config.JWTSettings.PublicJWKS()
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=300")

		if err := json.NewEncoder(w).Encode(Keys); err != nil {
			return
		}
	}
}
