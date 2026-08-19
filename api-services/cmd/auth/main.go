package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"idiom-api-services/api/auth/config"
	"idiom-api-services/api/auth/handlers"
	api "idiom-api-services/api/auth/routes"
	"idiom-api-services/packages/database/postgres"
	"idiom-api-services/packages/email"
	"idiom-api-services/packages/jwt"
	web "idiom-api-services/web/auth/routes"
)

func main() {
	ctx := context.Background()
	pg, err := postgres.Init(ctx, config.PostgresDSN)
	if err != nil {
		log.Fatalf("failed to initialize postgres: %v", err)
	}
	defer pg.Close()

	appConfig := config.AppConfig{
		JWTSettings: jwt.NewJWTSettings(
			config.JWTPrivateKey,
			config.JWTPublicKey,
			config.TokenExpirationInSeconds,
			config.TokenIssuer,
		),
		EmailSender:        email.NewHTTPMailSender(),
		AuthBaseURL:        config.AuthBaseURL,
		VerificationSecret: config.VerificationSecret,
		PostgresDB:         pg,
	}

	authHandler := handlers.NewHandler(appConfig)
	mux := http.NewServeMux()

	// API v1 subrouter
	apiMux := http.NewServeMux()
	api.Mount(apiMux, authHandler, appConfig)

	webMux := http.NewServeMux()
	web.Mount(webMux)

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", apiMux))
	mux.Handle("/auth/v1/", http.StripPrefix("/auth/v1", apiMux))
	mux.Handle("/", webMux)

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
