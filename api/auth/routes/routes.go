package routes

import (
	"net/http"

	"idiom-api-services/api/auth/handlers"
)

func Mount(mux *http.ServeMux) {
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handlers.LoginHandler(w, r)
	})

	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handlers.LogoutHandler(w, r)
	})

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handlers.RegisterHandler(w, r)
	})

	mux.HandleFunc("/password-reset", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handlers.PasswordResetHandler(w, r)
	})

	mux.HandleFunc("/me", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetCurrentUserHandler(w, r)
		case http.MethodPut:
			handlers.UpdateCurrentUserHandler(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
}
