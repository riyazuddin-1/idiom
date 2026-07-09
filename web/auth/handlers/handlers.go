package handlers

import (
	"html/template"
	"net/http"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/auth/templates/login.html")
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/auth/templates/register.html")
}

func PasswordResetPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/auth/templates/password-reset.html"))
	_ = tmpl.Execute(w, map[string]string{
		"Scope": r.URL.Query().Get("scope"),
	})
}
