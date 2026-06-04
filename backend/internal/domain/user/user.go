package user

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

const (
	argon2Memory      uint32 = 65536
	argon2Iterations  uint32 = 3
	argon2Parallelism uint8  = 2
	argon2SaltLen            = 16
	argon2KeyLen      uint32 = 32
	minPasswordLen           = 8
)

var validRoles = map[string]struct{}{
	"owner":    {},
	"operator": {},
	"client":   {},
}

// User is the aggregate root of the Identity bounded context.
type User struct {
	id             uuid.UUID
	name           string
	email          string
	passwordHash   []byte // stored format: salt (16 B) || argon2id key (32 B)
	role           string
	consentGiven   bool
	consentGivenAt time.Time
	createdAt      time.Time
	updatedAt      time.Time
}

// NewUser creates and validates a new User, hashing the password with Argon2id.
func NewUser(name, email, password, role string) (*User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if !isValidEmail(email) {
		return nil, ErrEmailInvalid
	}
	if len(password) < minPasswordLen {
		return nil, ErrPasswordTooShort
	}
	if _, ok := validRoles[role]; !ok {
		return nil, ErrInvalidRole
	}

	hash, err := hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	now := time.Now().UTC()
	return &User{
		id:           uuid.New(),
		name:         name,
		email:        email,
		passwordHash: hash,
		role:         role,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

// Reconstitute rebuilds a User from persisted data without re-hashing.
// passwordHash must be the hex-encoded bytes as stored in the database.
func Reconstitute(
	id, name, email, passwordHash, role string,
	consentGiven bool,
	consentGivenAt, createdAt, updatedAt time.Time,
) *User {
	parsedID, _ := uuid.Parse(id)
	decoded, _ := hex.DecodeString(passwordHash)
	return &User{
		id:             parsedID,
		name:           name,
		email:          email,
		passwordHash:   decoded,
		role:           role,
		consentGiven:   consentGiven,
		consentGivenAt: consentGivenAt,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}

// CheckPassword verifies a plaintext password against the stored hash in constant time.
func (u *User) CheckPassword(plain string) bool {
	if len(u.passwordHash) < argon2SaltLen+int(argon2KeyLen) {
		return false
	}
	salt := u.passwordHash[:argon2SaltLen]
	stored := u.passwordHash[argon2SaltLen:]
	derived := argon2.IDKey([]byte(plain), salt, argon2Iterations, argon2Memory, argon2Parallelism, argon2KeyLen)
	return subtle.ConstantTimeCompare(stored, derived) == 1
}

// ID returns the user's UUID.
func (u *User) ID() uuid.UUID { return u.id }

// Name returns the user's display name.
func (u *User) Name() string { return u.name }

// Email returns the user's email address (always lowercase).
func (u *User) Email() string { return u.email }

// Role returns the user's role: owner, operator, or client.
func (u *User) Role() string { return u.role }

// PasswordHash returns the raw hash bytes for persistence only — never include in HTTP responses.
func (u *User) PasswordHash() []byte { return u.passwordHash }

// ConsentGiven returns whether the user gave data-treatment consent (Ley 1581 de 2012).
func (u *User) ConsentGiven() bool { return u.consentGiven }

// ConsentGivenAt returns when data-treatment consent was recorded.
func (u *User) ConsentGivenAt() time.Time { return u.consentGivenAt }

// CreatedAt returns the UTC creation timestamp.
func (u *User) CreatedAt() time.Time { return u.createdAt }

// UpdatedAt returns the UTC last-update timestamp.
func (u *User) UpdatedAt() time.Time { return u.updatedAt }

// hashPassword generates a random salt and derives a key with Argon2id.
// Returned format: salt (16 B) || key (32 B).
func hashPassword(password string) ([]byte, error) {
	salt := make([]byte, argon2SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("generate salt: %w", err)
	}
	key := argon2.IDKey([]byte(password), salt, argon2Iterations, argon2Memory, argon2Parallelism, argon2KeyLen)
	return append(salt, key...), nil
}

// isValidEmail performs a minimal structural check (user@domain.tld).
func isValidEmail(email string) bool {
	if strings.Count(email, "@") != 1 {
		return false
	}
	at := strings.Index(email, "@")
	if at < 1 {
		return false
	}
	rest := email[at+1:]
	dot := strings.LastIndex(rest, ".")
	return dot > 0 && dot < len(rest)-1
}
