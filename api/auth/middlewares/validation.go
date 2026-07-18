package middlewares

import (
	"encoding/json"
	"errors"
	response "idiom-api-services/packages/responses"
	"net/http"
	"net/mail"
	"strings"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

type PasswordResetEmailRequest struct {
	Email string
}

type PasswordResetRequest struct {
	Scope           string `json:"scope"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type TokenRequest struct {
	Code      string `json:"code"`
	ProjectID string `json:"project_id"`
	APIKey    string `json:"api_key"`
}

func ValidateLoginRequest(w http.ResponseWriter, r *http.Request) (*LoginRequest, error) {
	var req LoginRequest

	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid request payload")
			return nil, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid request payload")
			return nil, err
		}
		req.Email = r.FormValue("email")
		req.Password = r.FormValue("password")
	}

	if req.Email == "" || req.Password == "" {
		response.Error(w, http.StatusBadRequest, "Email and password are required")
		return nil, errors.New("missing email or password")
	}

	email, err := normalizeEmail(req.Email)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid email address")
		return nil, err
	}
	req.Email = email

	return &req, nil
}

func ValidateRegisterRequest(w http.ResponseWriter, r *http.Request) (*RegisterRequest, error) {
	var req RegisterRequest

	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid request payload")
			return nil, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid request payload")
			return nil, err
		}
		req.Email = r.FormValue("email")
		req.Password = r.FormValue("password")
		req.FirstName = r.FormValue("first_name")
		req.LastName = r.FormValue("last_name")
		req.DisplayName = r.FormValue("display_name")
		req.AvatarURL = r.FormValue("avatar_url")
	}

	if req.Email == "" || req.Password == "" {
		response.Error(w, http.StatusBadRequest, "Email and password are required")
		return nil, errors.New("missing email or password")
	}

	email, err := normalizeEmail(req.Email)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid email address")
		return nil, err
	}
	req.Email = email

	return &req, nil
}

func ValidatePasswordResetEmailRequest(w http.ResponseWriter, r *http.Request) (*PasswordResetEmailRequest, error) {
	email, err := normalizeEmail(r.URL.Query().Get("email"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid email address")
		return nil, err
	}

	return &PasswordResetEmailRequest{
		Email: email,
	}, nil
}

func ValidatePasswordResetRequest(w http.ResponseWriter, r *http.Request) (*PasswordResetRequest, error) {
	var req PasswordResetRequest

	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid request payload")
			return nil, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid request payload")
			return nil, err
		}
		req.Scope = r.FormValue("scope")
		req.Password = r.FormValue("password")
		req.ConfirmPassword = r.FormValue("confirm_password")
	}

	if req.Scope == "" || req.Password == "" || req.ConfirmPassword == "" {
		response.Error(w, http.StatusBadRequest, "Scope, password, and confirm password are required")
		return nil, errors.New("missing password reset fields")
	}

	if req.Password != req.ConfirmPassword {
		response.Error(w, http.StatusBadRequest, "Passwords do not match")
		return nil, errors.New("passwords do not match")
	}

	return &req, nil
}

func ValidateTokenRequest(w http.ResponseWriter, r *http.Request) (*TokenRequest, error) {
	var req TokenRequest

	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid request payload")
			return nil, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid request payload")
			return nil, err
		}
		req.Code = r.FormValue("code")
		req.ProjectID = r.FormValue("project_id")
		req.APIKey = r.FormValue("api_key")
	}

	req.Code = strings.TrimSpace(req.Code)
	req.ProjectID = strings.TrimSpace(req.ProjectID)
	req.APIKey = strings.TrimSpace(req.APIKey)

	if req.Code == "" || req.ProjectID == "" {
		response.Error(w, http.StatusBadRequest, "Code and project_id are required")
		return nil, errors.New("missing token exchange fields")
	}

	return &req, nil
}

func normalizeEmail(email string) (string, error) {
	normalized := strings.ToLower(email)

	address, err := mail.ParseAddress(normalized)
	if err != nil {
		return "", err
	}

	if address.Name != "" || address.Address != normalized {
		return "", errors.New("email must be a plain email address")
	}

	return normalized, nil
}
