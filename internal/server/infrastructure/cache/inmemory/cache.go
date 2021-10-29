package inmemory

import (
	"fmt"
	"github.com/hotafrika/griz-backend/internal/server/domain"
	"sync"
	"time"
)

var ErrNotExist = fmt.Errorf("key not exist")

// Cache implements inmemory cache type
type Cache struct {
	data map[string]string
	rw   sync.RWMutex
}

var _ domain.Cacher = (*Cache)(nil)

// NewCache creates Cache
func NewCache() *Cache {
	return &Cache{
		data: make(map[string]string),
	}
}

// Get receives value by key
func (c *Cache) Get(key fmt.Stringer) (string, error) {
	c.rw.RLock()
	value, ok := c.data[key.String()]
	c.rw.RUnlock()
	if !ok {
		return "", ErrNotExist
	}
	return value, nil
}

// Set sets key - value in memory
func (c *Cache) Set(key fmt.Stringer, value string, ttl time.Duration) error {
	c.rw.Lock()
	c.data[key.String()] = value
	c.rw.Unlock()
	return nil
}

// Delete removes value by key
func (c *Cache) Delete(key fmt.Stringer) error {
	c.rw.Lock()
	delete(c.data, key.String())
	c.rw.Unlock()
	return nil
}
