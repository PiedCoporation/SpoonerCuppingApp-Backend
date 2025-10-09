package abstractions

import (
	"backend/internal/contracts/common"
	"backend/internal/contracts/user"
	"context"

	"github.com/google/uuid"
)

type (
	IUserAuthService interface {
		Register(ctx context.Context, vo user.RegisterUserVO) error
		ResendEmailVerifyRegister(ctx context.Context, email string) error
		VerifyRegister(ctx context.Context, userID uuid.UUID) (string, string, error)
		Login(ctx context.Context, vo user.LoginUserReq) (*common.Result[user.LoginUserRes])
		Logout(ctx context.Context, userID uuid.UUID, refreshToken string) error
		ForgotPassword(ctx context.Context, email string) error
		ChangePassword(ctx context.Context, vo user.ChangePasswordVO) error
		RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	}

	IUserSettingService interface {
		GetByID(ctx context.Context, userID uuid.UUID) (*user.UserSettingRes, error)
		Update(ctx context.Context, userID uuid.UUID, req user.UpdateUserSettingReq) error
	}
)
