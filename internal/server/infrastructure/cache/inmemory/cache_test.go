package inmemory

import (
	"fmt"
	"github.com/hotafrika/griz-backend/internal/server/infrastructure/cache"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCache_Set(t *testing.T) {
	tests := []struct {
		name      string
		key       fmt.Stringer
		value     string
		mapLength int
	}{
		{
			name:      "first value",
			key:       cache.AuthToken{Key: "abc"},
			value:     "1",
			mapLength: 1,
		},
		{
			name:      "second value",
			key:       cache.SocialUrl{Key: "http://second.com"},
			value:     "2",
			mapLength: 2,
		},
		{
			name:      "third value",
			key:       cache.HashUrl{Key: "http://third.com"},
			value:     "3",
			mapLength: 3,
		},
	}

	c := NewCache()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.Set(tt.key, tt.value, 0)
			assert.NoError(t, err)
			assert.Equal(t, tt.mapLength, len(c.data), "len of underlying map")
		})
	}
}

func TestCache_Get(t *testing.T) {
	tests := []struct {
		name  string
		key   fmt.Stringer
		value string
	}{
		{
			name:  "first value",
			key:   cache.AuthToken{Key: "abc"},
			value: "1",
		},
		{
			name:  "second value",
			key:   cache.SocialUrl{Key: "http://second.com"},
			value: "2",
		},
		{
			name:  "third value",
			key:   cache.HashUrl{Key: "http://third.com"},
			value: "3",
		},
	}

	c := NewCache()
	for _, tt := range tests {
		err := c.Set(tt.key, tt.value, 0)
		assert.NoError(t, err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := c.Get(tt.key)
			if assert.NoError(t, err) {
				assert.Equal(t, tt.value, s)
			}
		})
	}
}

func TestCache_Delete(t *testing.T) {
	tests := []struct {
		name  string
		key   fmt.Stringer
		value string
	}{
		{
			name:  "first value",
			key:   cache.AuthToken{Key: "abc"},
			value: "1",
		},
		{
			name:  "second value",
			key:   cache.SocialUrl{Key: "http://second.com"},
			value: "2",
		},
		{
			name:  "third value",
			key:   cache.HashUrl{Key: "http://third.com"},
			value: "3",
		},
	}

	c := NewCache()
	for _, tt := range tests {
		err := c.Set(tt.key, tt.value, 0)
		assert.NoError(t, err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.Delete(tt.key)
			assert.NoError(t, err)
			_, err = c.Get(tt.key)
			assert.Error(t, err)
		})
	}
}
