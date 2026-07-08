package handlers

import "net/http"

func LoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/auth/templates/login.html")
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/auth/templates/register.html")
}
