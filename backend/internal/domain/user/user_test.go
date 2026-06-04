package user_test

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/praxisvr/grama/internal/domain/user"
)

func TestNewUser_Valid(t *testing.T) {
	u, err := user.NewUser("Alice", "Alice@Example.com", "secure123", "owner")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if u.Email() != "alice@example.com" {
		t.Errorf("email not lowercased: got %s", u.Email())
	}
	if u.Role() != "owner" {
		t.Errorf("unexpected role: got %s", u.Role())
	}
	if u.ID().String() == "" {
		t.Error("expected non-empty UUID")
	}
	if u.Name() != "Alice" {
		t.Errorf("unexpected name: got %s", u.Name())
	}
}

func TestNewUser_InvalidEmail(t *testing.T) {
	_, err := user.NewUser("Bob", "not-an-email", "secure123", "client")
	if err != user.ErrEmailInvalid {
		t.Fatalf("expected ErrEmailInvalid, got %v", err)
	}
}

func TestNewUser_PasswordTooShort(t *testing.T) {
	_, err := user.NewUser("Bob", "bob@example.com", "short", "client")
	if err != user.ErrPasswordTooShort {
		t.Fatalf("expected ErrPasswordTooShort, got %v", err)
	}
}

func TestNewUser_InvalidRole(t *testing.T) {
	_, err := user.NewUser("Bob", "bob@example.com", "secure123", "superadmin")
	if err != user.ErrInvalidRole {
		t.Fatalf("expected ErrInvalidRole, got %v", err)
	}
}

func TestCheckPassword_Correct(t *testing.T) {
	u, err := user.NewUser("Carol", "carol@example.com", "mypassword", "client")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !u.CheckPassword("mypassword") {
		t.Error("CheckPassword returned false for the correct password")
	}
}

func TestCheckPassword_Wrong(t *testing.T) {
	u, err := user.NewUser("Carol", "carol@example.com", "mypassword", "client")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.CheckPassword("wrongpassword") {
		t.Error("CheckPassword returned true for a wrong password")
	}
}

func TestReconstitute(t *testing.T) {
	original, err := user.NewUser("Dave", "dave@example.com", "mypassword", "operator")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Simulate DB round-trip: the repository hex-encodes before INSERT and hex-decodes before Reconstitute.
	hexHash := hex.EncodeToString(original.PasswordHash())

	now := time.Now().UTC()
	reconstructed := user.Reconstitute(
		original.ID().String(),
		original.Name(),
		original.Email(),
		hexHash,
		original.Role(),
		false,
		now,
		original.CreatedAt(),
		original.UpdatedAt(),
	)

	if reconstructed.Email() != original.Email() {
		t.Errorf("email mismatch: got %s, want %s", reconstructed.Email(), original.Email())
	}
	if reconstructed.Role() != original.Role() {
		t.Errorf("role mismatch: got %s, want %s", reconstructed.Role(), original.Role())
	}
	if !reconstructed.CheckPassword("mypassword") {
		t.Error("Reconstitute: CheckPassword failed for the original password — hash was re-derived instead of restored")
	}
	if reconstructed.CheckPassword("wrongpassword") {
		t.Error("Reconstitute: CheckPassword accepted a wrong password")
	}
}
