package user

import (
	"context"
	"errors"

	"github.com/praxisvr/grama/internal/ports/repository"
)

var ErrNameRequired = errors.New("name is required")

type UpdateProfileInput struct {
	Name string // only name is mutable; email and role are immutable
}

type UpdateProfileUseCase struct {
	repo repository.UserRepository
}

func NewUpdateProfileUseCase(repo repository.UserRepository) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{repo: repo}
}

func (uc *UpdateProfileUseCase) Execute(ctx context.Context, userID string, input UpdateProfileInput) error {
	if input.Name == "" {
		return ErrNameRequired
	}
	return uc.repo.UpdateName(ctx, userID, input.Name)
}
