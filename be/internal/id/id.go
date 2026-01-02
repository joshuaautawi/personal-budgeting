package id

import (
	"crypto/rand"
	"encoding/hex"
)

type Generator interface {
	NewID() string
}

// RandomHex generates 16 random bytes and returns a 32-char hex string.
type RandomHex struct{}

func (RandomHex) NewID() string {
	var b [16]byte
	_, _ = rand.Read(b[:]) // best-effort; zero bytes are still valid IDs for local dev
	return hex.EncodeToString(b[:])
}
