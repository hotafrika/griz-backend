package domain

import (
	"context"
	"github.com/hotafrika/griz-backend/internal/server/domain/entities"
	"github.com/pkg/errors"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	// Get (ctx, UserID) -> (User, error)
	Get(context.Context, uint64) (entities.User, error)
	// Create (ctx, User) -> (UserID, error)
	Create(context.Context, entities.User) (uint64, error)
	// GetByUsernameAndPass (ctx, User) -> (UserID, error)
	GetByUsernameAndPass(context.Context, entities.User) (uint64, error)
}

var ErrCodeNotFound = errors.New("code not found")

type CodeRepository interface {
	// List (ctx, UserID, offset, limit) -> ([]Code, error)
	List(context.Context, uint64, int64, int64) ([]entities.Code, error)
	// ListAll (ctx, UserID) -> ([]Code, error)
	ListAll(context.Context, uint64) ([]entities.Code, error)
	// Get (ctx, CodeID) -> (Code, error)
	Get(context.Context, uint64) (entities.Code, error)
	// Create (ctx, Code) -> (CodeID, error)
	Create(context.Context, entities.Code) (uint64, error)
	// Update (ctx, Code) -> (error)
	Update(context.Context, entities.Code) error
	// Delete (ctx, CodeID) -> (error)
	Delete(context.Context, uint64) error
}
