package qrdecoder

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestMakiuchi_Decode(t *testing.T) {
	tests := []struct {
		filename string
		wantRes  string
		wantErr  bool
	}{
		{
			filename: "img1.png", // Look at QR format
			//wantRes: "http://q-r.to/QRCODE",
			wantRes: "",
			wantErr: true,
		},
		{
			filename: "img2.png",
			wantRes:  "https://urlgeni.us/instagram/coca-cola",
			wantErr:  false,
		},
		{
			filename: "img3.png",
			wantRes:  "https://urlgeni.us/amazon/nyc-2022?pqr",
			wantErr:  false,
		},
		{
			filename: "img4.png",
			wantRes:  "https://urlgeni.us/instagram/coca-cola?qr",
			wantErr:  false,
		},
		{
			filename: "img5.png",
			wantRes:  "http://onelink.to/d4vtgw",
			wantErr:  false,
		},
		{
			filename: "img6.jpg",
			wantRes:  "http://www.qrstuff.com",
			wantErr:  false,
		},
		{
			filename: "img7.bmp",
			wantRes:  "",
			wantErr:  true,
		},
		{
			filename: "ronaldo.png",
			wantRes:  "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		var m Makiuchi
		t.Run(tt.filename, func(t *testing.T) {
			file, err := os.Open(path.Join("testdata", tt.filename))
			assert.NoError(t, err)
			defer file.Close()
			var b bytes.Buffer
			_, err = b.ReadFrom(file)
			assert.NoError(t, err)
			res, err := m.Decode(b.Bytes())
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantRes, string(res))
		})
	}
}
