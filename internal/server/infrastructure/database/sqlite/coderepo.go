package sqlite

import (
	"context"
	"database/sql"
	"github.com/hotafrika/griz-backend/internal/server/domain"
	"github.com/hotafrika/griz-backend/internal/server/domain/entities"
	"github.com/pkg/errors"
)

// CodeRepository is SQL implementation
type CodeRepository struct {
	db *sql.DB
}

var _ domain.CodeRepository = (*CodeRepository)(nil)

// NewCodeRepository creates new CodeRepository
func NewCodeRepository(db *sql.DB) CodeRepository {
	return CodeRepository{
		db: db,
	}
}

// List returns codes with pagination
func (c CodeRepository) List(ctx context.Context, userID uint64, offset int64, limit int64) ([]entities.Code, error) {
	return c.ListAll(ctx, userID)
}

// ListAll returns all codes
func (c CodeRepository) ListAll(ctx context.Context, userID uint64) ([]entities.Code, error) {
	rows, err := c.db.QueryContext(ctx, `SELECT id, link, hash FROM codes WHERE user_id=?`, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCodeNotFound
		}
		return nil, err
	}
	var id uint64
	var link string
	var hash string
	codes := make([]entities.Code, 0)
	for rows.Next() {
		err = rows.Scan(&id, &link, &hash)
		if err != nil {
			return nil, err
		}
		codes = append(codes, entities.Code{ID: id, SrcURL: link, Hash: hash, UserID: userID})
	}
	return codes, nil
}

// Get returns code by id
func (c CodeRepository) Get(ctx context.Context, id uint64) (entities.Code, error) {
	var code entities.Code
	var link string
	var hash string
	var userID uint64
	err := c.db.QueryRowContext(ctx, `SELECT link, hash, user_id from codes WHERE id=?`, id).
		Scan(&link, &hash, &userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return code, domain.ErrCodeNotFound
		}
		return code, err
	}
	return entities.Code{
		ID:     id,
		SrcURL: link,
		Hash:   hash,
		UserID: userID,
	}, nil
}

// GetByHash returns code by hash
func (c CodeRepository) GetByHash(ctx context.Context, hash string) (entities.Code, error) {
	var code entities.Code
	var id uint64
	var userID uint64
	var link string
	err := c.db.QueryRowContext(ctx, `SELECT id, link, user_id from codes WHERE hash=?`, hash).
		Scan(&id, &link, &userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return code, domain.ErrCodeNotFound
		}
		return code, err
	}
	return entities.Code{ID: id, SrcURL: link, Hash: hash, UserID: userID}, nil
}

// Create creates new code
func (c CodeRepository) Create(ctx context.Context, code entities.Code) (uint64, error) {
	result, err := c.db.ExecContext(ctx, `INSERT INTO codes(link, user_id) VALUES (?, ?)`, code.SrcURL, code.UserID)
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

// Update updates existing code
func (c CodeRepository) Update(ctx context.Context, code entities.Code) error {
	result, err := c.db.ExecContext(ctx,
		`UPDATE codes SET link=?, hash=?, user_id=? WHERE id=?`,
		code.SrcURL,
		code.Hash,
		code.UserID,
		code.ID)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return domain.ErrCodeNotFound
	}
	return nil
}

// Delete removes code
func (c CodeRepository) Delete(ctx context.Context, id uint64) error {
	result, err := c.db.ExecContext(ctx,
		`DELETE FROM codes WHERE id=?`,
		id)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return domain.ErrCodeNotFound
	}
	return nil
}
