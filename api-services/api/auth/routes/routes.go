package api

import (
	"net/http"

	"idiom-api-services/api/auth/config"
	"idiom-api-services/api/auth/handlers"
	"idiom-api-services/api/auth/middlewares"
	response "idiom-api-services/packages/responses"
)

func Mount(mux *http.ServeMux) {
	mux.HandleFunc("/{project_id}/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		handlers.LoginHandler(w, r)
	})

	mux.HandleFunc("/session/token", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		handlers.SessionTokenHandler(w, r)
	})

	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		handlers.TokenHandler(w, r)
	})

	mux.Handle(
		"POST /logout",
		middlewares.Authenticate(config.JWTSettings)(
			http.HandlerFunc(handlers.LogoutHandler),
		),
	)

	mux.HandleFunc("/{project_id}/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		handlers.RegisterHandler(w, r)
	})

	mux.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		handlers.VerifyHandler(w, r)
	})

	mux.HandleFunc("/{project_id}/password-reset", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.SendPasswordResetHandler(w, r)
			return
		case http.MethodPost:
			handlers.UpdatePasswordHandler(w, r)
			return
		default:
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	mux.Handle("/me", middlewares.Authenticate(config.JWTSettings)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetCurrentUserHandler(w, r)
			case http.MethodPut:
				handlers.UpdateCurrentUserHandler(w, r)
			default:
				response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
		}),
	))
}
