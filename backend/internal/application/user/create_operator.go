package user

import (
	"context"
	"fmt"

	domainuser "github.com/praxisvr/grama/internal/domain/user"
	"github.com/praxisvr/grama/internal/ports/repository"
)

type CreateOperatorInput struct {
	Name     string
	Email    string
	Password string
	// VenueID string -- TODO GRAM-17 Sprint 2: assign operator to a venue
}

type CreateOperatorOutput struct {
	ID    string
	Name  string
	Email string
	Role  string
}

type CreateOperatorUseCase struct {
	repo repository.UserRepository
}

func NewCreateOperatorUseCase(repo repository.UserRepository) *CreateOperatorUseCase {
	return &CreateOperatorUseCase{repo: repo}
}

func (uc *CreateOperatorUseCase) Execute(
	ctx context.Context, ownerID string, input CreateOperatorInput,
) (*CreateOperatorOutput, error) {
	// Owners accept data-processing consent on behalf of their operators (internal staff).
	u, err := domainuser.NewUser(input.Name, input.Email, input.Password, "operator")
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, u); err != nil {
		fmt.Printf("create operator: repo.Create failed for owner %s: %v\n", ownerID, err)
		return nil, fmt.Errorf("operator creation failed")
	}

	// TODO GRAM-17 Sprint 2: persist venue_operators(owner_id, operator_id, venue_id)

	return &CreateOperatorOutput{
		ID:    u.ID().String(),
		Name:  u.Name(),
		Email: u.Email(),
		Role:  u.Role(),
	}, nil
}
