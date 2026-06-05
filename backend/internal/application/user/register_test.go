package user_test

import (
	"context"
	"errors"
	"testing"

	appuser "github.com/praxisvr/grama/internal/application/user"
	domainuser "github.com/praxisvr/grama/internal/domain/user"
)

// mockUserRepo is a simple in-memory UserRepository for unit tests.
// All test files in this package share this mock via package user_test.
type mockUserRepo struct {
	users   map[string]*domainuser.User
	emails  map[string]*domainuser.User
	nameLog map[string]string // tracks UpdateName calls: id -> newName
	err     error             // if set, Create returns this error
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users:   make(map[string]*domainuser.User),
		emails:  make(map[string]*domainuser.User),
		nameLog: make(map[string]string),
	}
}

func (m *mockUserRepo) Create(_ context.Context, u *domainuser.User) error {
	if m.err != nil {
		return m.err
	}
	m.users[u.ID().String()] = u
	m.emails[u.Email()] = u
	return nil
}

func (m *mockUserRepo) FindByEmail(_ context.Context, email string) (*domainuser.User, error) {
	if u, ok := m.emails[email]; ok {
		return u, nil
	}
	return nil, domainuser.ErrUserNotFound
}

func (m *mockUserRepo) FindByID(_ context.Context, id string) (*domainuser.User, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return nil, domainuser.ErrUserNotFound
}

func (m *mockUserRepo) UpdateName(_ context.Context, id, name string) error {
	if m.err != nil {
		return m.err
	}
	m.nameLog[id] = name
	return nil
}

func (m *mockUserRepo) FindByRole(_ context.Context, role string) ([]*domainuser.User, error) {
	var result []*domainuser.User
	for _, u := range m.users {
		if u.Role() == role {
			result = append(result, u)
		}
	}
	return result, nil
}

func TestRegisterUser_Valid_Owner(t *testing.T) {
	uc := appuser.NewRegisterUseCase(newMockUserRepo())
	out, err := uc.Execute(context.Background(), appuser.RegisterInput{
		Name:         "Alice",
		Email:        "alice@example.com",
		Password:     "password123",
		Role:         "owner",
		ConsentGiven: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Email != "alice@example.com" {
		t.Errorf("email mismatch: got %s", out.Email)
	}
	if out.Role != "owner" {
		t.Errorf("role mismatch: got %s", out.Role)
	}
}

func TestRegisterUser_Valid_Client(t *testing.T) {
	uc := appuser.NewRegisterUseCase(newMockUserRepo())
	out, err := uc.Execute(context.Background(), appuser.RegisterInput{
		Name:         "Bob",
		Email:        "bob@example.com",
		Password:     "password123",
		Role:         "client",
		ConsentGiven: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Role != "client" {
		t.Errorf("expected client role, got %s", out.Role)
	}
}

func TestRegisterUser_ConsentFalse(t *testing.T) {
	uc := appuser.NewRegisterUseCase(newMockUserRepo())
	_, err := uc.Execute(context.Background(), appuser.RegisterInput{
		Name:         "Carol",
		Email:        "carol@example.com",
		Password:     "password123",
		Role:         "client",
		ConsentGiven: false,
	})
	if !errors.Is(err, appuser.ErrConsentRequired) {
		t.Fatalf("expected ErrConsentRequired, got %v", err)
	}
}

func TestRegisterUser_OperatorBlocked(t *testing.T) {
	uc := appuser.NewRegisterUseCase(newMockUserRepo())
	_, err := uc.Execute(context.Background(), appuser.RegisterInput{
		Name:         "Dave",
		Email:        "dave@example.com",
		Password:     "password123",
		Role:         "operator",
		ConsentGiven: true,
	})
	if !errors.Is(err, appuser.ErrOperatorSelfRegister) {
		t.Fatalf("expected ErrOperatorSelfRegister, got %v", err)
	}
}

func TestRegisterUser_InvalidEmail(t *testing.T) {
	uc := appuser.NewRegisterUseCase(newMockUserRepo())
	_, err := uc.Execute(context.Background(), appuser.RegisterInput{
		Name:         "Eve",
		Email:        "not-an-email",
		Password:     "password123",
		Role:         "client",
		ConsentGiven: true,
	})
	if !errors.Is(err, domainuser.ErrEmailInvalid) {
		t.Fatalf("expected ErrEmailInvalid, got %v", err)
	}
}
