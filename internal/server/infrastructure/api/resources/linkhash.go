package resources

import (
	"errors"
	"net/url"
	"strings"
)

// LinkHashRequest is used fot parsing requests
type LinkHashRequest struct {
	Type string
	URL  string `json:"url"`
}

// Parse return parsed token
func (sl LinkHashRequest) Parse() (string, error) {
	link, err := url.ParseRequestURI(sl.URL)
	if err != nil {
		return "", err
	}
	if !strings.HasSuffix(sl.URL, "https://griz.grizzlytics.com/app?d=") {
		return "", errors.New("not griz link")
	}
	//if !(link.Host == "griz.grizzlytics.com") {
	//	return "", errors.New("not griz link")
	//}
	m, err := url.ParseQuery(link.RawQuery)
	if err != nil {
		return "", errors.New("query is not right")
	}
	res := m.Get("d")
	if res == "" {
		return "", errors.New("query doesn't contain necessary params")
	}
	return res, nil
}

// LinkHashResponse serves responses
type LinkHashResponse struct {
	URL string `json:"url"`
}
