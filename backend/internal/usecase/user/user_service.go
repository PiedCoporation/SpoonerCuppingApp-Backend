package user

import (
	"context"

	"github.com/google/uuid"
)

type (
	UserAuthService interface {
		Register(ctx context.Context, vo RegisterUserVO) error
		ResendEmailVerifyRegister(ctx context.Context, email string) error
		VerifyRegister(ctx context.Context, userID uuid.UUID) (string, string, error)
		Login(ctx context.Context, email string) error
		VerifyLogin(ctx context.Context, userID uuid.UUID) (string, string, error)
		Logout(ctx context.Context, userID uuid.UUID, refreshToken string) error
		RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	}
)
