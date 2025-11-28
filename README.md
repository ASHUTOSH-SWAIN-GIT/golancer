# GoLancer

A high-performance HTTP load balancer written in Go, featuring health checks, round-robin load balancing, and active connection tracking.

## Features

- **Round-Robin Load Balancing**: Distributes incoming requests evenly across healthy backend servers
- **Health Checks**: Periodic health monitoring of backend servers with configurable intervals
- **Active Connection Tracking**: Monitors active connections per backend server
- **Automatic Failover**: Automatically routes traffic away from unhealthy servers
- **Reverse Proxy**: Built on Go's `httputil.ReverseProxy` for efficient request forwarding
- **Thread-Safe**: Uses mutexes and atomic operations for concurrent safety

## Architecture

The load balancer consists of three main components:

1. **Backend**: Represents a backend server with health status and connection tracking
2. **ServerPool**: Manages a collection of backends and implements round-robin selection
3. **Health Checker**: Periodically checks backend health via HTTP health endpoints

## Project Structure

```
golancer/
├── main.go                          # Main entry point and HTTP server setup
├── internal/
│   └── balancer/
│       ├── serverpool.go            # Server pool and load balancing logic
│       └── healthcheck.go           # Health check implementation
└── go.mod                           # Go module definition
```

## Prerequisites

- Go 1.25.4 or higher
- Backend servers with a `/health` endpoint that returns HTTP 200 when healthy

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd golancer
```

2. Build the project:
```bash
go build
```

3. Run the load balancer:
```bash
./golancer
```

Or run directly with Go:
```bash
go run main.go
```

## Configuration

Edit `main.go` to configure:

- **Load Balancer Port**: Change `LbPort` constant (default: `8080`)
- **Backend URLs**: Modify the `backendURLs` slice to include your backend servers
- **Health Check Interval**: Adjust `HealthCheckInterval` constant (default: `10 seconds`)

Example configuration:
```go
const LbPort = "8080"
const HealthCheckInterval = 10 * time.Second

var backendURLs = []string{
    "http://localhost:8081",
    "http://localhost:8082",
    "http://localhost:8083",
}
```

## Health Check Endpoint

Backend servers must implement a health check endpoint at `/health` that:
- Returns HTTP 200 status code when the server is healthy
- Responds within 5 seconds (health check timeout)

Example health endpoint implementation:
```go
http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
})
```

## Usage

1. Start your backend servers on ports 8081, 8082, and 8083 (or modify the configuration)
2. Ensure each backend has a `/health` endpoint
3. Start the load balancer:
```bash
go run main.go
```

4. Send requests to the load balancer:
```bash
curl http://localhost:8080/your-endpoint
```

The load balancer will:
- Route requests to healthy backends using round-robin
- Track active connections per backend
- Automatically skip unhealthy servers
- Return 502 Bad Gateway if no healthy servers are available

## How It Works

1. **Request Routing**: When a request arrives, the load balancer selects the next healthy backend using round-robin algorithm
2. **Connection Tracking**: Active connections are incremented when a request is forwarded and decremented when the response is received
3. **Health Monitoring**: A background goroutine periodically checks each backend's `/health` endpoint
4. **Failover**: Unhealthy backends are automatically excluded from the rotation until they recover

## Logging

The load balancer logs:
- Backend addition events
- Health check start/status
- Backend health status (up/down)
- Proxy errors
- Server startup information

## Error Handling

- **No Healthy Backends**: Returns HTTP 502 Bad Gateway with an error message
- **Invalid Backend URLs**: Fails to start with a fatal error
- **Proxy Errors**: Logs errors and returns 502 to the client

## Thread Safety

The implementation uses:
- `sync.RWMutex` for read/write locks on shared data structures
- `sync/atomic` for lock-free operations on connection counters and round-robin index
- Goroutines for concurrent health checks


