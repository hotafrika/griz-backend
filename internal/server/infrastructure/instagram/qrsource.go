package instagram

import (
	"github.com/hotafrika/griz-backend/internal/server/domain"
	"github.com/hotafrika/griz-backend/internal/server/infrastructure/instagram/photo"
)

// QRSource is for getting QR codes from instagram post
type QRSource struct {
	photoSource domain.PhotoSourcer
}

func NewQRSource() *QRSource {
	return &QRSource{photoSource: photo.NewPhotoSource()}
}


