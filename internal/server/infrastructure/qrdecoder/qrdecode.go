package qrdecoder

import (
	"bytes"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/pkg/errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
)

// Makiuchi is a qr decoder of Makiuchi user
type Makiuchi struct {
}

// Decode decodes slice of bytes (potential image) to result string as slice of bytes
func (m Makiuchi) Decode(b []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, errors.Wrap(err, "unable to decode bytes to image: ")
	}
	bitmap, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create binary bitmap: ")
	}
	result, err := qrcode.NewQRCodeReader().DecodeWithoutHints(bitmap)
	if err != nil {
		return nil, errors.Wrap(err, "unable to decode binary bitmap: ")
	}
	return []byte(result.GetText()), nil
}
