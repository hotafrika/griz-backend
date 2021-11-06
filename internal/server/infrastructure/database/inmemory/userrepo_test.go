package inmemory

import (
	"context"
	"github.com/hotafrika/griz-backend/internal/server/domain/entities"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	tests := []struct {
		user       entities.User
		wantLen1   int
		wantLen2   int
		wantLastID uint64
	}{
		{
			user: entities.User{
				Username: "1",
				Password: "1",
				Email:    "1",
			},
			wantLen1:   1,
			wantLen2:   1,
			wantLastID: 1,
		},
		{
			user: entities.User{
				Username: "2",
				Password: "2",
				Email:    "2",
			},
			wantLen1:   2,
			wantLen2:   2,
			wantLastID: 2,
		},
		{
			user: entities.User{
				Username: "3",
				Password: "3",
				Email:    "3",
			},
			wantLen1:   3,
			wantLen2:   3,
			wantLastID: 3,
		},
	}

	ur := NewUserRepository()
	ctx := context.TODO()
	for _, tt := range tests {
		t.Run(tt.user.Username, func(t *testing.T) {
			ur.Create(ctx, tt.user)
			assert.Equal(t, tt.wantLen1, len(ur.users))
			assert.Equal(t, tt.wantLen2, len(ur.usersByName))
			assert.Equal(t, tt.wantLastID, ur.lastID)
		})
	}
}

func TestUserRepository_Get(t *testing.T) {
	tests := []struct {
		user       entities.User
		wantLen1   int
		wantLen2   int
		wantLastID uint64
	}{
		{
			user: entities.User{
				Username: "1",
				Password: "1",
				Email:    "1",
			},
			wantLen1:   1,
			wantLen2:   1,
			wantLastID: 1,
		},
		{
			user: entities.User{
				Username: "2",
				Password: "2",
				Email:    "2",
			},
			wantLen1:   2,
			wantLen2:   2,
			wantLastID: 2,
		},
		{
			user: entities.User{
				Username: "3",
				Password: "3",
				Email:    "3",
			},
			wantLen1:   3,
			wantLen2:   3,
			wantLastID: 3,
		},
	}

	ur := NewUserRepository()
	ids := make([]uint64, 0, len(tests))
	ctx := context.TODO()
	for _, tt := range tests {
		id, err := ur.Create(ctx, tt.user)
		if assert.NoError(t, err) {
			ids = append(ids, id)
		}
	}

	for _, id := range ids {
		t.Run(strconv.FormatInt(int64(id), 10), func(t *testing.T) {
			_, err := ur.Get(ctx, id)
			assert.NoError(t, err)
		})
	}
}

func TestUserRepository_GetByUsernameAndPass(t *testing.T) {
	tests := []struct {
		user       entities.User
		wantLen1   int
		wantLen2   int
		wantLastID uint64
	}{
		{
			user: entities.User{
				Username: "1",
				Password: "1",
				Email:    "1",
			},
			wantLen1:   1,
			wantLen2:   1,
			wantLastID: 1,
		},
		{
			user: entities.User{
				Username: "2",
				Password: "2",
				Email:    "2",
			},
			wantLen1:   2,
			wantLen2:   2,
			wantLastID: 2,
		},
		{
			user: entities.User{
				Username: "3",
				Password: "3",
				Email:    "3",
			},
			wantLen1:   3,
			wantLen2:   3,
			wantLastID: 3,
		},
	}

	ur := NewUserRepository()
	ctx := context.TODO()
	for _, tt := range tests {
		_, err := ur.Create(ctx, tt.user)
		assert.NoError(t, err)
	}

	for _, tt := range tests {
		t.Run(tt.user.Username, func(t *testing.T) {
			_, err := ur.GetByUsernameAndPass(ctx, tt.user)
			assert.NoError(t, err)
		})
	}
}
