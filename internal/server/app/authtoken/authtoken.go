package authtoken

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"time"
)

// JWT for creating tokens
type JWT struct {
	Duration time.Duration
	key      []byte
}

// NewJWT creates new JWT tokenizer
func NewJWT(key []byte, d time.Duration) JWT {
	return JWT{
		Duration: d,
		key:      key,
	}
}

// NewJWTFromString creates new JWT tokenizer
func NewJWTFromString(skey string, d time.Duration) JWT {
	return NewJWT([]byte(skey), d)
}

// MakeByID creates new token
func (j JWT) MakeByID(id uint64) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["id"] = id
	atClaims["exp"] = time.Now().Add(j.Duration).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString(j.key)
	if err != nil {
		return "", errors.Wrap(err, "SignedString: ")
	}
	return token, nil
}
