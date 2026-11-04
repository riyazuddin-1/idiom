package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"

	"idiom-api-services/api/auth/config"
	"idiom-api-services/api/auth/handlers"
	"idiom-api-services/api/auth/middlewares"
	api "idiom-api-services/api/auth/routes"
	web "idiom-api-services/web/auth/routes"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()
	if err := config.Initialize(ctx); err != nil {
		log.Fatal(err)
	}
	defer config.Close()

	mux := http.NewServeMux()

	// API v1 subrouter
	apiMux := http.NewServeMux()
	api.Mount(apiMux)

	webMux := http.NewServeMux()
	web.Mount(webMux)

	apiHandler := middlewares.CORS(http.StripPrefix("/api/v1", apiMux))
	mux.Handle("/api/v1/", apiHandler)
	mux.Handle("/", webMux)

	mux.Handle(
		"GET /.well-known/jwks.json",
		handlers.NewJWKSHandler(),
	)

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
