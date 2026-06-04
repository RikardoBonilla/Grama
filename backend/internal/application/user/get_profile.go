package user

import (
	"context"
	"fmt"
	"time"

	"github.com/praxisvr/grama/internal/ports/repository"
)

type GetProfileOutput struct {
	ID        string
	Name      string
	Email     string
	Role      string
	CreatedAt time.Time
}

type GetProfileUseCase struct {
	repo repository.UserRepository
}

func NewGetProfileUseCase(repo repository.UserRepository) *GetProfileUseCase {
	return &GetProfileUseCase{repo: repo}
}

func (uc *GetProfileUseCase) Execute(ctx context.Context, userID string) (*GetProfileOutput, error) {
	u, err := uc.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get profile: %w", err)
	}
	// NEVER include PasswordHash in the output.
	return &GetProfileOutput{
		ID:        u.ID().String(),
		Name:      u.Name(),
		Email:     u.Email(),
		Role:      u.Role(),
		CreatedAt: u.CreatedAt(),
	}, nil
}
