package resources

import (
	"github.com/pkg/errors"
	"net/url"
)

// CodeCreateRequest ...
type CodeCreateRequest struct {
	URL string `json:"url"`
}

// Validate ...
func (r CodeCreateRequest) Validate() error {
	_, err := url.ParseRequestURI(r.URL)
	return errors.Wrap(err, "URL validation: ")
}

// CodeCreateResponse ...
type CodeCreateResponse struct {
	ID uint64 `json:"id"`
}

// UpdateCodeRequest ...
type UpdateCodeRequest struct {
	ID  uint64 `json:"id"`
	URL string `json:"url"`
}

// Validate ...
func (r UpdateCodeRequest) Validate() error {
	_, err := url.ParseRequestURI(r.URL)
	if err != nil {
		return errors.Wrap(err, "URL validation: ")
	}
	if r.ID == 0 {
		return errors.New("ID validation: zero value")
	}
	return nil
}
