package handlers

import (
	"encoding/json"
	"idiom-api-services/api/auth/config"
	"idiom-api-services/api/auth/helpers"
	"idiom-api-services/api/auth/middlewares"
	apikeys "idiom-api-services/domains/api_keys"
	"idiom-api-services/domains/identities"
	"idiom-api-services/domains/projects"
	"idiom-api-services/domains/sessions"
	"idiom-api-services/packages/crypto"
	response "idiom-api-services/packages/responses"
	"log"
	"net/http"
	"strings"
	"time"
)

type Handler struct {
	config      config.AppConfig
	repo        *identities.Repository
	sessionRepo *sessions.Repository
	projectRepo *projects.Repository
	apiKeyRepo  *apikeys.Repository
}

type loginCodeScope struct {
	Operation  string `json:"operation"`
	ProjectID  string `json:"pid"`
	IdentityID string `json:"sub"`
	SessionID  string `json:"sid"`
	ExpiresAt  int64  `json:"exp"`
}

func NewHandler(config config.AppConfig) *Handler {
	return &Handler{
		config:      config,
		repo:        identities.NewRepository(config.PostgresDB),
		sessionRepo: sessions.NewRepository(config.PostgresDB),
		projectRepo: projects.NewRepository(config.PostgresDB),
		apiKeyRepo:  apikeys.NewRepository(config.PostgresDB),
	}
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	routeProjectID := r.PathValue("project_id")
	path := r.URL.Path
	var err error
	r, err = middlewares.VerifyProject(h.projectRepo, w, r)
	if err != nil {
		log.Printf("login project verification failed project=%q path=%q: %v", routeProjectID, path, err)
		return
	}

	req, err := middlewares.ValidateLoginRequest(w, r)
	if err != nil {
		return
	}

	project, ok := middlewares.ProjectFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusInternalServerError, "Project is not configured")
		return
	}
	projectID := project.ID

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

	session, err := sessions.Create(
		r.Context(),
		h.sessionRepo,
		identity,
		r.RemoteAddr,
		r.UserAgent(),
	)
	if err != nil {
		log.Printf("login session start failed project=%q email=%q identity=%q: %v", projectID, req.Email, identity.ID, err)
		response.Error(w, http.StatusInternalServerError, "Failed to start session")
		return
	}

	code, err := crypto.EncryptJSON(loginCodeScope{
		Operation:  config.OperationLogin,
		ProjectID:  identity.ProjectID,
		IdentityID: identity.ID,
		SessionID:  session.ID,
		ExpiresAt:  time.Now().UTC().Add(60 * time.Second).Unix(),
	}, h.config.VerificationSecret)
	if err != nil {
		log.Printf("login code creation failed project=%q email=%q identity=%q session=%q: %v", projectID, req.Email, identity.ID, session.ID, err)
		response.Error(w, http.StatusInternalServerError, "Failed to start session")
		return
	}

	log.Printf("login succeeded project=%q email=%q identity=%q session=%q", projectID, req.Email, identity.ID, session.ID)

	if response.Redirect(w, r, project.RedirectURLs, code) {
		return
	}

	response.OK(w, "Login successful", map[string]interface{}{
		"code": code,
		"user": map[string]string{
			"email": identity.Email,
		},
	})
}

func (h *Handler) TokenHandler(w http.ResponseWriter, r *http.Request) {
	req, err := middlewares.ValidateTokenRequest(w, r)
	if err != nil {
		return
	}

	project, ok, err := projects.ResolveActive(
		r.Context(),
		h.projectRepo,
		req.ProjectID,
	)
	if err != nil {
		log.Printf(
			"token exchange project resolution errored project=%q: %v",
			req.ProjectID,
			err,
		)
		response.Error(
			w,
			http.StatusInternalServerError,
			"Failed to resolve project",
		)
		return
	}

	if !ok {
		log.Printf(
			"token exchange failed project=%q: project not found",
			req.ProjectID,
		)
		response.Error(
			w,
			http.StatusUnauthorized,
			"Invalid authorization code",
		)
		return
	}

	// r, err = middlewares.VerifyAPIKey(
	// 	h.apiKeyRepo,
	// 	project.ID,
	// 	w,
	// 	r,
	// )
	// if err != nil {
	// 	log.Printf(
	// 		"token exchange api key verification failed project=%q: %v",
	// 		project.ID,
	// 		err,
	// 	)
	// 	return
	// }

	var payload loginCodeScope

	if err := crypto.DecryptJSON(
		req.Code,
		h.config.VerificationSecret,
		&payload,
	); err != nil {
		log.Printf(
			"token exchange failed project=%q: invalid code: %v",
			project.ID,
			err,
		)
		response.Error(
			w,
			http.StatusUnauthorized,
			"Invalid authorization code",
		)
		return
	}

	if payload.Operation != config.OperationLogin ||
		payload.ProjectID != project.ID ||
		payload.IdentityID == "" ||
		payload.SessionID == "" ||
		payload.ExpiresAt < time.Now().UTC().Unix() {

		log.Printf(
			"token exchange failed route_project=%q code_project=%q identity=%q session=%q operation=%q",
			project.ID,
			payload.ProjectID,
			payload.IdentityID,
			payload.SessionID,
			payload.Operation,
		)

		response.Error(
			w,
			http.StatusUnauthorized,
			"Invalid authorization code",
		)
		return
	}

	tokens, err := sessions.IssueTokens(
		r.Context(),
		h.sessionRepo,
		h.config.JWTSettings,
		payload.SessionID,
		payload.IdentityID,
		payload.ProjectID,
	)
	if err != nil {
		log.Printf(
			"token exchange issue failed project=%q session=%q identity=%q: %v",
			project.ID,
			payload.SessionID,
			payload.IdentityID,
			err,
		)
		response.Error(
			w,
			http.StatusInternalServerError,
			"Failed to issue tokens",
		)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/api/v1/token/refresh",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	log.Printf(
		"token exchange succeeded project=%q identity=%q session=%q",
		project.ID,
		tokens.IdentityID,
		tokens.SessionID,
	)

	response.OK(w, "Token issued successfully", map[string]interface{}{
		"accessToken":  tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
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

	response.OK(w, "Logout successful", nil)
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	routeProjectID := r.PathValue("project_id")
	path := r.URL.Path
	var err error
	r, err = middlewares.VerifyProject(h.projectRepo, w, r)
	if err != nil {
		log.Printf("register project verification failed project=%q path=%q: %v", routeProjectID, path, err)
		return
	}

	req, err := middlewares.ValidateRegisterRequest(w, r)
	if err != nil {
		return
	}

	project, ok := middlewares.ProjectFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusInternalServerError, "Project is not configured")
		return
	}
	projectID := project.ID

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

	response.Created(w, "User registered successfully", identity)
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

	response.OK(w, "Email verified successfully", nil)
}

func (h *Handler) SendPasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	routeProjectID := r.PathValue("project_id")
	path := r.URL.Path
	var err error
	r, err = middlewares.VerifyProject(h.projectRepo, w, r)
	if err != nil {
		log.Printf("password reset request project verification failed project=%q path=%q: %v", routeProjectID, path, err)
		return
	}

	req, err := middlewares.ValidatePasswordResetEmailRequest(w, r)
	if err != nil {
		return
	}

	project, ok := middlewares.ProjectFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusInternalServerError, "Project is not configured")
		return
	}
	projectID := project.ID

	if err := helpers.SendPasswordResetEmail(r.Context(), h.config, projectID, req.Email); err != nil {
		log.Printf("failed to send password reset email to %s: %v", req.Email, err)
		response.Error(w, http.StatusInternalServerError, "Failed to send password reset email")
		return
	}

	log.Printf("password reset request succeeded project=%q email=%q", projectID, req.Email)

	response.OK(w, "Password reset email sent", nil)
}

func (h *Handler) UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	routeProjectID := r.PathValue("project_id")
	path := r.URL.Path
	var err error
	r, err = middlewares.VerifyProject(h.projectRepo, w, r)
	if err != nil {
		log.Printf("password update project verification failed project=%q path=%q: %v", routeProjectID, path, err)
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

	project, ok := middlewares.ProjectFromContext(r.Context())
	if !ok || payload.ProjectID != project.ID {
		projectID := ""
		if ok {
			projectID = project.ID
		}
		log.Printf("password update failed route_project=%q scope_project=%q email=%q", projectID, payload.ProjectID, payload.Email)
		response.Error(w, http.StatusBadRequest, "Invalid password reset scope")
		return
	}
	projectID := project.ID

	if err := identities.UpdatePassword(r.Context(), h.repo, projectID, payload.Email, req.Password); err != nil {
		log.Printf("password update failed project=%q email=%q: %v", projectID, payload.Email, err)
		response.Error(w, http.StatusInternalServerError, "Failed to update password")
		return
	}

	log.Printf("password update succeeded project=%q email=%q", projectID, payload.Email)

	response.OK(w, "Password updated successfully", nil)
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
	var refreshTokenValue string

	cookie, err := r.Cookie("refresh_token")
	if err == nil && cookie.Value != "" {
		refreshTokenValue = cookie.Value
	}

	if refreshTokenValue == "" {
		var body struct {
			RefreshToken string `json:"refresh_token"`
		}
		if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
			if err := json.NewDecoder(r.Body).Decode(&body); err == nil {
				refreshTokenValue = body.RefreshToken
			}
		} else if err := r.ParseForm(); err == nil {
			refreshTokenValue = r.FormValue("refresh_token")
		}
	}

	if refreshTokenValue == "" {
		log.Printf("token refresh failed: missing refresh token")
		response.Error(w, http.StatusUnauthorized, "Refresh token not found")
		return
	}

	tokens, err := sessions.Refresh(r.Context(), h.sessionRepo, h.config.JWTSettings, refreshTokenValue)

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

	response.OK(w, "Token refreshed successfully", map[string]interface{}{
		"accessToken":  tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
	})
}

// OAuth
