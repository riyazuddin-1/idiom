package middlewares

import (
	"net/http"
	"strings"
)

func CORS(next http.Handler) http.Handler {
	origins := map[string][]string{
		"idiom": {
			"https://idiom-identity.onrender.com",
			"https://idiom-console.onrender.com",
		},
		"notelet": {
			"https://notelet.vercel.app",
		},
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Handle preflight requests.
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set(
				"Access-Control-Allow-Methods",
				"GET, POST, PUT, PATCH, DELETE, OPTIONS",
			)
			w.Header().Set(
				"Access-Control-Allow-Headers",
				"Content-Type, Authorization, X-Project-ID",
			)
			w.Header().Set("Access-Control-Max-Age", "600")
			w.Header().Add("Vary", "Origin")

			w.WriteHeader(http.StatusNoContent)
			return
		}

		origin := r.Header.Get("Origin")
		projectID := r.Header.Get("X-Project-ID")

		allowed_origins, exists := origins[projectID]
		if !exists {
			http.Error(w, "Invalid project", http.StatusForbidden)
			return
		}

		allowed := false

		for _, allowedOrigin := range allowed_origins {
			if strings.TrimSpace(allowedOrigin) == origin {
				allowed = true
				break
			}
		}

		if !allowed {
			http.Error(w, "Origin not allowed", http.StatusForbidden)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Vary", "Origin")

		next.ServeHTTP(w, r)
	})
}
