package token

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// AES creates tokens by key according to aes crypto library
type AES struct {
	aesblock cipher.Block
}

// NewAES creates new AES with set key
func NewAES(key string) (*AES, error) {
	keyB := []byte(key)
	if len(keyB) != aes.BlockSize {
		return nil, errors.New("key has to be 16 bytes long")
	}
	block, err := aes.NewCipher(keyB)
	if err != nil {
		return nil, err
	}
	return &AES{aesblock: block}, nil
}

// Create returns token as a string
func (a AES) Create(id uint64) (string, error) {
	if id < 1 {
		return "", errors.New("minimal value for id is 1")
	}

	src := strconv.FormatUint(id, 16) //Always max 16 bytes
	src = fmt.Sprintf("%16s", src)

	srcBytes := []byte(src)
	if len(srcBytes) != aes.BlockSize {
		return "", errors.New("src has to be 16 bytes long")
	}

	dst := make([]byte, aes.BlockSize)
	a.aesblock.Encrypt(dst, []byte(src))
	return "v01" + hex.EncodeToString(dst), nil
}

// Decode decodes token to id (uint64)
func (a AES) Decode(token string) (uint64, error) {
	if !strings.HasPrefix(token, "v01") {
		return 0, errors.New("token has wrong format")
	}
	b, err := hex.DecodeString(strings.TrimLeft(token, "v01"))
	if err != nil {
		return 0, err
	}
	if len(b) != aes.BlockSize {
		return 0, errors.New("decoded token has to be 16 bytes")
	}
	dst := make([]byte, aes.BlockSize)
	a.aesblock.Decrypt(dst, b)
	return strconv.ParseUint(strings.TrimLeft(string(dst), " "), 16, 64)
}
