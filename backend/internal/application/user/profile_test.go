package user_test

import (
	"context"
	"errors"
	"testing"

	appuser "github.com/praxisvr/grama/internal/application/user"
	domainuser "github.com/praxisvr/grama/internal/domain/user"
)

func TestGetProfile_Valid(t *testing.T) {
	repo := newMockUserRepo()
	u, err := domainuser.NewUser("Alice", "alice@example.com", "password123", "owner")
	if err != nil {
		t.Fatal(err)
	}
	repo.users[u.ID().String()] = u
	repo.emails[u.Email()] = u

	uc := appuser.NewGetProfileUseCase(repo)
	out, err := uc.Execute(context.Background(), u.ID().String())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Email != "alice@example.com" {
		t.Errorf("email mismatch: %s", out.Email)
	}
	if out.Role != "owner" {
		t.Errorf("role mismatch: %s", out.Role)
	}
	// Verify output struct has no PasswordHash field.
	// The GetProfileOutput type itself enforces this at compile time.
}

func TestUpdateProfile_Valid(t *testing.T) {
	repo := newMockUserRepo()
	u, _ := domainuser.NewUser("Alice", "alice@example.com", "password123", "owner")
	repo.users[u.ID().String()] = u

	uc := appuser.NewUpdateProfileUseCase(repo)
	err := uc.Execute(context.Background(), u.ID().String(), appuser.UpdateProfileInput{Name: "Alice Bonilla"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.nameLog[u.ID().String()] != "Alice Bonilla" {
		t.Errorf("UpdateName not called with correct name")
	}
}

func TestUpdateProfile_EmptyName(t *testing.T) {
	uc := appuser.NewUpdateProfileUseCase(newMockUserRepo())
	err := uc.Execute(context.Background(), "any-id", appuser.UpdateProfileInput{Name: ""})
	if !errors.Is(err, appuser.ErrNameRequired) {
		t.Fatalf("expected ErrNameRequired, got %v", err)
	}
}

func TestCreateOperator_Valid(t *testing.T) {
	uc := appuser.NewCreateOperatorUseCase(newMockUserRepo())
	out, err := uc.Execute(context.Background(), "owner-123", appuser.CreateOperatorInput{
		Name:     "Carlos",
		Email:    "carlos@grama.com",
		Password: "operador123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Role != "operator" {
		t.Errorf("expected role=operator, got %s", out.Role)
	}
}

func TestCreateOperator_EmailInvalid(t *testing.T) {
	uc := appuser.NewCreateOperatorUseCase(newMockUserRepo())
	_, err := uc.Execute(context.Background(), "owner-123", appuser.CreateOperatorInput{
		Name:     "Carlos",
		Email:    "not-an-email",
		Password: "operador123",
	})
	if !errors.Is(err, domainuser.ErrEmailInvalid) {
		t.Fatalf("expected ErrEmailInvalid, got %v", err)
	}
}
