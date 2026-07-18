package handlers

import (
	"idiom-api-services/api/auth/config"
	"idiom-api-services/api/auth/helpers"
	"idiom-api-services/api/auth/middlewares"
	"idiom-api-services/domains/identities"
	"idiom-api-services/domains/projects"
	"idiom-api-services/domains/sessions"
	"idiom-api-services/packages/crypto"
	response "idiom-api-services/packages/responses"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	config      config.AppConfig
	repo        *identities.Repository
	sessionRepo *sessions.Repository
	projectRepo *projects.Repository
}

func NewHandler(config config.AppConfig) *Handler {
	return &Handler{
		config:      config,
		repo:        identities.NewRepository(config.PostgresDB),
		sessionRepo: sessions.NewRepository(config.PostgresDB),
		projectRepo: projects.NewRepository(config.PostgresDB),
	}
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	r, err = middlewares.VerifyProject(h.projectRepo, w, r)
	if err != nil {
		log.Printf("login project verification failed project=%q path=%q: %v", r.PathValue("project_id"), r.URL.Path, err)
		return
	}

	req, err := middlewares.ValidateLoginRequest(w, r)
	if err != nil {
		return
	}

	projectID, ok := middlewares.ProjectIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusInternalServerError, "Project is not configured")
		return
	}

	identity, ok, err := identities.Login(r.Context(), h.repo, projectID, req.Email, req.Password)
	if err != nil || !ok {
		log.Printf("login failed project=%q email=%q", projectID, req.Email)
		response.Error(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if h.config.JWTSettings == nil {
		response.Error(w, http.StatusInternalServerError, "JWT settings are not configured")
		return
	}

	tokens, err := sessions.Start(
		r.Context(),
		h.sessionRepo,
		h.config.JWTSettings,
		identity,
		r.RemoteAddr,
		r.UserAgent(),
	)
	if err != nil {
		log.Printf("login session start failed project=%q email=%q identity=%q: %v", projectID, req.Email, identity.ID, err)
		response.Error(w, http.StatusInternalServerError, "Failed to start session")
		return
	}

	log.Printf("login succeeded project=%q email=%q identity=%q session=%q", projectID, req.Email, identity.ID, tokens.SessionID)

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/api/v1/token/refresh",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"message":     "Login successful",
		"accessToken": tokens.AccessToken,
		"user": map[string]string{
			"email": identity.Email,
		},
	})
}

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := middlewares.UserFromContext(r.Context())

	if !ok {
		log.Printf("logout failed: missing authenticated user")
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	revoked, err := sessions.Revoke(r.Context(), h.sessionRepo, user.SessionID)

	if err != nil || !revoked {
		log.Printf("logout failed project=%q identity=%q session=%q: %v", user.ProjectID, user.IdentityID, user.SessionID, err)
		response.Error(w, http.StatusInternalServerError, "Failed to logout")
		return
	}

	log.Printf("logout succeeded project=%q identity=%q session=%q", user.ProjectID, user.IdentityID, user.SessionID)

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/v1/token/refresh",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})

	response.Message(w, http.StatusOK, "Logout successful")
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	r, err = middlewares.VerifyProject(h.projectRepo, w, r)
	if err != nil {
		log.Printf("register project verification failed project=%q path=%q: %v", r.PathValue("project_id"), r.URL.Path, err)
		return
	}

	req, err := middlewares.ValidateRegisterRequest(w, r)
	if err != nil {
		return
	}

	projectID, ok := middlewares.ProjectIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusInternalServerError, "Project is not configured")
		return
	}

	identity, err := identities.Register(r.Context(), h.repo, identities.RegisterInput{
		ProjectID:   projectID,
		Email:       req.Email,
		Password:    req.Password,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		DisplayName: req.DisplayName,
		AvatarURL:   req.AvatarURL,
	})
	if err != nil {
		log.Printf("register failed project=%q email=%q: %v", projectID, req.Email, err)
		response.Error(w, http.StatusInternalServerError, "Failed to register user")
		return
	}

	log.Printf("register succeeded project=%q email=%q identity=%q", projectID, req.Email, identity.ID)

	// Temporarily disabled until production email delivery is configured.
	// if err := helpers.SendVerificationEmail(r.Context(), h.config, identity); err != nil {
	// 	log.Printf("failed to send verification email to %s: %v", identity.Email, err)
	// 	response.Error(w, http.StatusInternalServerError, "Failed to send verification email")
	// 	return
	// }

	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"message": "User registered successfully.",
		"user":    identity,
	})
}

func (h *Handler) VerifyHandler(w http.ResponseWriter, r *http.Request) {
	scope := r.URL.Query().Get("scope")
	if scope == "" {
		log.Printf("email verification failed: missing scope")
		response.Error(w, http.StatusBadRequest, "Missing verification scope")
		return
	}

	var payload helpers.Scope
	if err := crypto.DecryptJSON(scope, h.config.VerificationSecret, &payload); err != nil {
		log.Printf("email verification failed: invalid scope: %v", err)
		response.Error(w, http.StatusBadRequest, "Invalid verification scope")
		return
	}

	if payload.Operation != config.OperationEmailVerification {
		log.Printf("email verification failed project=%q email=%q: invalid operation %q", payload.ProjectID, payload.Email, payload.Operation)
		response.Error(w, http.StatusBadRequest, "Invalid verification scope")
		return
	}

	if payload.ProjectID == "" {
		log.Printf("email verification failed email=%q: missing project", payload.Email)
		response.Error(w, http.StatusBadRequest, "Invalid verification scope")
		return
	}

	if err := identities.VerifyEmail(r.Context(), h.repo, payload.ProjectID, payload.Email); err != nil {
		log.Printf("email verification failed project=%q email=%q: %v", payload.ProjectID, payload.Email, err)
		response.Error(w, http.StatusInternalServerError, "Failed to verify email")
		return
	}

	log.Printf("email verification succeeded project=%q email=%q", payload.ProjectID, payload.Email)

	response.Message(w, http.StatusOK, "Email verified successfully")
}

func (h *Handler) SendPasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	r, err = middlewares.VerifyProject(h.projectRepo, w, r)
	if err != nil {
		log.Printf("password reset request project verification failed project=%q path=%q: %v", r.PathValue("project_id"), r.URL.Path, err)
		return
	}

	req, err := middlewares.ValidatePasswordResetEmailRequest(w, r)
	if err != nil {
		return
	}

	projectID, ok := middlewares.ProjectIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusInternalServerError, "Project is not configured")
		return
	}

	if err := helpers.SendPasswordResetEmail(r.Context(), h.config, projectID, req.Email); err != nil {
		log.Printf("failed to send password reset email to %s: %v", req.Email, err)
		response.Error(w, http.StatusInternalServerError, "Failed to send password reset email")
		return
	}

	log.Printf("password reset request succeeded project=%q email=%q", projectID, req.Email)

	response.Message(w, http.StatusOK, "Password reset email sent")
}

func (h *Handler) UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	r, err = middlewares.VerifyProject(h.projectRepo, w, r)
	if err != nil {
		log.Printf("password update project verification failed project=%q path=%q: %v", r.PathValue("project_id"), r.URL.Path, err)
		return
	}

	req, err := middlewares.ValidatePasswordResetRequest(w, r)
	if err != nil {
		return
	}

	var payload helpers.Scope
	if err := crypto.DecryptJSON(req.Scope, h.config.VerificationSecret, &payload); err != nil {
		log.Printf("password update failed: invalid scope: %v", err)
		response.Error(w, http.StatusBadRequest, "Invalid password reset scope")
		return
	}

	if payload.Operation != config.OperationPasswordReset {
		log.Printf("password update failed project=%q email=%q: invalid operation %q", payload.ProjectID, payload.Email, payload.Operation)
		response.Error(w, http.StatusBadRequest, "Invalid password reset scope")
		return
	}

	projectID, ok := middlewares.ProjectIDFromContext(r.Context())
	if !ok || payload.ProjectID != projectID {
		log.Printf("password update failed route_project=%q scope_project=%q email=%q", projectID, payload.ProjectID, payload.Email)
		response.Error(w, http.StatusBadRequest, "Invalid password reset scope")
		return
	}

	if err := identities.UpdatePassword(r.Context(), h.repo, projectID, payload.Email, req.Password); err != nil {
		log.Printf("password update failed project=%q email=%q: %v", projectID, payload.Email, err)
		response.Error(w, http.StatusInternalServerError, "Failed to update password")
		return
	}

	log.Printf("password update succeeded project=%q email=%q", projectID, payload.Email)

	response.Message(w, http.StatusOK, "Password updated successfully")
}

func (h *Handler) GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	if user, ok := middlewares.UserFromContext(r.Context()); ok {
		log.Printf("current user read requested project=%q identity=%q session=%q", user.ProjectID, user.IdentityID, user.SessionID)
	}
	// Implement logic to get current user information here
}

func (h *Handler) UpdateCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	if user, ok := middlewares.UserFromContext(r.Context()); ok {
		log.Printf("current user update requested project=%q identity=%q session=%q", user.ProjectID, user.IdentityID, user.SessionID)
	}
	// Implement logic to update current user information here
}

// Session management
func (h *Handler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		log.Printf("token refresh failed: missing refresh token")
		response.Error(w, http.StatusUnauthorized, "Refresh token not found")
		return
	}

	tokens, err := sessions.Refresh(r.Context(), h.sessionRepo, h.config.JWTSettings, refreshToken.Value)

	if err != nil {
		log.Printf("token refresh failed: %v", err)
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/api/v1/token/refresh",
			HttpOnly: true,
			Secure:   true,
			MaxAge:   -1,
			Expires:  time.Unix(0, 0),
		})
		response.Error(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	log.Printf("token refresh succeeded project=%q identity=%q session=%q", tokens.ProjectID, tokens.IdentityID, tokens.SessionID)

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/api/v1/token/refresh",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"message":     "Token refreshed successfully",
		"accessToken": tokens.AccessToken,
	})
}

// OAuth
