package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"idiom-api-services/api/console/config"
	api "idiom-api-services/api/console/routes"
	web "idiom-api-services/web/console/routes"
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

	// Root mux
	mux := http.NewServeMux()

	// API routes
	apiMux := http.NewServeMux()
	api.Mount(apiMux)

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
