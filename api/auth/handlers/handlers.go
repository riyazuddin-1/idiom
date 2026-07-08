package handlers

import (
	"encoding/json"
	"idiom-api-services/api/auth/config"
	"idiom-api-services/api/auth/middlewares"
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
	req, err := middlewares.ValidateLoginRequest(w, r)
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
		ProjectID:   dummyProjectID,
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

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully",
		"user":    identity,
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
