package domain

import (
	"context"
	"fmt"
	"time"
)

var ErrCacheNotExist = fmt.Errorf("key not exist")

type Cacher interface {
	Get(context.Context, fmt.Stringer) (string, error)
	Set(context.Context, fmt.Stringer, string, time.Duration) error
	Delete(context.Context, fmt.Stringer) error
}
