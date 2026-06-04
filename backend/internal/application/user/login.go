package user

import (
	"context"
	"errors"
	"log"
	"time"

	domainuser "github.com/praxisvr/grama/internal/domain/user"
	"github.com/praxisvr/grama/internal/ports/repository"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
	UserID       string
	UserName     string
	UserRole     string
}

type LoginUseCase struct {
	userRepo   repository.UserRepository
	tokenRepo  repository.TokenRepository
	jwtManager jwtIssuer
}

func NewLoginUseCase(
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
	jwtManager jwtIssuer,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwtManager: jwtManager,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	u, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, domainuser.ErrUserNotFound) {
			// Generic response — do not reveal whether the email exists.
			return nil, ErrInvalidCredentials
		}
		log.Printf("login: FindByEmail error: %v", err)
		return nil, ErrInvalidCredentials
	}

	if !u.CheckPassword(input.Password) {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := uc.jwtManager.GenerateAccessToken(u.ID().String(), u.Role())
	if err != nil {
		return nil, err
	}

	plain, hash, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().UTC().Add(7 * 24 * time.Hour)
	if err := uc.tokenRepo.SaveRefreshToken(ctx, u.ID().String(), hash, expiresAt); err != nil {
		return nil, err
	}

	return &LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: plain,
		ExpiresIn:    900,
		UserID:       u.ID().String(),
		UserName:     u.Name(),
		UserRole:     u.Role(),
	}, nil
}
