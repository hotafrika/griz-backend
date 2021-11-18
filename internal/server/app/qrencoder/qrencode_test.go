package qrencoder

import (
	qrdecoder2 "github.com/hotafrika/griz-backend/internal/server/app/qrdecoder"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestYeqown_Encode(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		blockSize uint8
		wantErr   bool
	}{
		{
			name:      "valid long block 6",
			data:      "https://someapp.somedomain.com/apps?q=123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
			wantErr:   false,
			blockSize: 6,
		},
		{
			name:      "valid small block 6",
			data:      "https://someapp.somedomain.com/apps?q=123",
			wantErr:   false,
			blockSize: 6,
		},
		{
			name:      "valid long block 12",
			data:      "https://someapp.somedomain.com/apps?q=123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
			wantErr:   false,
			blockSize: 12,
		},
		{
			name:      "valid small block 12",
			data:      "https://someapp.somedomain.com/apps?q=123",
			wantErr:   false,
			blockSize: 12,
		},
	}
	m := qrdecoder2.Makiuchi{}
	for _, tt := range tests {
		y := NewYeqown(WithQRWidth(tt.blockSize), WithFileImagePNG("GrizLogo.png"))
		t.Run(tt.name, func(t *testing.T) {
			res, err := y.Encode([]byte(tt.data))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NotNil(t, res)

			res2, err := m.Decode(res)
			assert.NoError(t, err)
			assert.Equal(t, tt.data, string(res2))
		})
	}
}
