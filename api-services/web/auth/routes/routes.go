package web

import (
	response "idiom-api-services/packages/responses"
	"idiom-api-services/web/auth/handlers"
	"net/http"
)

func Mount(mux *http.ServeMux) {
	mux.HandleFunc("/{project_id}/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		handlers.LoginPage(w, r)
	})
	mux.HandleFunc("/{project_id}/logout", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("/{project_id}/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		handlers.RegisterPage(w, r)
	})
	mux.HandleFunc("/{project_id}/forgot-password", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		handlers.ForgotPasswordPage(w, r)
	})
	mux.HandleFunc("/{project_id}/password-reset", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		handlers.PasswordResetPage(w, r)
	})
}
