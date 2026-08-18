package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"idiom-api-services/api/console/config"
	api "idiom-api-services/api/console/routes"
	"idiom-api-services/packages/database/postgres"
	"idiom-api-services/packages/jwt"
	web "idiom-api-services/web/console/routes"
)

func main() {
	ctx := context.Background()

	// Database
	pg, err := postgres.Init(ctx, config.PostgresDSN)
	if err != nil {
		log.Fatalf("failed to initialize postgres: %v", err)
	}
	defer pg.Close()

	// Shared app configuration
	appConfig := config.AppConfig{
		PostgresDB: pg,
		JWTSettings: jwt.NewJWTSettings(
			"",
			config.JWTPublicKey,
			0,
			config.TokenIssuer,
		),
		AuthBaseURL: config.AuthBaseURL,
		AuthProject: config.AuthProject,
	}

	// Root mux
	mux := http.NewServeMux()

	// API routes
	apiMux := http.NewServeMux()
	api.Mount(apiMux, appConfig)

	// Web routes (React SPA)
	webMux := http.NewServeMux()
	web.Mount(webMux, config.DistDir)

	// Mount routers
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", apiMux))
	mux.Handle("/", webMux)

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// Render provides PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	log.Printf("Console server listening on :%s", port)
	log.Fatal(server.ListenAndServe())
}
