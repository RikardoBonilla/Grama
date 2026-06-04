package repository

import (
	"context"

	"github.com/praxisvr/grama/internal/domain/user"
)

type UserRepository interface {
	Create(ctx context.Context, u *user.User) error
	FindByEmail(ctx context.Context, email string) (*user.User, error)
	FindByID(ctx context.Context, id string) (*user.User, error)
	UpdateName(ctx context.Context, id string, name string) error
	FindByRole(ctx context.Context, role string) ([]*user.User, error)
}
