package api

import (
	"net/http"

	"idiom-api-services/api/auth/config"
	"idiom-api-services/api/auth/handlers"
	"idiom-api-services/api/auth/middlewares"
	response "idiom-api-services/packages/responses"
)

func Mount(mux *http.ServeMux, handler *handlers.Handler, appConfig config.AppConfig) {
	mux.HandleFunc("/{project_id}/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		handler.LoginHandler(w, r)
	})

	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		handler.TokenHandler(w, r)
	})

	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		r, err := middlewares.VerifyUserToken(appConfig.JWTSettings, w, r)
		if err != nil {
			response.OK(w, "Logout successful", nil)
			return
		}

		handler.LogoutHandler(w, r)
	})

	mux.HandleFunc("/{project_id}/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		handler.RegisterHandler(w, r)
	})

	mux.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
		handler.VerifyHandler(w, r)
	})

	mux.HandleFunc("/{project_id}/password-reset", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.SendPasswordResetHandler(w, r)
			return
		case http.MethodPost:
			handler.UpdatePasswordHandler(w, r)
			return
		default:
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	mux.HandleFunc("/me", func(w http.ResponseWriter, r *http.Request) {
		r, err := middlewares.VerifyUserToken(appConfig.JWTSettings, w, r)
		if err != nil {
			return
		}

		switch r.Method {
		case http.MethodGet:
			handler.GetCurrentUserHandler(w, r)
			return
		case http.MethodPut:
			handler.UpdateCurrentUserHandler(w, r)
			return
		default:
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	mux.HandleFunc("/token/refresh", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		handler.RefreshTokenHandler(w, r)
	})
}
