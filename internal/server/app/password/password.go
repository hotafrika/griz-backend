package password

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// Encryptor works with passwords hashing
type Encryptor struct {
	key []byte
}

// NewEncryptor creates new Encryptor
func NewEncryptor(b []byte) Encryptor {
	return Encryptor{
		key: []byte(b),
	}
}

// NewEncryptorByString creates new Encryptor
func NewEncryptorByString(s string) Encryptor {
	return NewEncryptor([]byte(s))
}

// Encode encodes bytes to hex byte slice
func (e Encryptor) Encode(raw []byte) ([]byte, error) {
	h := hmac.New(sha256.New, e.key)
	_, err := h.Write(raw)
	if err != nil {
		return nil, err
	}
	src := h.Sum(nil)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return dst, nil
}

// EncodeString encodes string to hex byte slice
func (e Encryptor) EncodeString(s string) ([]byte, error) {
	return e.Encode([]byte(s))
}
