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

	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

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

	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.DisplayName = strings.TrimSpace(req.DisplayName)
	req.AvatarURL = strings.TrimSpace(req.AvatarURL)

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

func normalizeEmail(email string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(email))

	address, err := mail.ParseAddress(normalized)
	if err != nil {
		return "", err
	}

	if address.Name != "" || address.Address != normalized {
		return "", errors.New("email must be a plain email address")
	}

	return normalized, nil
}
