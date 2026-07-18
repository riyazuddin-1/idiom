package web

import (
	"html/template"
	"idiom-api-services/api/console/config"
	"net/http"
	"strings"
)

func Mount(mux *http.ServeMux, appConfig config.AppConfig) {
	mux.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		http.Redirect(w, r, strings.TrimRight(appConfig.AuthBaseURL, "/")+"/"+appConfig.AuthProject+"/login", http.StatusSeeOther)
	})

	mux.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		tmpl := template.Must(template.ParseFiles("web/console/templates/dashboard.html"))
		_ = tmpl.Execute(w, map[string]string{
			"AuthBaseURL": strings.TrimRight(appConfig.AuthBaseURL, "/"),
			"AuthProject": appConfig.AuthProject,
		})
	})

}
