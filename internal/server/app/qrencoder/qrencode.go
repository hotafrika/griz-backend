package qrencoder

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/yeqown/go-qrcode"
	"path"
)

// YeqownOption is options for encoder of type Yeqown
type YeqownOption func(*Yeqown)

// WithQRWidth adds size for blocks
func WithQRWidth(n uint8) YeqownOption {
	return func(yeqown *Yeqown) {
		yeqown.options = append(yeqown.options, qrcode.WithQRWidth(n))
	}
}

// WithFileImagePNG adds PNG image in the center of QR.
// Width of file couldn't be more than 1/5 of QR width
// Put this image to img folder. Use just image name in this option.
func WithFileImagePNG(filename string) YeqownOption {
	return func(yeqown *Yeqown) {
		yeqown.options = append(yeqown.options, qrcode.WithLogoImageFilePNG(path.Join("img", filename)))
	}
}

// Yeqown type of QR code encoder
type Yeqown struct {
	options []qrcode.ImageOption
}

// NewYeqown creates new QR encoder
func NewYeqown(options ...YeqownOption) Yeqown {
	y := Yeqown{}
	for _, option := range options {
		option(&y)
	}
	return y
}

// DefaultYeqown creates new default encoder
func DefaultYeqown() Yeqown {
	y := Yeqown{}
	options := []YeqownOption{WithQRWidth(6), WithFileImagePNG("GrizLogo.png")}
	for _, option := range options {
		option(&y)
	}
	return y
}

// Encode returns slice of bytes with QR code
func (y Yeqown) Encode(b []byte) ([]byte, error) {
	qr, err := qrcode.New(string(b), y.options...)
	if err != nil {
		return nil, errors.Wrap(err, "qr encoder generation: ")
	}
	var buf bytes.Buffer

	err = qr.SaveTo(&buf)
	if err != nil {
		return nil, errors.Wrap(err, "qr encoder buffering: ")
	}
	return buf.Bytes(), nil
}
