package inmemory

import (
	"context"
	"github.com/hotafrika/griz-backend/internal/server/domain/entities"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestCodeRepository_Create(t *testing.T) {
	tests := []struct {
		code       entities.Code
		wantLen    int
		wantLastID uint64
	}{
		{
			code: entities.Code{
				UserID: 1,
				SrcURL: "1",
				Hash:   "abc",
			},
			wantLen:    1,
			wantLastID: 1,
		},
		{
			code: entities.Code{
				UserID: 2,
				SrcURL: "2",
				Hash:   "cde",
			},
			wantLen:    2,
			wantLastID: 2,
		},
		{
			code: entities.Code{
				UserID: 1,
				SrcURL: "3",
				Hash:   "def",
			},
			wantLen:    3,
			wantLastID: 3,
		},
	}

	cr := NewCodeRepository()
	ctx := context.TODO()

	for _, tt := range tests {
		t.Run(tt.code.SrcURL, func(t *testing.T) {
			id, err := cr.Create(ctx, tt.code)
			assert.Equal(t, tt.wantLen, len(cr.codes))
			if assert.NoError(t, err) {
				assert.Equal(t, tt.wantLastID, id)
			}
		})
	}
}

func TestCodeRepository_ListAll(t *testing.T) {
	tests := []struct {
		code       entities.Code
		wantLen    int
		wantLastID uint64
	}{
		{
			code: entities.Code{
				UserID: 1,
				SrcURL: "1",
				Hash:   "abc",
			},
			wantLen: 1,
		},
		{
			code: entities.Code{
				UserID: 2,
				SrcURL: "2",
				Hash:   "cde",
			},
			wantLen: 1,
		},
		{
			code: entities.Code{
				UserID: 1,
				SrcURL: "3",
				Hash:   "def",
			},
			wantLen: 2,
		},
	}

	cr := NewCodeRepository()
	ctx := context.TODO()

	for _, tt := range tests {
		t.Run(tt.code.SrcURL, func(t *testing.T) {
			_, err := cr.Create(ctx, tt.code)
			if assert.NoError(t, err) {
				list, err := cr.ListAll(ctx, tt.code.UserID)
				if assert.NoError(t, err) {
					assert.Equal(t, tt.wantLen, len(list))
				}
			}
		})
	}
}

func TestCodeRepository_Get(t *testing.T) {
	tests := []struct {
		code       entities.Code
		wantLen    int
		wantLastID uint64
	}{
		{
			code: entities.Code{
				UserID: 1,
				SrcURL: "1",
				Hash:   "abc",
			},
		},
		{
			code: entities.Code{
				UserID: 2,
				SrcURL: "2",
				Hash:   "cde",
			},
		},
		{
			code: entities.Code{
				UserID: 1,
				SrcURL: "3",
				Hash:   "def",
			},
		},
	}

	cr := NewCodeRepository()
	ctx := context.TODO()
	ids := make([]uint64, 0, len(tests))

	for _, tt := range tests {
		id, err := cr.Create(ctx, tt.code)
		if assert.NoError(t, err) {
			ids = append(ids, id)
		}
	}

	for _, id := range ids {
		t.Run(strconv.FormatInt(int64(id), 10), func(t *testing.T) {
			_, err := cr.Get(ctx, id)
			assert.NoError(t, err)
		})
	}
}

func TestCodeRepository_GetByToken(t *testing.T) {
	tests := []struct {
		code       entities.Code
		wantLen    int
		wantLastID uint64
	}{
		{
			code: entities.Code{
				UserID: 1,
				SrcURL: "1",
				Hash:   "abc",
			},
		},
		{
			code: entities.Code{
				UserID: 2,
				SrcURL: "2",
				Hash:   "cde",
			},
		},
		{
			code: entities.Code{
				UserID: 1,
				SrcURL: "3",
				Hash:   "def",
			},
		},
	}

	cr := NewCodeRepository()
	ctx := context.TODO()
	tokens := make([]string, 0, len(tests))

	for _, tt := range tests {
		_, err := cr.Create(ctx, tt.code)
		if assert.NoError(t, err) {
			tokens = append(tokens, tt.code.Hash)
		}
	}

	for _, token := range tokens {
		t.Run(token, func(t *testing.T) {
			_, err := cr.GetByHash(ctx, token)
			assert.NoError(t, err)
		})
	}
}

func TestCodeRepository_Update(t *testing.T) {
	tests := []struct {
		code       entities.Code
		wantLen    int
		wantLastID uint64
	}{
		{
			code: entities.Code{
				UserID: 1,
				SrcURL: "1",
				Hash:   "abc",
			},
		},
		{
			code: entities.Code{
				UserID: 2,
				SrcURL: "2",
				Hash:   "cde",
			},
		},
		{
			code: entities.Code{
				UserID: 1,
				SrcURL: "3",
				Hash:   "def",
			},
		},
	}

	cr := NewCodeRepository()
	ctx := context.TODO()
	id, err := cr.Create(ctx, entities.Code{
		UserID: 1,
		SrcURL: "4",
		Hash:   "efg",
	})
	assert.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.code.SrcURL, func(t *testing.T) {
			tt.code.ID = id
			err = cr.Update(ctx, tt.code)
			if assert.NoError(t, err) {
				val, err := cr.Get(ctx, id)
				if assert.NoError(t, err) {
					assert.Equal(t, tt.code, val)
				}
			}
		})
	}
}

func TestCodeRepository_Delete(t *testing.T) {
	tests := []struct {
		code       entities.Code
		wantLen    int
		wantLastID uint64
	}{
		{
			code: entities.Code{
				ID:     1,
				UserID: 1,
				SrcURL: "1",
				Hash:   "abc",
			},
			wantLen:    2,
			wantLastID: 3,
		},
		{
			code: entities.Code{
				ID:     2,
				UserID: 2,
				SrcURL: "2",
				Hash:   "cde",
			},
			wantLen:    1,
			wantLastID: 3,
		},
		{
			code: entities.Code{
				ID:     3,
				UserID: 1,
				SrcURL: "3",
				Hash:   "def",
			},
			wantLen:    0,
			wantLastID: 3,
		},
	}

	cr := NewCodeRepository()
	ctx := context.TODO()

	for _, tt := range tests {
		_, err := cr.Create(ctx, tt.code)
		assert.NoError(t, err)
	}

	for _, tt := range tests {
		t.Run(tt.code.SrcURL, func(t *testing.T) {
			err := cr.Delete(ctx, tt.code.ID)
			if assert.NoError(t, err) {
				assert.Equal(t, tt.wantLen, len(cr.codes))
				assert.Equal(t, tt.wantLastID, cr.lastID)
			}
		})
	}
}
