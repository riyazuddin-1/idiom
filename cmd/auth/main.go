package main

import (
	"log"
	"net/http"

	routes "idiom-api-services/api/auth/routes"
)

func main() {
	mux := http.NewServeMux()

	// API v1 subrouter
	api := http.NewServeMux()
	routes.Mount(api)

	// Mount under /api/v1/
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", api))

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
