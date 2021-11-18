package authtoken

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestJWT_MakeByID(t *testing.T) {
	tests := []struct {
		id uint64
	}{
		{
			id: 0,
		},
		{
			id: 1,
		},
		{
			id: 100000000000,
		},
	}
	j := NewJWTFromString("abc", time.Minute)
	for _, tt := range tests {
		t.Run(strconv.FormatInt(int64(tt.id), 10), func(t *testing.T) {
			_, err := j.MakeByID(tt.id)
			assert.NoError(t, err)
		})
	}
}
