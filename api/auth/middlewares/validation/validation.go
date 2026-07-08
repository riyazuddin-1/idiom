package validation

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

	return &req, nil
}
