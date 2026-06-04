package user

import (
	"context"
	"errors"
	"time"

	domainuser "github.com/praxisvr/grama/internal/domain/user"
	"github.com/praxisvr/grama/internal/ports/repository"
)

var ErrInvalidToken = errors.New("invalid or expired token")

type RefreshInput struct {
	RefreshToken string
}

type RefreshOutput struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}

type RefreshUseCase struct {
	userRepo   repository.UserRepository
	tokenRepo  repository.TokenRepository
	jwtManager jwtIssuer
}

func NewRefreshUseCase(
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
	jwtManager jwtIssuer,
) *RefreshUseCase {
	return &RefreshUseCase{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwtManager: jwtManager,
	}
}

func (uc *RefreshUseCase) Execute(ctx context.Context, input RefreshInput) (*RefreshOutput, error) {
	hash := hashRefreshToken(input.RefreshToken)

	stored, err := uc.tokenRepo.FindRefreshToken(ctx, hash)
	if err != nil {
		return nil, ErrInvalidToken
	}
	if stored.RevokedAt != nil || time.Now().After(stored.ExpiresAt) {
		return nil, ErrInvalidToken
	}

	// Rotating refresh: revoke the current token before issuing a new pair.
	if err := uc.tokenRepo.RevokeRefreshToken(ctx, hash); err != nil {
		return nil, err
	}

	u, err := uc.userRepo.FindByID(ctx, stored.UserID)
	if err != nil {
		if errors.Is(err, domainuser.ErrUserNotFound) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	accessToken, err := uc.jwtManager.GenerateAccessToken(u.ID().String(), u.Role())
	if err != nil {
		return nil, err
	}

	plain, newHash, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().UTC().Add(7 * 24 * time.Hour)
	if err := uc.tokenRepo.SaveRefreshToken(ctx, u.ID().String(), newHash, expiresAt); err != nil {
		return nil, err
	}

	return &RefreshOutput{
		AccessToken:  accessToken,
		RefreshToken: plain,
		ExpiresIn:    900,
	}, nil
}
