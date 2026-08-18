package handlers

import (
	"html/template"
	"net/http"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "web/auth/templates/login.html")
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "web/auth/templates/register.html")
}

func ForgotPasswordPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "web/auth/templates/forgot-password.html")
}

func PasswordResetPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "web/auth/templates/password-reset.html")
}

func renderTemplate(w http.ResponseWriter, r *http.Request, path string) {
	projectID := r.PathValue("project_id")
	tmpl := template.Must(template.ParseFiles(path))
	_ = tmpl.Execute(w, map[string]string{
		"ProjectID": projectID,
		"Scope":     r.URL.Query().Get("scope"),
	})
}
