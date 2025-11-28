package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
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

	// setting up the backend servers
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

	// the request director
	proxyDirector := func(req *http.Request) {
		targetBackend := pool.GetNextHealthyServer()

		if targetBackend == nil {
			req.Host = ""
			return
		}
		targetURL := targetBackend.URL

		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host

		req.Header.Set("X-Forwarded-For", req.RemoteAddr)
		req.Header.Set("X-Forwarded-Host", req.Host)

		req.Host = targetURL.Host
		targetBackend.IncrementActiveConnections()

	}

	proxy := &httputil.ReverseProxy{
		Director: proxyDirector,
		ModifyResponse: func(r *http.Response) error {
			host := r.Request.Host
			for _, b := range pool.GetBackends() {
				if b.URL.Host == host {
					b.DecrementActiveConnections()
					break
				}
			}
			return nil
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("ERROR: Proxy failed for request to %s: %v", r.URL.Host, err)
			w.WriteHeader(http.StatusBadGateway) // 502 error
			_, _ = w.Write([]byte("502 Bad Gateway: No healthy backend available."))
		},
	}

	log.Printf("INFO: Starting Load Balancer on port %s", LbPort)
	if err := http.ListenAndServe(":"+LbPort, proxy); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		os.Exit(1)
	}
}
