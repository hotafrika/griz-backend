package token

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExtractHashFromLink(t *testing.T) {
	tests := []struct {
		url      string
		wantErr  bool
		wantHash string
	}{
		{
			url:      "https://griz.grizzlytics.com/app?d=v015cf58619ad623291c8c3b26c108720f7",
			wantErr:  false,
			wantHash: "v015cf58619ad623291c8c3b26c108720f7",
		},
		{
			url:      "https://localhost/app?d=v015cf58619ad623291c8c3b26c108720f7",
			wantErr:  true,
			wantHash: "",
		},
		{
			url:      "https://localhost/app?b=v015cf58619ad623291c8c3b26c108720f7",
			wantErr:  true,
			wantHash: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			link, err := ExtractHashFromLink(tt.url)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, link, tt.wantHash)
				}
			}
		})
	}
}
