package routes

import (
	"idiom-api-services/api/console/config"
	"idiom-api-services/api/console/handlers"
	response "idiom-api-services/packages/responses"
	"net/http"
)

func Mount(mux *http.ServeMux, appConfig config.AppConfig) {
	handler := handlers.NewHandler(appConfig)
	// Organization routes
	mux.HandleFunc("/organizations", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.ListOrganizationsHandler(w, r)
			return
		case http.MethodPost:
			handler.CreateOrganizationHandler(w, r)
			return
		default:
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	mux.HandleFunc("/organizations/{oid}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			handler.UpdateOrganizationHandler(w, r)
			return
		case http.MethodDelete:
			handler.DeleteOrganizationHandler(w, r)
			return
		default:
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	mux.HandleFunc("/organizations/{oid}/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		handler.UpdateOrganizationStatusHandler(w, r)
	})

	// Organization member routes
	mux.HandleFunc("/organizations/{oid}/members", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		handler.ListOrganizationMembersHandler(w, r)
	})

	mux.HandleFunc("/organizations/{oid}/members/{uid}/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		handler.UpdateOrganizationMemberStatusHandler(w, r)
	})

	mux.HandleFunc("/organizations/{oid}/members/{uid}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		handler.DeleteOrganizationMemberHandler(w, r)
	})

	// Organization project routes
	mux.HandleFunc("/organizations/{oid}/projects", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.ListOrganizationProjectsHandler(w, r)
			return
		case http.MethodPost:
			handler.CreateOrganizationProjectHandler(w, r)
			return
		default:
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	mux.HandleFunc("/organizations/{oid}/projects/{pid}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetOrganizationProjectHandler(w, r)
			return
		case http.MethodPut:
			handler.UpdateOrganizationProjectHandler(w, r)
			return
		case http.MethodDelete:
			handler.DeleteOrganizationProjectHandler(w, r)
			return
		default:
			response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})
}
