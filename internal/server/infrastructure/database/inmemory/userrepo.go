package inmemory

import (
	"context"
	"github.com/hotafrika/griz-backend/internal/server/domain"
	"github.com/hotafrika/griz-backend/internal/server/domain/entities"
	"github.com/pkg/errors"
	"sync"
)

// UserRepository is inmemory implementation
type UserRepository struct {
	users       map[uint64]entities.User
	usersByName map[string]uint64
	lastID      uint64
	rmu         sync.RWMutex
}

var _ domain.UserRepository = (*UserRepository)(nil)

// NewUserRepository creates new UserRepository
func NewUserRepository() *UserRepository {
	return &UserRepository{
		users:       make(map[uint64]entities.User),
		usersByName: make(map[string]uint64),
		lastID:      0,
	}
}

// Get returns User by ID
func (u *UserRepository) Get(ctx context.Context, u2 uint64) (entities.User, error) {
	u.rmu.RLock()
	v, ok := u.users[u2]
	u.rmu.RUnlock()
	if !ok {
		return entities.User{}, domain.ErrUserNotFound
	}

	// remove password
	v.Password = ""

	return v, nil
}

// Create adds new user to repo
func (u *UserRepository) Create(ctx context.Context, user entities.User) (uint64, error) {
	// check username
	u.rmu.RLock()
	_, ok := u.usersByName[user.Username]
	u.rmu.RUnlock()
	if ok {
		return 0, errors.New("username already exists")
	}

	u.rmu.Lock()
	newID := u.lastID + 1
	user.ID = newID
	u.users[newID] = user
	u.usersByName[user.Username] = newID
	u.lastID = newID
	u.rmu.Unlock()
	return newID, nil
}

// GetByUsernameAndPass returns User by Username and Pass
func (u *UserRepository) GetByUsernameAndPass(ctx context.Context, user entities.User) (uint64, error) {
	// check username
	u.rmu.RLock()
	v, ok := u.usersByName[user.Username]
	u.rmu.RUnlock()
	if !ok {
		return 0, domain.ErrUserNotFound
	}

	// check pass
	u.rmu.RLock()
	repUser, ok := u.users[v]
	u.rmu.RUnlock()
	if !ok {
		return 0, domain.ErrUserNotFound
	}
	if repUser.Password != user.Password {
		return 0, domain.ErrUserNotFound
	}

	return v, nil
}
