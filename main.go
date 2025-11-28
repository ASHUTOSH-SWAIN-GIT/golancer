package main

import (
	"log"
	"net/url"
	"time"

	"golancer.ashutosh.net/internal/balancer"
)

const LbPort = "8080"

const HealthCheckInterval = 10 * time.Second

var backendURLs = []string{
	"http://localhost:8081",
	"http://localhost:8082",
	"http://localhost:8083",
}

func main() {
	pool := &balancer.ServerPool{}

	for _, rawURL := range backendURLs {
		parsedURL, err := url.Parse(rawURL)
		if err != nil {
			log.Fatalf("Invalid backend URL: %s, Error: %v", rawURL, err)
		}
		backend := &balancer.Backend{
			URL: parsedURL,
		}
		backend.Setalive(true)
		pool.AddBackend(backend)
		log.Printf("INFO: Added backend: %s", rawURL)
	}

	balancer.StartHealthCheckLoop(pool, HealthCheckInterval)
	log.Printf("INFO: Health check started with interval: %v", HealthCheckInterval)

}
