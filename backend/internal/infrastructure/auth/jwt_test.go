package auth_test

import (
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/praxisvr/grama/internal/infrastructure/auth"
)

func TestGenerateAndValidate_Valid(t *testing.T) {
	mgr := auth.NewJWTManager("super-secret-key-min-32-chars-xx")
	token, err := mgr.GenerateAccessToken("user-123", "owner")
	if err != nil {
		t.Fatalf("GenerateAccessToken: %v", err)
	}

	claims, err := mgr.ValidateAccessToken(token)
	if err != nil {
		t.Fatalf("ValidateAccessToken: %v", err)
	}
	if claims.Subject != "user-123" {
		t.Errorf("subject mismatch: got %s", claims.Subject)
	}
	if claims.Role != "owner" {
		t.Errorf("role mismatch: got %s", claims.Role)
	}
}

func TestValidate_AlgNone(t *testing.T) {
	// Build a token with alg:none manually to simulate the attack.
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(
		fmt.Sprintf(`{"sub":"attacker","exp":%d}`, time.Now().Add(time.Hour).Unix()),
	))
	algNoneToken := header + "." + payload + "."

	mgr := auth.NewJWTManager("super-secret-key-min-32-chars-xx")
	_, err := mgr.ValidateAccessToken(algNoneToken)
	if err == nil {
		t.Fatal("expected error for alg:none token, got nil — attack not blocked")
	}
}

func TestValidate_ExpiredToken(t *testing.T) {
	// Forge an already-expired HS256 token signed with the correct secret.
	secret := "super-secret-key-min-32-chars-xx"
	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-123",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			ID:        "test-jti",
		},
		Role: "client",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	mgr := auth.NewJWTManager(secret)
	_, err = mgr.ValidateAccessToken(signed)
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
}

func TestValidate_WrongSecret(t *testing.T) {
	mgr := auth.NewJWTManager("super-secret-key-min-32-chars-xx")
	token, err := mgr.GenerateAccessToken("user-123", "client")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}

	mgrWrong := auth.NewJWTManager("wrong-secret-key-min-32-chars-xx")
	_, err = mgrWrong.ValidateAccessToken(token)
	if err == nil {
		t.Fatal("expected error for wrong secret, got nil")
	}
}
