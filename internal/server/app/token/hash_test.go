package token

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestAES_Create(t *testing.T) {
	var key = "1234567890123456"
	tests := []struct {
		name      string
		id        uint64
		wantErr   bool
		wantToken string
	}{
		{
			name:      "zero",
			id:        0,
			wantErr:   true,
			wantToken: "v011771075a74696a04c3e718880cf09112",
		},
		{
			name:      "one",
			id:        1,
			wantErr:   false,
			wantToken: "v01bf00641e37ce8da7d29cd4eab69b96ea",
		},
		{
			name:      "one hundred",
			id:        100,
			wantErr:   false,
			wantToken: "v015064e7064888d52339da98116111e07a",
		},
		{
			name:      "one million",
			id:        1000000,
			wantErr:   false,
			wantToken: "v0129e576ad8c4e2a9de4a99ebecd032c6c",
		},
		{
			name:      "max",
			id:        math.MaxUint64,
			wantErr:   false,
			wantToken: "v01e6a7dd3e202395f08573bf2c863792d4",
		},
	}
	a, err := NewAES(key)
	if assert.NoError(t, err) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				fmt.Println(tt.id)
				token, err := a.Create(tt.id)
				if tt.wantErr {
					assert.Error(t, err)
					return
				}
				fmt.Println("token: ", token)
				if assert.NoError(t, err) {
					assert.Equal(t, tt.wantToken, token)
				}
			})
		}
	}
}

func TestAES_Decode(t *testing.T) {
	var key = "1234567890123456"
	tests := []struct {
		name    string
		wantId  uint64
		wantErr bool
		token   string
	}{
		{
			name:    "zero",
			wantId:  0,
			wantErr: true,
			token:   "v011771075a74696a04c3e718880cf09112",
		},
		{
			name:    "one",
			wantId:  1,
			wantErr: false,
			token:   "v01bf00641e37ce8da7d29cd4eab69b96ea",
		},
		{
			name:    "one hundred",
			wantId:  100,
			wantErr: false,
			token:   "v015064e7064888d52339da98116111e07a",
		},
		{
			name:    "one million",
			wantId:  1000000,
			wantErr: false,
			token:   "v0129e576ad8c4e2a9de4a99ebecd032c6c",
		},
		{
			name:    "max",
			wantId:  math.MaxUint64,
			wantErr: false,
			token:   "v01e6a7dd3e202395f08573bf2c863792d4",
		},
	}
	a, err := NewAES(key)
	if assert.NoError(t, err) {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				id, err := a.Decode(tt.token)
				if tt.wantErr {
					assert.Error(t, err)
					return
				}
				if assert.NoError(t, err) {
					assert.Equal(t, tt.wantId, id)
				}
			})
		}
	}
}
