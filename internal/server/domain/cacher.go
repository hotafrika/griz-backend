package domain

import (
	"fmt"
	"time"
)

type Cacher interface {
	Get(fmt.Stringer) (string, error)
	Set(fmt.Stringer, string, time.Duration) error
	Delete(fmt.Stringer) error
}
