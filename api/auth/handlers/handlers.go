package handlers

import (
	"encoding/json"
	"idiom-api-services/api/auth/config"
	"idiom-api-services/api/auth/helpers"
	"idiom-api-services/api/auth/middlewares"
	"idiom-api-services/domains/identities"
	"idiom-api-services/packages/crypto"
	"log"
	"net/http"
	"strings"
)

type Handler struct {
	config config.AppConfig
}

type verificationScope struct {
	ProjectID string `json:"project_id"`
	Email     string `json:"email"`
}

func NewHandler(config config.AppConfig) *Handler {
	return &Handler{
		config: config,
	}
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	req, err := middlewares.ValidateLoginRequest(w, r)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	ok, err := identities.Login(r.Context(), req.Email, req.Password, config.ProjectID)
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

	token, err := h.config.JWTSettings.CreateToken(req.Email, config.ProjectID)
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

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Implement logout logic here
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	req, err := middlewares.ValidateRegisterRequest(w, r)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	identity, err := identities.Register(r.Context(), identities.RegisterInput{
		Email:       req.Email,
		Password:    req.Password,
		ProjectID:   config.ProjectID,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		DisplayName: req.DisplayName,
		AvatarURL:   req.AvatarURL,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to register user",
		})
		return
	}

	if err := helpers.SendVerificationEmail(r.Context(), h.config, identity); err != nil {
		log.Printf("failed to send verification email to %s: %v", identity.Email, err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to send verification email",
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully. Please verify your email.",
		"user":    identity,
	})
}

func (h *Handler) VerifyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	scope := strings.TrimSpace(r.URL.Query().Get("scope"))
	if scope == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Missing verification scope",
		})
		return
	}

	var payload verificationScope
	if err := crypto.DecryptJSON(scope, h.config.VerificationSecret, &payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid verification scope",
		})
		return
	}

	if err := identities.VerifyEmail(r.Context(), payload.Email, payload.ProjectID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to verify email",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Email verified successfully",
	})
}

func (h *Handler) PasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	// Implement password reset logic here
}

func (h *Handler) GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	// Implement logic to get current user information here
}

func (h *Handler) UpdateCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	// Implement logic to update current user information here
}
