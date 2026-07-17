package web

import (
	response "idiom-api-services/packages/responses"
	"idiom-api-services/web/auth/handlers"
	"net/http"
)

func Mount(mux *http.ServeMux) {
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		handlers.LoginPage(w, r)
	})
	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		handlers.RegisterPage(w, r)
	})
	mux.HandleFunc("/password-reset", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		handlers.PasswordResetPage(w, r)
	})
}
