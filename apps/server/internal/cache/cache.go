package cache

import (
	"sync"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
)

// Config holds cache configuration
type Config struct {
	DefaultTTL      time.Duration
	MaxEntries      int
	CleanupInterval time.Duration
}

// Manager handles all caching operations
type Manager struct {
	mu     sync.RWMutex
	cache  *expirable.LRU[string, *Entry]
	config Config
}

// Entry represents a cached item
type Entry struct {
	Value     interface{}
	ExpiresAt time.Time
	Key       string
}

// New creates a new cache manager
func New(cfg Config) *Manager {
	return &Manager{
		cache:  expirable.NewLRU[string, *Entry](cfg.MaxEntries, nil, cfg.CleanupInterval),
		config: cfg,
	}
}

// Get retrieves a value from cache
func (m *Manager) Get(key string) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, ok := m.cache.Get(key)
	if !ok {
		return nil, false
	}

	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}

	return entry.Value, true
}

// Set stores a value in cache with default TTL
func (m *Manager) Set(key string, value interface{}) {
	m.SetWithTTL(key, value, m.config.DefaultTTL)
}

// SetWithTTL stores a value with a specific TTL
func (m *Manager) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cache.Add(key, &Entry{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
		Key:       key,
	})
}

// Delete removes a value from cache
func (m *Manager) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cache.Remove(key)
}

// Clear removes all entries matching a pattern
func (m *Manager) Clear(pattern string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Simple prefix matching
	keys := m.cache.Keys()
	for _, key := range keys {
		if pattern == "*" || matchPattern(pattern, key) {
			m.cache.Remove(key)
		}
	}
}

// matchPattern checks if a key matches a simple pattern (*, : wildcards)
func matchPattern(pattern, key string) bool {
	if pattern == "*" {
		return true
	}
	// Handle prefix:* pattern
	if len(pattern) > 0 && pattern[len(pattern)-1] == '*' {
		prefix := pattern[:len(pattern)-1]
		return len(key) >= len(prefix) && key[:len(prefix)] == prefix
	}
	return pattern == key
}

// GetOrCompute gets a cached value or computes and caches it
func (m *Manager) GetOrCompute(key string, ttl time.Duration, compute func() (interface{}, error)) (interface{}, error) {
	// Try cache first
	if value, ok := m.Get(key); ok {
		return value, nil
	}

	// Compute
	value, err := compute()
	if err != nil {
		return nil, err
	}

	// Cache result
	m.SetWithTTL(key, value, ttl)
	return value, nil
}

// Stats returns cache statistics
func (m *Manager) Stats() (keys int, maxEntries int) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.cache.Len(), m.config.MaxEntries
}
