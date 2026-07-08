package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"idiom-api-services/api/auth/config"
	"idiom-api-services/api/auth/handlers"
	api "idiom-api-services/api/auth/routes"
	"idiom-api-services/packages/jwt"
	web "idiom-api-services/web/auth/routes"
)

func main() {
	ctx := context.Background()
	appConfig := config.AppConfig{
		JWTSettings: jwt.NewJWTSettings("your-secret-key"),
	}

	authHandler := handlers.NewHandler(appConfig)
	mux := http.NewServeMux()

	// API v1 subrouter
	apiMux := http.NewServeMux()
	api.Mount(apiMux, authHandler)

	webMux := http.NewServeMux()
	web.Mount(webMux)

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", apiMux))
	mux.Handle("/web/", http.StripPrefix("/web", webMux))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	log.Println("Server listening on :8080")
	log.Fatal(server.ListenAndServe())
}
