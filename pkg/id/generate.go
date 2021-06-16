package id

import (
	"encoding/hex"
	"math/rand"
	"time"
)

func New(length int) (string, error) {
	gen := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, (length + 1) / 2)

	if _, err := gen.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b)[:length]
}
