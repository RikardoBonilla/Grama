package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
)

// GenerateRefreshToken creates a new random refresh token.
// Returns the plaintext token (sent to the client) and its SHA-256 hex hash (stored in DB).
// SHA-256 is sufficient here because the token is already a high-entropy random UUID.
func GenerateRefreshToken() (plain string, hash string, err error) {
	plain = uuid.NewString()
	if plain == "" {
		return "", "", fmt.Errorf("uuid generation failed")
	}
	hash = HashRefreshToken(plain)
	return plain, hash, nil
}

// HashRefreshToken returns the SHA-256 hex digest of the plaintext token.
func HashRefreshToken(plain string) string {
	sum := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(sum[:])
}
