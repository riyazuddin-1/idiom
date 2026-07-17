package web

import (
	"idiom-api-services/api/auth/middlewares"
	response "idiom-api-services/packages/responses"
	"idiom-api-services/web/auth/handlers"
	"net/http"
)

func Mount(mux *http.ServeMux) {
	withProject := func(r *http.Request) *http.Request {
		projectID := r.PathValue("project_id")
		return r.WithContext(middlewares.ContextWithProjectID(r.Context(), projectID))
	}

	mux.HandleFunc("/{project_id}/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		r = withProject(r)
		handlers.LoginPage(w, r)
	})
	mux.HandleFunc("/{project_id}/logout", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("/{project_id}/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		r = withProject(r)
		handlers.RegisterPage(w, r)
	})
	mux.HandleFunc("/{project_id}/password-reset", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		r = withProject(r)
		handlers.PasswordResetPage(w, r)
	})
}
