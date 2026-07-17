package handlers

import (
	"html/template"
	"idiom-api-services/api/auth/middlewares"
	"net/http"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "web/auth/templates/login.html")
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "web/auth/templates/register.html")
}

func PasswordResetPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r, "web/auth/templates/password-reset.html")
}

func renderTemplate(w http.ResponseWriter, r *http.Request, path string) {
	projectID, _ := middlewares.ProjectIDFromContext(r.Context())
	tmpl := template.Must(template.ParseFiles(path))
	_ = tmpl.Execute(w, map[string]string{
		"ProjectID": projectID,
		"Scope":     r.URL.Query().Get("scope"),
	})
}
