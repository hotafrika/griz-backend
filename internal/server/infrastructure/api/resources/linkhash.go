package resources

import (
	"github.com/hotafrika/griz-backend/internal/server/app/token"
)

// LinkHashRequest is used fot parsing requests
type LinkHashRequest struct {
	Type string
	URL  string `json:"url"`
}

// Parse return parsed token
func (sl LinkHashRequest) Parse() (string, error) {
	return token.ExtractHashFromLink(sl.URL)
}

// LinkHashResponse serves responses
type LinkHashResponse struct {
	URL string `json:"url"`
}
