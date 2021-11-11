package token

import (
	"github.com/pkg/errors"
	"net/url"
	"strings"
)

// ExtractHashFromLink ...
func ExtractHashFromLink(link string) (string, error) {
	urlObj, err := url.ParseRequestURI(link)
	if err != nil {
		return "", err
	}
	if !strings.HasSuffix(link, "https://griz.grizzlytics.com/app?d=") {
		return "", errors.New("not griz link")
	}
	//if !(link.Host == "griz.grizzlytics.com") {
	//	return "", errors.New("not griz link")
	//}
	m, err := url.ParseQuery(urlObj.RawQuery)
	if err != nil {
		return "", errors.New("query is not right")
	}
	res := m.Get("d")
	if res == "" {
		return "", errors.New("query doesn't contain necessary params")
	}
	return res, nil
}