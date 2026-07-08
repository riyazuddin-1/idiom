package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"idiom-api-services/api/auth/config"
	"idiom-api-services/api/auth/handlers"
	routes "idiom-api-services/api/auth/routes"
	"idiom-api-services/packages/jwt"
)

func main() {
	ctx := context.Background()
	appConfig := config.AppConfig{
		JWTSettings: jwt.NewJWTSettings("your-secret-key"),
	}

	authHandler := handlers.NewHandler(appConfig)
	mux := http.NewServeMux()

	// API v1 subrouter
	api := http.NewServeMux()
	routes.Mount(api, authHandler)

	// Mount under /api/v1/
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", api))

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
