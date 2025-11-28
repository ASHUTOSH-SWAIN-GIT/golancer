package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	port := os.Getenv("PORT")
	log.Printf("Request received on port %s for URL: %s", port, r.URL.Path)

	// Simulate some work being done
	fmt.Fprintf(w, "Hello from Backend Server running on port %s!\n", port)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Our Health Checker probes this endpoint and expects a 200 OK
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Error: PORT environment variable not set.")
	}

	http.HandleFunc("/", handler)
	http.HandleFunc("/health", healthCheckHandler)

	log.Printf("Starting Backend Server on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed on port %s: %v", port, err)
	}
}
