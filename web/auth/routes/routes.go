package web

import (
	"idiom-api-services/web/auth/handlers"
	"net/http"
)

func Mount(mux *http.ServeMux) {
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid route", http.StatusMethodNotAllowed)
			return
		}

		handlers.LoginPage(w, r)
	})
	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid route", http.StatusMethodNotAllowed)
			return
		}

		handlers.RegisterPage(w, r)
	})
}
