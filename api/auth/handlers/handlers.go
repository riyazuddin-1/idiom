package handlers

import (
	"encoding/json"
	"idiom-api-services/api/auth/config"
	"idiom-api-services/api/auth/middlewares/validation"
	"idiom-api-services/domains/identities"
	"net/http"
)

const (
	dummyProjectID = "projectID"
)

type Handler struct {
	config config.AppConfig
}

func NewHandler(config config.AppConfig) *Handler {
	return &Handler{
		config: config,
	}
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	req, err := validation.ValidateLoginRequest(w, r)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	ok, err := identities.Login(r.Context(), req.Email, req.Password, dummyProjectID)
	if err != nil || !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid email or password",
		})
		return
	}

	if h.config.JWTSettings == nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "JWT settings are not configured",
		})
		return
	}

	token, err := h.config.JWTSettings.CreateToken(req.Email, dummyProjectID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to create token",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"token":   token,
		"user": map[string]string{
			"email": req.Email,
		},
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Implement logout logic here
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Implement registration logic here
}

func PasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	// Implement password reset logic here
}

func GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	// Implement logic to get current user information here
}

func UpdateCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	// Implement logic to update current user information here
}
