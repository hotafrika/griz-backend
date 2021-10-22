package domain

import "context"

type PhotoSourcer interface {
	GetPhotos(context.Context, string) ([]string, error)
}
