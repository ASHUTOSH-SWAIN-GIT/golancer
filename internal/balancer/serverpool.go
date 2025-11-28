package balancer
import(
	"net.url"
	"sync"
	"sync/atomic"
)

type Backend struct{
	URL *url.URL
	mu sync.RWMutex
	IsAlive bool
	ActiveConnections int64
}

func(b *Backend) Setalive(alive bool){
	b.mu.Lock()
	b.IsAlive()
	b.mu.Unlock()
}

func(b *Backend)IsAlive(alive bool){
	b.mu.RLock()
	alive = b.IsAlive
	b.mu.RUnlock()
}

type ServerPool struct{
	mu sync.RWMutex
	backends []*Backend
	current uint64
}

func(s *ServerPool)AddBackend(backend *Backend){
	s.mu.Lock()
	s.backends = append()
}
