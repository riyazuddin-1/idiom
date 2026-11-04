package routes

import (
	"idiom-api-services/api/auth/middlewares"
	"idiom-api-services/api/console/config"
	"idiom-api-services/api/console/handlers"
	response "idiom-api-services/packages/responses"
	"net/http"
)

func Mount(mux *http.ServeMux) {
	// Organization routes
	mux.Handle("/organizations", middlewares.Authenticate(config.JWTSettings)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.ListOrganizationsHandler(w, r)
				return
			case http.MethodPost:
				handlers.CreateOrganizationHandler(w, r)
				return
			default:
				response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
		}),
	))

	mux.Handle("/organizations/{oid}",
		middlewares.Authenticate(config.JWTSettings)(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodPut:
					handlers.UpdateOrganizationHandler(w, r)
					return
				case http.MethodDelete:
					handlers.DeleteOrganizationHandler(w, r)
					return
				default:
					response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
				}
			}),
		))

	mux.Handle("PUT /organizations/{oid}/status",
		middlewares.Authenticate(config.JWTSettings)(
			http.HandlerFunc(handlers.UpdateOrganizationStatusHandler),
		),
	)

	// Organization member routes
	mux.Handle("GET /organizations/{oid}/members",
		middlewares.Authenticate(config.JWTSettings)(
			http.HandlerFunc(handlers.ListOrganizationMembersHandler),
		),
	)

	mux.Handle("PUT /organizations/{oid}/members/{uid}/status",
		middlewares.Authenticate(config.JWTSettings)(
			http.HandlerFunc(handlers.UpdateOrganizationMemberStatusHandler),
		),
	)

	mux.Handle("DELETE /organizations/{oid}/members/{uid}",
		middlewares.Authenticate(config.JWTSettings)(
			http.HandlerFunc(handlers.DeleteOrganizationMemberHandler),
		),
	)

	// Organization project routes
	mux.Handle("/organizations/{oid}/projects", middlewares.Authenticate(config.JWTSettings)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.ListOrganizationProjectsHandler(w, r)
				return
			case http.MethodPost:
				handlers.CreateOrganizationProjectHandler(w, r)
				return
			default:
				response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
		}),
	))

	mux.Handle("/organizations/{oid}/projects/{pid}", middlewares.Authenticate(config.JWTSettings)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.GetOrganizationProjectHandler(w, r)
				return
			case http.MethodPut:
				handlers.UpdateOrganizationProjectHandler(w, r)
				return
			case http.MethodDelete:
				handlers.DeleteOrganizationProjectHandler(w, r)
				return
			default:
				response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
			}
		}),
	))
}
