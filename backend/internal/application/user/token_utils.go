package user

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/google/uuid"
)

// jwtIssuer is the port for JWT generation. Defined here (consumer side) so the
// application layer never imports from infrastructure. *auth.JWTManager satisfies
// this interface implicitly via Go structural typing.
type jwtIssuer interface {
	GenerateAccessToken(userID, role string) (string, error)
}

// generateRefreshToken creates a random plain-text UUID token and its SHA-256 hex hash.
// SHA-256 is sufficient: the input is already a high-entropy UUID.
func generateRefreshToken() (plain, hash string, err error) {
	plain = uuid.NewString()
	hash = hashRefreshToken(plain)
	return plain, hash, nil
}

// hashRefreshToken returns the SHA-256 hex digest of the plain token.
func hashRefreshToken(plain string) string {
	sum := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(sum[:])
}
