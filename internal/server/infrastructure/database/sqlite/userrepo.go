package sqlite

import (
	"context"
	"database/sql"
	"github.com/hotafrika/griz-backend/internal/server/domain"
	"github.com/hotafrika/griz-backend/internal/server/domain/entities"
	"github.com/pkg/errors"
)

// UserRepository is SQL implementation
type UserRepository struct {
	db *sql.DB
	//insertStmt *sql.Stmt
}

var _ domain.UserRepository = (*UserRepository)(nil)

// NewUserRepository creates new repository
func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{
		db: db,
	}
}

// Get returns user by userID
func (u UserRepository) Get(ctx context.Context, id uint64) (entities.User, error) {
	var user entities.User
	var username string
	var email string
	err := u.db.QueryRowContext(ctx, `SELECT username, email from users WHERE id=?`, id).
		Scan(&username, &email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, domain.ErrUserNotFound
		}
		return user, err
	}
	return entities.User{
		ID:       id,
		Username: username,
		Email:    email,
	}, nil
}

// Create creates new user
func (u UserRepository) Create(ctx context.Context, user entities.User) (uint64, error) {
	result, err := u.db.ExecContext(ctx, `INSERT INTO users(username, password, email) VALUES (?,?,?)`, user.Username, user.Password, user.Email)
	if err != nil {
		return 0, err
	}
	// TODO maybe replace with getting user by username and pass
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}

// GetByUsernameAndPass returns userID by username and pass
func (u UserRepository) GetByUsernameAndPass(ctx context.Context, user entities.User) (uint64, error) {
	var id uint64
	var password string
	err := u.db.QueryRowContext(ctx, `SELECT id, password from users WHERE username=?`, user.Username).
		Scan(&id, &password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, domain.ErrUserNotFound
		}
		return 0, err
	}
	if password != user.Password {
		return 0, domain.ErrUserNotFound
	}
	return id, nil
}
