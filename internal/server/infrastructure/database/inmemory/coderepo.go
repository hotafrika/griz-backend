package inmemory

import (
	"context"
	"github.com/hotafrika/griz-backend/internal/server/domain"
	"github.com/hotafrika/griz-backend/internal/server/domain/entities"
	"sync"
)

// CodeRepository represents inmemory repo
type CodeRepository struct {
	lastID    uint64
	codes     map[uint64]entities.Code
	userCodes map[uint64]uint64
	rmu       sync.RWMutex
}

var _ domain.CodeRepository = (*CodeRepository)(nil)

// NewCodeRepository creates new CodeRepository
func NewCodeRepository() *CodeRepository {
	return &CodeRepository{
		lastID:    0,
		codes:     make(map[uint64]entities.Code),
		userCodes: nil,
		rmu:       sync.RWMutex{},
	}
}

// List returns codes by userID with offset and limit
func (c *CodeRepository) List(ctx context.Context, u uint64, i int64, i2 int64) ([]entities.Code, error) {
	return c.ListAll(ctx, u)
}

// ListAll returns codes by userID
func (c *CodeRepository) ListAll(ctx context.Context, u uint64) ([]entities.Code, error) {
	var codes []entities.Code
	c.rmu.RLock()
	for _, code := range c.codes {
		if u == code.UserID {
			codes = append(codes, code)
		}
	}
	c.rmu.RUnlock()
	if len(codes) == 0 {
		return nil, domain.ErrCodeNotFound
	}
	return codes, nil
}

// Get returns code by its ID
func (c *CodeRepository) Get(ctx context.Context, u uint64) (entities.Code, error) {
	c.rmu.RLock()
	v, ok := c.codes[u]
	c.rmu.RUnlock()
	if !ok {
		return entities.Code{}, domain.ErrCodeNotFound
	}
	return v, nil
}

// GetByHash returns code by its token
func (c *CodeRepository) GetByHash(ctx context.Context, token string) (entities.Code, error) {
	result := entities.Code{}
	ok := false
	c.rmu.RLock()
	for _, v := range c.codes {
		if v.Hash == token {
			result = v
			ok = true
			break
		}
	}
	c.rmu.RUnlock()
	if !ok {
		return entities.Code{}, domain.ErrCodeNotFound
	}
	return result, nil
}

// Create adds new code to repo
func (c *CodeRepository) Create(ctx context.Context, code entities.Code) (uint64, error) {
	c.rmu.Lock()
	newID := c.lastID + 1
	c.codes[newID] = code
	c.lastID = newID
	c.rmu.Unlock()
	return newID, nil
}

// Update updates existing code in repo
func (c *CodeRepository) Update(ctx context.Context, code entities.Code) error {
	c.rmu.Lock()
	defer c.rmu.Unlock()
	_, ok := c.codes[code.ID]
	if !ok {
		return domain.ErrCodeNotFound
	}
	c.codes[code.ID] = code
	return nil
}

// Delete removes code from repo
func (c *CodeRepository) Delete(ctx context.Context, u uint64) error {
	c.rmu.Lock()
	defer c.rmu.Unlock()
	_, ok := c.codes[u]
	if !ok {
		return domain.ErrCodeNotFound
	}
	delete(c.codes, u)
	return nil
}
