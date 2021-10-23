package photo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/hotafrika/griz-backend/internal/server/domain"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var notNecScript = errors.New("script is not necessary")

//const instagramURL = "https://www.instagram.com/p/%v/embed/captioned/"
const instagramURL = "https://www.instagram.com/p/%v/embed/"

var _ domain.PhotoSourcer = (*Source)(nil)

// Source is the source of photos from Instagram
type Source struct {
	keyValidator *regexp.Regexp
	client       *resty.Client
	URL          string
}

// NewPhotoSource creates Source
func NewPhotoSource() Source {
	re := regexp.MustCompile(`^[0-9A-Za-z]+$`)
	client := resty.New().SetTimeout(10 * time.Second)
	return Source{
		keyValidator: re,
		client:       client,
		URL:          instagramURL,
	}
}

// GetPhotos returns links to instagram photos
func (s Source) GetPhotos(ctx context.Context, key string) (links []string, err error) {
	err = s.validateKey(key)
	if err != nil {
		return nil, errors.Wrap(err, "key validation: ")
	}

	ua := browser.Chrome()
	URL := fmt.Sprintf(s.URL, key)

	res, err := s.client.R().
		SetHeader("user-agent", ua).
		SetContext(ctx).
		Get(URL)

	if err != nil {
		return nil, errors.Wrap(err, "http request: ")
	}
	if res.StatusCode() != http.StatusOK {
		return nil, errors.New("http request status not OK")
	}

	var embedRes EmbedResponse

	err = parseBodyByScript(&embedRes, bytes.NewReader(res.Body()))

	return embedRes.getURLs(), err
}

func (s Source) validateKey(key string) error {
	if !s.keyValidator.Match([]byte(key)) {
		return errors.New("key has wrong format")
	}
	return nil
}

func parseBodyByScript(er *EmbedResponse, r io.Reader) (err error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return errors.Wrap(err, "goquery parsing: ")
	}

	doc.Find("script").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		err := checkScript(er, selection.Text())
		if err != nil {
			if !errors.Is(err, notNecScript) {
				//TODO log here
			}
			return true
		}
		return false
	})

	if er.IsEmpty() {
		src, ok := doc.Find(".EmbeddedMediaImage").First().Attr("src")
		if ok {
			er.Media.DisplayURL = src
		}
	}

	return
}

func checkScript(er *EmbedResponse, scriptContent string) (err error) {
	necPrefix := "window.__additionalDataLoaded('extra',"
	if !strings.HasPrefix(scriptContent, necPrefix) {
		return notNecScript
	}
	res := strings.TrimPrefix(scriptContent, necPrefix)
	res = strings.TrimSuffix(res, ");")

	err = json.Unmarshal([]byte(res), er)
	return
}
