package routeros

import (
	"fmt"
	"sync"
	"time"

	ros "github.com/go-routeros/routeros/v3"
)

// Pool manages one RouterOS client per named router session.
type Pool struct {
	mu      sync.Mutex
	clients map[string]*entry
}

type entry struct {
	client   *ros.Client
	lastUsed time.Time
}

// NewPool creates a new connection pool.
func NewPool() *Pool {
	p := &Pool{clients: make(map[string]*entry)}
	go p.cleanupLoop()
	return p
}

// Get returns (or lazily creates) a RouterOS client for the given session.
// On connection failure it attempts a single reconnect.
func (p *Pool) Get(sessionName, host, user, password string) (*ros.Client, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	e, ok := p.clients[sessionName]
	if ok {
		// Test existing connection with a lightweight command.
		if _, err := e.client.RunArgs([]string{"/system/identity/print"}); err == nil {
			e.lastUsed = time.Now()
			return e.client, nil
		}
		// Connection dead — close and reconnect.
		e.client.Close()
		delete(p.clients, sessionName)
	}

	c, err := ros.Dial(host, user, password)
	if err != nil {
		return nil, fmt.Errorf("pool connect %s: %w", sessionName, err)
	}
	p.clients[sessionName] = &entry{client: c, lastUsed: time.Now()}
	return c, nil
}

// Close removes and closes a specific session's connection.
func (p *Pool) Close(sessionName string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if e, ok := p.clients[sessionName]; ok {
		e.client.Close()
		delete(p.clients, sessionName)
	}
}

// CloseAll closes every open connection.
func (p *Pool) CloseAll() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for name, e := range p.clients {
		e.client.Close()
		delete(p.clients, name)
	}
}

// cleanupLoop periodically closes connections idle for more than 10 minutes.
func (p *Pool) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		p.mu.Lock()
		cutoff := time.Now().Add(-10 * time.Minute)
		for name, e := range p.clients {
			if e.lastUsed.Before(cutoff) {
				e.client.Close()
				delete(p.clients, name)
			}
		}
		p.mu.Unlock()
	}
}
