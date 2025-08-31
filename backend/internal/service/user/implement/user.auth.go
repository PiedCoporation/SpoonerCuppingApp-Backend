package implement

import (
	"backend/config"
	"backend/internal/cache/rolecache"
	"backend/internal/constants/enum/jwtpurpose"
	"backend/internal/constants/enum/rolename"
	"backend/internal/constants/errorcode"
	"backend/internal/domain/commons"
	"backend/internal/domain/entities"
	"backend/internal/repository"
	"backend/internal/service/user"
	"backend/pkg/utils/jwt"
	"backend/pkg/utils/sendto"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type userAuthService struct {
	config   *config.Config
	userRepo repository.UserRepository
	rtRepo   repository.RefreshTokenRepository
}

func NewUserAuthService(
	config *config.Config,
	userRepo repository.UserRepository,
	rtRepo repository.RefreshTokenRepository,
) user.UserAuthService {
	return &userAuthService{
		config:   config,
		userRepo: userRepo,
		rtRepo:   rtRepo,
	}
}

// Register implements user.UserAuthService.
func (us *userAuthService) Register(ctx context.Context, vo user.RegisterUserVO) error {
	g, gCtx := errgroup.WithContext(ctx)

	// check if email exists
	g.Go(func() error {
		exists, err := us.userRepo.IsEmailTaken(gCtx, vo.Email, uuid.Nil)
		if err != nil {
			return err
		}
		if exists {
			return errorcode.ErrEmailExists
		}
		return nil
	})

	// check if phone exists
	g.Go(func() error {
		exists, err := us.userRepo.IsPhoneTaken(gCtx, vo.Phone, uuid.Nil)
		if err != nil {
			return err
		}
		if exists {
			return errorcode.ErrPhoneExists
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	// get user role id
	defaultRole, ok := rolecache.Get(string(rolename.User))
	if !ok {
		return errorcode.ErrUnexpectedCreatingUser
	}

	// create user
	now := time.Now()
	user := &entities.User{
		Entity:     commons.Entity{ID: uuid.New(), IsDeleted: false},
		FirstName:  vo.FirstName,
		LastName:   vo.LastName,
		Email:      vo.Email,
		Phone:      vo.Phone,
		IsVerified: false,
		Auditable:  commons.Auditable{CreatedAt: now, UpdatedAt: now},
		RoleID:     defaultRole.ID,
	}

	// insert user into db
	if err := us.userRepo.Create(ctx, user); err != nil {
		return err
	}

	// gene email verify jwt
	token, err := jwt.GenerateEmailToken([]byte(us.config.JWT.RegisterTokenKey),
		us.config.JWT.RegisterTokenExpiresIn, user.ID, jwtpurpose.Register)
	if err != nil {
		return err
	}

	verifyLink := fmt.Sprintf("%s/v1/users/register/verify?token=%s",
		us.config.HTTP.Url, token)

	// send verify email to activate account
	if err := sendto.SendTemplateEmailOtp(&us.config.SMTP, []string{vo.Email},
		"register-verify.html", map[string]any{"verifyLink": verifyLink},
	); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// ResendEmailVerifyRegister implements user.UserAuthService.
func (us *userAuthService) ResendEmailVerifyRegister(ctx context.Context, email string) error {
	// get user from db
	user, err := us.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}

	// check if user is verified or not
	if user.IsVerified {
		return errorcode.ErrAccountIsVerified
	}

	// gene email verify jwt
	token, err := jwt.GenerateEmailToken([]byte(us.config.JWT.RegisterTokenKey),
		us.config.JWT.RegisterTokenExpiresIn, user.ID, jwtpurpose.Register)
	if err != nil {
		return err
	}

	verifyLink := fmt.Sprintf("%s/v1/users/register/verify?token=%s",
		us.config.HTTP.Url, token)

	// send verify email to activate account
	if err := sendto.SendTemplateEmailOtp(&us.config.SMTP, []string{email},
		"register-verify.html", map[string]any{"verifyLink": verifyLink},
	); err != nil {
		return err
	}

	return nil
}

// VerifyRegister implements user.UserAuthService.
func (us *userAuthService) VerifyRegister(ctx context.Context, userID uuid.UUID) (string, string, error) {
	// get user from db
	user, err := us.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", "", err
	}

	// re-check if user is verified or not
	if user.IsVerified {
		return "", "", errorcode.ErrAccountIsVerified
	}

	// exec in db
	if err := us.userRepo.UpdateEmailVerified(ctx, user.ID, true); err != nil {
		return "", "", err
	}

	// gene ac and rt
	accessToken, refreshToken, err := jwt.GenerateAcAndRtTokens(&us.config.JWT, user.ID)
	if err != nil {
		return "", "", err
	}

	// insert rt to db
	if err := us.insertRefreshToken(ctx, userID, refreshToken); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// Login implements user.UserAuthService.
func (us *userAuthService) Login(ctx context.Context, email string) error {
	// get user from db
	user, err := us.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}

	// check if user is verified or not
	if !user.IsVerified {
		return errorcode.ErrAccountIsNotVerified
	}

	// gene login verify jwt
	token, err := jwt.GenerateEmailToken([]byte(us.config.JWT.LoginTokenKey),
		us.config.JWT.LoginTokenExpiresIn, user.ID, jwtpurpose.Login)
	if err != nil {
		return err
	}

	verifyLink := fmt.Sprintf("%s/v1/users/login/verify?token=%s",
		us.config.HTTP.Url, token)

	// send verify link to login
	if err := sendto.SendTemplateEmailOtp(&us.config.SMTP, []string{email},
		"login-verify.html", map[string]any{"verifyLink": verifyLink},
	); err != nil {
		return err
	}

	return nil
}

// VerifyLogin implements user.UserAuthService.
func (us *userAuthService) VerifyLogin(ctx context.Context, userID uuid.UUID) (string, string, error) {
	// get user from db
	user, err := us.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", "", err
	}

	// re-check if user is verified or not
	if !user.IsVerified {
		return "", "", errorcode.ErrAccountIsNotVerified
	}

	// gene ac and rt
	accessToken, refreshToken, err := jwt.GenerateAcAndRtTokens(&us.config.JWT, user.ID)
	if err != nil {
		return "", "", err
	}

	if err := us.insertRefreshToken(ctx, userID, refreshToken); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// Logout implements user.UserAuthService.
func (us *userAuthService) Logout(ctx context.Context, userID uuid.UUID, refreshToken string) error {
	// decode rt
	claims, err := jwt.ValidateToken([]byte(us.config.JWT.RefreshTokenKey),
		refreshToken, jwtpurpose.Refresh)
	if err != nil {
		return err
	}

	// compare userID from ac and rt
	if claims.Subject != userID.String() {
		return errorcode.ErrInvalidToken
	}

	// check if revoked or not
	if _, err := us.rtRepo.GetByTokenAndUserID(ctx, refreshToken, userID); err != nil {
		return errorcode.ErrInvalidToken
	}

	// revoke
	if err = us.rtRepo.Revoke(ctx, refreshToken); err != nil {
		return err
	}

	return nil
}

// RefreshToken implements user.UserAuthService.
func (us *userAuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	// validate token
	claims, err := jwt.ValidateToken([]byte(us.config.JWT.RefreshTokenKey),
		refreshToken, jwtpurpose.Refresh)
	if err != nil {
		return "", "", errorcode.ErrInvalidToken
	}

	// parse sub to uuid
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return "", "", errorcode.ErrInvalidToken
	}

	if _, err := us.rtRepo.GetByTokenAndUserID(ctx, refreshToken, userID); err != nil {
		return "", "", errorcode.ErrInvalidToken
	}

	// gene ac and rt
	accessToken, newRefreshToken, err := jwt.GenerateAcAndRtTokens(&us.config.JWT, userID)
	if err != nil {
		return "", "", err
	}

	if err := us.insertRefreshToken(ctx, userID, newRefreshToken); err != nil {
		return "", "", err
	}

	// revoke old rt
	if err := us.rtRepo.Revoke(ctx, refreshToken); err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

// helper

// insertRefreshToken
func (us *userAuthService) insertRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string) error {
	claims, err := jwt.ValidateToken([]byte(us.config.JWT.RefreshTokenKey),
		refreshToken, jwtpurpose.Refresh)
	if err != nil {
		return err
	}

	now := time.Now()
	rt := &entities.RefreshToken{
		ID:        uuid.New(),
		Token:     refreshToken,
		IssuedAt:  claims.IssuedAt.Time,
		ExpiresAt: claims.ExpiresAt.Time,
		Revoked:   false,
		Auditable: commons.Auditable{CreatedAt: now, UpdatedAt: now},
		UserID:    userID,
	}

	if err := us.rtRepo.Create(ctx, rt); err != nil {
		return err
	}
	return nil
}
