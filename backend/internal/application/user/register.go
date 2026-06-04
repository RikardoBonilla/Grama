package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	domainuser "github.com/praxisvr/grama/internal/domain/user"
	"github.com/praxisvr/grama/internal/ports/repository"
)

var (
	ErrConsentRequired      = errors.New("consent required")
	ErrOperatorSelfRegister = errors.New("operators cannot self-register")
)

type RegisterInput struct {
	Name         string
	Email        string
	Password     string
	Role         string
	ConsentGiven bool
}

type RegisterOutput struct {
	ID        string
	Name      string
	Email     string
	Role      string
	CreatedAt time.Time
}

type RegisterUseCase struct {
	repo repository.UserRepository
}

func NewRegisterUseCase(repo repository.UserRepository) *RegisterUseCase {
	return &RegisterUseCase{repo: repo}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, input RegisterInput) (*RegisterOutput, error) {
	if !input.ConsentGiven {
		return nil, ErrConsentRequired
	}
	if input.Role == "operator" {
		return nil, ErrOperatorSelfRegister
	}

	u, err := domainuser.NewUser(input.Name, input.Email, input.Password, input.Role)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, u); err != nil {
		// Do not reveal whether the email already exists.
		log.Printf("register: repo.Create failed: %v", err)
		return nil, fmt.Errorf("registration failed")
	}

	return &RegisterOutput{
		ID:        u.ID().String(),
		Name:      u.Name(),
		Email:     u.Email(),
		Role:      u.Role(),
		CreatedAt: u.CreatedAt(),
	}, nil
}
