package user

import (
	"context"

	"github.com/praxisvr/grama/internal/ports/repository"
)

type OperatorItem struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ListOperatorsOutput struct {
	Operators []OperatorItem
}

type ListOperatorsUseCase struct {
	repo repository.UserRepository
}

func NewListOperatorsUseCase(repo repository.UserRepository) *ListOperatorsUseCase {
	return &ListOperatorsUseCase{repo: repo}
}

func (uc *ListOperatorsUseCase) Execute(ctx context.Context, ownerID string) (*ListOperatorsOutput, error) {
	// TODO GRAM-17 Sprint 2: filter operators by owner's venues instead of listing all.
	users, err := uc.repo.FindByRole(ctx, "operator")
	if err != nil {
		return nil, err
	}

	items := make([]OperatorItem, 0, len(users))
	for _, u := range users {
		items = append(items, OperatorItem{
			ID:    u.ID().String(),
			Name:  u.Name(),
			Email: u.Email(),
		})
	}
	return &ListOperatorsOutput{Operators: items}, nil
}
