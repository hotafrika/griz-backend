package domain

import "context"

// QRSourcer is interface for getting QR-data from different sources
type QRSourcer interface {
	Get(context.Context, string) ([]string, error)
}
