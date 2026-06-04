package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	appuser "github.com/praxisvr/grama/internal/application/user"
	domainuser "github.com/praxisvr/grama/internal/domain/user"
	"github.com/praxisvr/grama/internal/ports/repository"
)

// mockTokenRepo is an in-memory TokenRepository for unit tests.
type mockTokenRepo struct {
	tokens map[string]*repository.RefreshToken
}

func newMockTokenRepo() *mockTokenRepo {
	return &mockTokenRepo{tokens: make(map[string]*repository.RefreshToken)}
}

func (m *mockTokenRepo) SaveRefreshToken(_ context.Context, userID, tokenHash string, expiresAt time.Time) error {
	m.tokens[tokenHash] = &repository.RefreshToken{
		ID:        "mock-id",
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
	return nil
}

func (m *mockTokenRepo) FindRefreshToken(_ context.Context, tokenHash string) (*repository.RefreshToken, error) {
	t, ok := m.tokens[tokenHash]
	if !ok {
		return nil, errors.New("not found")
	}
	return t, nil
}

func (m *mockTokenRepo) RevokeRefreshToken(_ context.Context, tokenHash string) error {
	now := time.Now()
	if t, ok := m.tokens[tokenHash]; ok {
		t.RevokedAt = &now
	}
	return nil
}

func (m *mockTokenRepo) RevokeAllUserTokens(_ context.Context, _ string) error {
	return nil
}

// mockJWT satisfies the unexported jwtIssuer interface used by LoginUseCase.
type mockJWT struct{}

func (mockJWT) GenerateAccessToken(_, _ string) (string, error) {
	return "mock-access-token", nil
}

func TestLoginUser_Valid(t *testing.T) {
	repo := newMockUserRepo()
	u, err := domainuser.NewUser("Alice", "alice@example.com", "password123", "owner")
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	repo.emails["alice@example.com"] = u
	repo.users[u.ID().String()] = u

	uc := appuser.NewLoginUseCase(repo, newMockTokenRepo(), mockJWT{})
	out, err := uc.Execute(context.Background(), appuser.LoginInput{
		Email:    "alice@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.AccessToken == "" {
		t.Error("expected non-empty access token")
	}
	if out.RefreshToken == "" {
		t.Error("expected non-empty refresh token")
	}
	if out.ExpiresIn != 900 {
		t.Errorf("expected ExpiresIn=900, got %d", out.ExpiresIn)
	}
}

func TestLoginUser_WrongPassword(t *testing.T) {
	repo := newMockUserRepo()
	u, _ := domainuser.NewUser("Alice", "alice@example.com", "password123", "owner")
	repo.emails["alice@example.com"] = u
	repo.users[u.ID().String()] = u

	uc := appuser.NewLoginUseCase(repo, newMockTokenRepo(), mockJWT{})
	_, err := uc.Execute(context.Background(), appuser.LoginInput{
		Email:    "alice@example.com",
		Password: "wrongpassword",
	})
	if !errors.Is(err, appuser.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLoginUser_UserNotFound(t *testing.T) {
	uc := appuser.NewLoginUseCase(newMockUserRepo(), newMockTokenRepo(), mockJWT{})
	_, err := uc.Execute(context.Background(), appuser.LoginInput{
		Email:    "ghost@example.com",
		Password: "password123",
	})
	// Must return the same generic error regardless of whether the email exists.
	if !errors.Is(err, appuser.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials (not a user-existence leak), got %v", err)
	}
}
