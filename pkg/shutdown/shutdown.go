package shutdown

import (
	"context"
	"log"
	"sync"
)

// ShutdownFunc is a function that will be called during shutdown
type ShutdownFunc func(ctx context.Context) error

// Manager handles graceful shutdown of the application
type Manager struct {
	handlers map[string]ShutdownFunc
	mu       sync.RWMutex
}

var (
	manager *Manager
	once    sync.Once
)

// GetManager returns the singleton instance of the shutdown manager
func GetManager() *Manager {
	once.Do(func() {
		manager = &Manager{
			handlers: make(map[string]ShutdownFunc),
		}
	})
	return manager
}

// Register adds a new shutdown handler with the given name
func (m *Manager) Register(name string, handler ShutdownFunc) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers[name] = handler
}

// Unregister removes a shutdown handler by name
func (m *Manager) Unregister(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.handlers, name)
}

// Shutdown triggers the shutdown process with the given context
func (m *Manager) Shutdown(ctx context.Context) {
	log.Println("Shutting down gracefully...")

	// Create a channel to signal completion
	done := make(chan interface{})

	go func() {
		m.mu.RLock()
		defer m.mu.RUnlock()

		var wg sync.WaitGroup
		for name, handler := range m.handlers {
			wg.Add(1)
			go func(name string, handler ShutdownFunc) {
				defer wg.Done()
				if err := handler(ctx); err != nil {
					log.Printf("Error shutting down %s: %v", name, err)
				} else {
					log.Printf("Successfully shut down %s", name)
				}
			}(name, handler)
		}
		wg.Wait()
		done <- nil // Signal completion by sending nil
	}()

	select {
	case <-ctx.Done():
		log.Println("Shutdown timeout reached, forcing exit")
	case <-done:
		log.Println("All handlers completed successfully")
	}
}
