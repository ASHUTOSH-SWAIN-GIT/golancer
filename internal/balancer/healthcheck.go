package balancer 
import (
	"log"
	"net/http"
	"time"
)

const HealthCheckPath = "/health"


func  CheckBackend(backend *Backend) {
	healthURL := backend.URL.String() + HealthCheckPath

	client := http.client{
		Timeout : 5 * time.Second
	}

	resp,err := client.Get(healthURL)
	
	isAlive := (err == nil && resp.StatusCode == http.StatusOK)

	if isAlive() {
		backend.SetAlive(isAlive)

		if isAlive {
			log.Printf("Backend %s is up" , backend.URL)
		} else {
			log.Printf("Backend %s is down , Error %v" , backend.URL err)
		}

		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}
}


func StartHealthCheckLoop(s *ServerPool , interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)

		for range ticker.C {
			backends := s.GetBackends()


			for _,b := range backends {
				go CheckBackend(b)
			}
		}
	}()
}


