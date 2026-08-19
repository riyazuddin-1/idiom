package middlewares

import (
	"net/http"
	"strings"
)

var defaultAllowedOrigins = map[string]struct{}{
	"https://idiom-console.onrender.com": {},
}

// CORS adds the headers required by browser clients and handles preflight
// requests before they reach individual API route handlers.
func CORS(allowedOrigins string, next http.Handler) http.Handler {
	configuredOrigins := parseOrigins(allowedOrigins)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if isAllowedOrigin(origin, configuredOrigins) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Add("Vary", "Origin")
		}

		if r.Method == http.MethodOptions {
			if origin == "" || !isAllowedOrigin(origin, configuredOrigins) {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "600")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func parseOrigins(value string) map[string]struct{} {
	origins := make(map[string]struct{}, len(defaultAllowedOrigins))
	for origin := range defaultAllowedOrigins {
		origins[origin] = struct{}{}
	}

	for _, origin := range strings.Split(value, ",") {
		origin = strings.TrimSpace(strings.TrimRight(origin, "/"))
		if origin != "" {
			origins[origin] = struct{}{}
		}
	}

	return origins
}

func isAllowedOrigin(origin string, allowedOrigins map[string]struct{}) bool {
	_, ok := allowedOrigins[strings.TrimRight(origin, "/")]
	return ok
}
