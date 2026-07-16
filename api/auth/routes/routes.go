package api

import (
	"net/http"

	"idiom-api-services/api/auth/config"
	"idiom-api-services/api/auth/handlers"
	"idiom-api-services/api/auth/middlewares"
)

func Mount(mux *http.ServeMux, handler *handlers.Handler, appConfig config.AppConfig) {
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.LoginHandler(w, r)
	})

	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		r, err := middlewares.VerifyUserToken(appConfig.JWTSettings, w, r)
		if err != nil {
			return
		}

		handler.LogoutHandler(w, r)
	})

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.RegisterHandler(w, r)
	})

	mux.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.VerifyHandler(w, r)
	})

	mux.HandleFunc("/password-reset", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.SendPasswordResetHandler(w, r)
		case http.MethodPost:
			handler.UpdatePasswordHandler(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
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
		case http.MethodPut:
			handler.UpdateCurrentUserHandler(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/token/refresh", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		handler.RefreshTokenHandler(w, r)
	})
}
