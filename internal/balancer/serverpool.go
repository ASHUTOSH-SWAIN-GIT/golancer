package balancer

import (
	"net/url"
	"sync"
	"sync/atomic"
)

type Backend struct {
	URL               *url.URL
	mu                sync.RWMutex
	isAlive           bool
	ActiveConnections int64
}

func (b *Backend) Setalive(alive bool) {
	b.mu.Lock()
	b.isAlive = alive
	b.mu.Unlock()
}

func (b *Backend) IsAlive() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.isAlive
}

type ServerPool struct {
	mu       sync.RWMutex
	backends []*Backend
	current  uint64
}

func (s *ServerPool) AddBackend(backend *Backend) {
	s.mu.Lock()
	s.backends = append(s.backends, backend)
	s.mu.Unlock()
}

func (s *ServerPool) GetBackends() []*Backend {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.backends
}

func (s *ServerPool) NextIndex() int {
	return int(atomic.AddUint64(&s.current, 1) % uint64(len(s.backends)))
}

func (s *ServerPool) GetNextHealthyServer() *Backend {
	s.mu.RLock()
	defer s.mu.Unlock()

	TotalServers := len(s.backends)
	if TotalServers == 0 {
		return nil
	}

	startIdx := s.NextIndex()

	for i := 0; i < TotalServers; i++ {
		idx := (startIdx + i) % TotalServers
		backend := s.backends[idx]
		if backend.IsAlive() {
			atomic.StoreUint64(&s.current, uint64(idx))
			return backend
		}
	}

	return nil
}

func (b *Backend) IncrementActiveConnections() {
	atomic.AddInt64(&b.ActiveConnections, 1)
}

func (b *Backend) DecrementActiveConnections() {
	atomic.AddInt64(&b.ActiveConnections, -1)
}
