package middlewares

import (
	"encoding/json"
	"errors"
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

func ValidateLoginRequest(w http.ResponseWriter, r *http.Request) (*LoginRequest, error) {
	var req LoginRequest

	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return nil, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return nil, err
		}
		req.Email = r.FormValue("email")
		req.Password = r.FormValue("password")
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return nil, errors.New("missing email or password")
	}

	email, err := normalizeEmail(req.Email)
	if err != nil {
		http.Error(w, "Invalid email address", http.StatusBadRequest)
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
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return nil, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
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
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return nil, errors.New("missing email or password")
	}

	email, err := normalizeEmail(req.Email)
	if err != nil {
		http.Error(w, "Invalid email address", http.StatusBadRequest)
		return nil, err
	}
	req.Email = email

	return &req, nil
}

func ValidatePasswordResetEmailRequest(w http.ResponseWriter, r *http.Request) (*PasswordResetEmailRequest, error) {
	email, err := normalizeEmail(r.URL.Query().Get("email"))
	if err != nil {
		http.Error(w, "Invalid email address", http.StatusBadRequest)
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
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return nil, err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return nil, err
		}
		req.Scope = r.FormValue("scope")
		req.Password = r.FormValue("password")
		req.ConfirmPassword = r.FormValue("confirm_password")
	}

	if req.Scope == "" || req.Password == "" || req.ConfirmPassword == "" {
		http.Error(w, "Scope, password, and confirm password are required", http.StatusBadRequest)
		return nil, errors.New("missing password reset fields")
	}

	if req.Password != req.ConfirmPassword {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return nil, errors.New("passwords do not match")
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
