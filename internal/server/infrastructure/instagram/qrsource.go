package instagram

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/hotafrika/griz-backend/internal/server/domain"
	"github.com/hotafrika/griz-backend/internal/server/infrastructure/instagram/photo"
	"github.com/pkg/errors"
	"sync"
	"time"
)

// QRSource is for getting QR codes from instagram post
type QRSource struct {
	decoder     domain.QRDecoder
	photoSource domain.PhotoSourcer
	client      *resty.Client
}

// NewQRSource ...
func NewQRSource() *QRSource {
	return &QRSource{
		photoSource: photo.NewPhotoSource(),
		client:      resty.New().SetTimeout(10 * time.Second),
	}
}

// GetFirstQR returns parsed data from first found code.
// It could be any kind of data. For our case we need to validate it.
func (qs QRSource) GetFirstQR(ctx context.Context, code string) (b []byte, err error) {
	links, err := qs.photoSource.GetPhotos(ctx, code)
	if err != nil {
		return nil, errors.Wrap(err, "qrencoder source: ")
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	c := make(chan []byte)
	var wg sync.WaitGroup

	for _, s := range links {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			b, err := qs.processImage(ctx, s)
			if err != nil {
				return
			}
			select {
			case c <- b:
			case <-ctx.Done():
				return
			}
		}(s)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	res, ok := <-c
	if !ok {
		return nil, errors.New("unable to find QR")
	}
	return res, nil
}

func (qs QRSource) processImage(ctx context.Context, link string) ([]byte, error) {
	b, err := qs.downloadImage(ctx, link)
	if err != nil {
		// TODO log here
		return nil, err
	}

	res, err := qs.decoder.Decode(b)
	if err != nil {
		// TODO log here
		return nil, err
	}

	return res, nil
}

func (qs QRSource) downloadImage(ctx context.Context, link string) ([]byte, error) {
	res, err := qs.client.R().SetContext(ctx).Get(link)
	if err != nil {
		return nil, err
	}
	if res.StatusCode() != 200 {
		return nil, errors.New("unable to download image. Broken link")
	}
	return res.Body(), nil
}
