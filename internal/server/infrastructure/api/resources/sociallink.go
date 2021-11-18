package resources

import (
	"github.com/pkg/errors"
	"net/url"
)

// SocialLinkRequest is used fot parsing requests
type SocialLinkRequest struct {
	Type string
	URL  string `json:"url"`
}

// ValidateIsInsta ...
func (sl SocialLinkRequest) ValidateIsInsta() error {
	link, err := url.ParseRequestURI(sl.URL)
	if err != nil {
		return err
	}
	if !(link.Host == "instagram.com" || link.Host == "www.instagram.com") {
		return errors.New("not instagram link")
	}
	return nil
}

// SocialLinkResponse serves responses
type SocialLinkResponse struct {
	URL string `json:"url"`
}
