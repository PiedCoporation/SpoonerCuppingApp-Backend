package usecases

import (
	"backend/global"
	"backend/internal/constants/enums/circlestyle"
	"backend/internal/constants/enums/jwtpurpose"
	"backend/internal/constants/enums/rolename"
	"backend/internal/constants/errorcode"
	"backend/internal/contracts/common"
	"backend/internal/contracts/user"
	"backend/internal/domains/commons"
	"backend/internal/domains/entities"
	"backend/internal/infrastructures/cache/rolecache"
	"backend/internal/mapper"
	repoAbstractions "backend/internal/persistents/abstractions"
	serviceAbstractions "backend/internal/usecases/abstractions"
	"backend/pkg/utils/jwt"
	"backend/pkg/utils/password"
	"backend/pkg/utils/sendto"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type userAuthService struct {
	uow      repoAbstractions.UserAuthUow
	userRepo repoAbstractions.IUserRepository
	rtRepo   repoAbstractions.IRefreshTokenRepository
}

func NewUserAuthService(
	uow repoAbstractions.UserAuthUow,
	userRepo repoAbstractions.IUserRepository,
	rtRepo repoAbstractions.IRefreshTokenRepository,
) serviceAbstractions.IUserAuthService {
	return &userAuthService{
		uow:      uow,
		userRepo: userRepo,
		rtRepo:   rtRepo,
	}
}

// Register implements user.IUserAuthService.
func (us *userAuthService) Register(ctx context.Context, vo user.RegisterUserVO) (*common.Result[user.UserRes]) {
	g, gCtx := errgroup.WithContext(ctx)

	fmt.Println("Register user", vo)

	var hashedPassword string

    // if email already exists: for unverified users, update info and resend verify; for verified users, block
    if existing, err := us.userRepo.GetByEmail(ctx, vo.Email); err == nil {
        if existing.IsVerified {
            return common.Failure[user.UserRes](&common.Error{Code: 400, Message: errorcode.ErrEmailExists.Error()})
        }

        // update latest submitted info
        newHashed, hashErr := password.HashPassword(vo.Password)
        if hashErr != nil {
            return common.Failure[user.UserRes](&common.Error{Code: 500, Message: hashErr.Error()})
        }
        if updErr := us.userRepo.Update(ctx, existing.ID, map[string]any{
            "first_name": vo.FirstName,
            "last_name":  vo.LastName,
            "password":   newHashed,
        }); updErr != nil {
            return common.Failure[user.UserRes](&common.Error{Code: 500, Message: updErr.Error()})
        }
        existing.FirstName = vo.FirstName
        existing.LastName = vo.LastName
        existing.Password = newHashed

        // resend verify for existing but unverified account
        token, genErr := jwt.GenerateEmailToken([]byte(global.Config.JWT.RegisterTokenKey),
            global.Config.JWT.RegisterTokenExpiresIn, existing.ID, jwtpurpose.Register)
			
        if genErr != nil {
            return common.Failure[user.UserRes](&common.Error{Code: 500, Message: "Failed to generate email verify jwt"})
        }

        verifyLink := fmt.Sprintf("%s/v1/users/register/verify?token=%s",
            global.Config.HTTP.Url, token)

        go func(email, vlink string) {
            if err := sendto.SendTemplateEmailOtp(&global.Config.SMTP, []string{email},
                "register-verify.html", map[string]any{"verifyLink": vlink},
            ); err != nil {
                fmt.Println(err)
            }
        }(vo.Email, verifyLink)

        mapped := mapper.MapUserToContractUserLoginResponse(existing)
        return common.Success(mapped)
    } else if !errors.Is(err, errorcode.ErrUserNotFound) {
        return common.Failure[user.UserRes](&common.Error{Code: 500, Message: err.Error()})
    }

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

	// hash pass
	g.Go(func() error {
		hp, err := password.HashPassword(vo.Password)
		if err != nil {
			return err
		}
		hashedPassword = hp
		return nil
	})

    if err := g.Wait(); err != nil {
        return common.Failure[user.UserRes](&common.Error{Code: 400, Message: err.Error()})
    }

    // get user role id
	defaultRole, ok := rolecache.Get(string(rolename.User))
    if !ok {
        return common.Failure[user.UserRes](&common.Error{Code: 500, Message: "Failed to get default role"})
    }

	// create user
	now := time.Now()
    createdUser := &entities.User{
		Entity:      commons.Entity{ID: uuid.New(), IsDeleted: false},
		FirstName:   vo.FirstName,
		LastName:    vo.LastName,
		Email:       vo.Email,
		Phone:       "",
		Password:    hashedPassword,
		CircleStyle: circlestyle.Minimal,
		IsVerified:  false,
		Auditable:   commons.Auditable{CreatedAt: now, UpdatedAt: now},
		RoleID:      defaultRole.ID,
	}

	// insert user into db
    if err := us.userRepo.Create(ctx, createdUser); err != nil {
        return common.Failure[user.UserRes](&common.Error{Code: 500, Message: "Failed to create user"})
    }

	// gene email verify jwt
    token, err := jwt.GenerateEmailToken([]byte(global.Config.JWT.RegisterTokenKey),
        global.Config.JWT.RegisterTokenExpiresIn, createdUser.ID, jwtpurpose.Register)
    if err != nil {
        return common.Failure[user.UserRes](&common.Error{Code: 500, Message: "Failed to generate email verify jwt"})
    }

	verifyLink := fmt.Sprintf("%s/v1/users/register/verify?token=%s",
		global.Config.HTTP.Url, token)

	// send verify email to activate account
    go func(email, vlink string) {
        if err := sendto.SendTemplateEmailOtp(&global.Config.SMTP, []string{email},
            "register-verify.html", map[string]any{"verifyLink": vlink},
        ); err != nil {
            fmt.Println(err)
        }
    }(vo.Email, verifyLink)

    // return user info for UI (no tokens on register)
    mapped := mapper.MapUserToContractUserLoginResponse(createdUser)
    return common.Success(mapped)
}

// ResendEmailVerifyRegister implements user.IUserAuthService.
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
	token, err := jwt.GenerateEmailToken([]byte(global.Config.JWT.RegisterTokenKey),
		global.Config.JWT.RegisterTokenExpiresIn, user.ID, jwtpurpose.Register)
	if err != nil {
		return err
	}

	verifyLink := fmt.Sprintf("%s/v1/users/register/verify?token=%s",
		global.Config.HTTP.Url, token)

	// send verify email to activate account
	if err := sendto.SendTemplateEmailOtp(&global.Config.SMTP, []string{email},
		"register-verify.html", map[string]any{"verifyLink": verifyLink},
	); err != nil {
		return err
	}

	return nil
}

// VerifyRegister implements user.IUserAuthService.
func (us *userAuthService) VerifyRegister(ctx context.Context, userID uuid.UUID) (string, string, error) {
	// get user from db
	user, err := us.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, errorcode.ErrNotFound) {
			err = errorcode.ErrUserNotFound
		}
		return "", "", err
	}

	// re-check if user is verified or not
	if user.IsVerified {
		return "", "", errorcode.ErrAccountIsVerified
	}

	// begin transaction
	rp, err := us.uow.Begin(ctx)
	if err != nil {
		return "", "", err
	}
	defer us.uow.Rollback()

	// exec in db
	if err := rp.UserRepository().Update(ctx, userID, map[string]any{
		"is_verified": true,
	}); err != nil {
		return "", "", err
	}

	// gene ac and rt
	accessToken, refreshToken, err := jwt.GenerateAcAndRtTokens(user.ID)
	if err != nil {
		return "", "", err
	}

	// insert rt to db
	if err := insertRefreshToken(ctx, userID,
		rp.RefreshTokenRepository(),
		refreshToken, []byte(global.Config.JWT.RefreshTokenKey)); err != nil {
		return "", "", err
	}

	// commit transaction
	if err := us.uow.Commit(); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// Login implements user.IUserAuthService.
func (us *userAuthService) Login(ctx context.Context, vo user.LoginUserReq) (*common.Result[user.LoginUserRes]) {
	// get user from db
	dbUser, err := us.userRepo.GetByEmail(ctx, vo.Email)
	if err != nil {
		return common.Failure[user.LoginUserRes](&common.Error{Code: 404, Message: "User not found"})
	}

	// check if user is verified or not
	if !dbUser.IsVerified {
		return common.Failure[user.LoginUserRes](&common.Error{Code: 403, Message: "Account is not verified, please check your email for verification"})
	}
	// check if user is deleted or not
	if dbUser.IsDeleted {
		return common.Failure[user.LoginUserRes](&common.Error{Code: 403, Message: "Account is deleted"})
	}

	// check password
	if !password.ComparePasswords(dbUser.Password, vo.Password) {
		return common.Failure[user.LoginUserRes](&common.Error{Code: 401, Message: "Invalid password"})
	}

	// gene ac and rt
	accessToken, refreshToken, err := jwt.GenerateAcAndRtTokens(dbUser.ID)
	if err != nil {
		return common.Failure[user.LoginUserRes](&common.Error{Code: 500, Message: "Generate token error"})
	}

	if err := insertRefreshToken(ctx, dbUser.ID,
		us.rtRepo, refreshToken, []byte(global.Config.JWT.RefreshTokenKey)); err != nil {
		return common.Failure[user.LoginUserRes](&common.Error{Code: 500, Message: "Insert refresh token error"})
	}

	userRes := mapper.MapUserToContractUserLoginResponse(dbUser)

	return common.Success(&user.LoginUserRes{AccessToken: accessToken, RefreshToken: refreshToken, User: *userRes})
}

// Logout implements user.IUserAuthService.
func (us *userAuthService) Logout(ctx context.Context, userID uuid.UUID, refreshToken string) error {
	// decode rt
	claims, err := jwt.ValidateToken([]byte(global.Config.JWT.RefreshTokenKey),
		refreshToken, jwtpurpose.Refresh)
	if err != nil {
		return err
	}

	// compare userID from ac and rt
	if claims.Subject != userID.String() {
		return errorcode.ErrInvalidToken
	}

	// check if revoked or not
	rt, err := us.rtRepo.GetByTokenAndUserID(ctx, refreshToken, userID)
	if err != nil {
		return errorcode.ErrInvalidToken
	}

	// revoke
	if err = us.rtRepo.Update(ctx, rt.ID, map[string]any{
		"revoked": true,
	}); err != nil {
		return err
	}

	return nil
}

// ForgotPassword implements abstractions.IUserAuthService.
func (us *userAuthService) ForgotPassword(ctx context.Context, email string) *common.Result[string] {
	// get user from db
	user, err := us.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return common.Failure[string](&common.Error{Code: 500, Message: "Get user error"})
	}

	// check if user is verified or not
	if !user.IsVerified {
		return common.Failure[string](&common.Error{Code: 403, Message: "Account is not verified, please check your email for verification"})
	}
	// check if user is deleted or not
	if user.IsDeleted {
		return common.Failure[string](&common.Error{Code: 403, Message: "Account is deleted"})
	}

	// generate a secure 6-digit numeric code
	max := big.NewInt(1000000)
	num, err := rand.Int(rand.Reader, max)
	if err != nil {
		return common.Failure[string](&common.Error{Code: 500, Message: "Generate code error"})
	}
	code := fmt.Sprintf("%06d", num.Int64())

	// set expiry (e.g., 10 minutes)
	minutes := 3
	expiresAt := time.Now().Add(time.Duration(minutes) * time.Minute)

	// persist code and expiry
	if err := us.userRepo.Update(ctx, user.ID, map[string]any{
		"forgot_password_code":        code,
		"forgot_password_expires_at": expiresAt,
	}); err != nil {
		return common.Failure[string](&common.Error{Code: 500, Message: "Update user error"})
	}

	// send email with code
	go func() {
		if err := sendto.SendTemplateEmailOtp(&global.Config.SMTP, []string{email},
			"forgot-password-verify-code.html", map[string]any{"code": code, "minutes": minutes},
		); err != nil {
			fmt.Println(err)
		}
	}()

	message := "Create code successfully"
	return common.Success(&message)
}

// ChangePassword implements abstractions.IUserAuthService.
func (us *userAuthService) ChangePassword(ctx context.Context, vo user.ChangePasswordVO) error {
	g, gCtx := errgroup.WithContext(ctx)

	var user *entities.User
	var hashedPassword string

	// get user from db
	g.Go(func() error {
		us, err := us.userRepo.GetByID(gCtx, vo.UserID)
		if err != nil {
			if errors.Is(err, errorcode.ErrNotFound) {
				return errorcode.ErrUserNotFound
			}
			return err
		}
		user = us
		return nil
	})

	// hash pass
	g.Go(func() error {
		hp, err := password.HashPassword(vo.Password)
		if err != nil {
			return err
		}
		hashedPassword = hp
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	// re-check if user is verified or not
	if !user.IsVerified {
		return errorcode.ErrAccountIsNotVerified
	}
	// re-check if user is deleted or not
	if user.IsDeleted {
		return errorcode.ErrAccountIsDeleted
	}

	// begin transaction
	rp, err := us.uow.Begin(ctx)
	if err != nil {
		return err
	}
	defer us.uow.Rollback()

	// update user password in db
	if err := rp.UserRepository().Update(ctx, user.ID, map[string]any{
		"password": hashedPassword,
	}); err != nil {
		return err
	}

	// revoke all old refresh token
	if err := rp.RefreshTokenRepository().RevokeAllByUserID(ctx, user.ID); err != nil {
		return err
	}

	// commit transaction
	if err := us.uow.Commit(); err != nil {
		return err
	}

	return nil
}

// VerifyForgotPasswordCode implements abstractions.IUserAuthService.
func (us *userAuthService) VerifyForgotPasswordCode(ctx context.Context, code string, email string) *common.Result[string] {
	// get user from db
	user, err := us.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return common.Failure[string](&common.Error{Code: 500, Message: "Get user error"})
	}

	// check if user is verified or not
	if !user.IsVerified {
		return common.Failure[string](&common.Error{Code: 403, Message: "Account is not verified, please check your email for verification"})
	}
	// check if user is deleted or not
	if user.IsDeleted {
		return common.Failure[string](&common.Error{Code: 403, Message: "Account is deleted"})
	}

	// check if code is expired
	if user.ForgotPasswordExpiresAt.Before(time.Now()) {
		return common.Failure[string](&common.Error{Code: 403, Message: "Code is expired"})
	}

	// check if code is correct
	if user.ForgotPasswordCode != code {
		return common.Failure[string](&common.Error{Code: 403, Message: "Code is incorrect"})
	}

	// update user forgot password code and expires at
	if err := us.userRepo.Update(ctx, user.ID, map[string]any{
		"forgot_password_code":        "",
		"forgot_password_expires_at": nil,
	}); err != nil {
		return common.Failure[string](&common.Error{Code: 500, Message: "Update user error"})
	}

	accessToken, _, err := jwt.GenerateAcAndRtTokens(user.ID)
	if err != nil {
		return common.Failure[string](&common.Error{Code: 500, Message: "Generate token error"})
	}

	return common.Success(&accessToken)
}

// RefreshToken implements user.IUserAuthService.
func (us *userAuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	// validate token
	claims, err := jwt.ValidateToken([]byte(global.Config.JWT.RefreshTokenKey),
		refreshToken, jwtpurpose.Refresh)
	if err != nil {
		return "", "", errorcode.ErrInvalidToken
	}

	// parse sub to uuid
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return "", "", errorcode.ErrInvalidToken
	}

	oldRefreshToken, err := us.rtRepo.GetByTokenAndUserID(ctx, refreshToken, userID)
	if err != nil {
		return "", "", errorcode.ErrInvalidToken
	}

	// gene ac and rt
	accessToken, newRefreshToken, err := jwt.GenerateAcAndRtTokens(userID)
	if err != nil {
		return "", "", err
	}

	// begin transaction
	rp, err := us.uow.Begin(ctx)
	if err != nil {
		return "", "", err
	}
	defer us.uow.Rollback()

	// insert new rt to db
	if err := insertRefreshToken(ctx, userID,
		rp.RefreshTokenRepository(), newRefreshToken, []byte(global.Config.JWT.RefreshTokenKey)); err != nil {
		return "", "", err
	}

	// revoke old rt
	if err := rp.RefreshTokenRepository().Update(ctx, oldRefreshToken.ID, map[string]any{
		"revoked": true,
	}); err != nil {
		return "", "", err
	}

	// commit
	if err := us.uow.Commit(); err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

// helper

// insertRefreshToken
func insertRefreshToken(
	ctx context.Context, userID uuid.UUID,
	rtRepo repoAbstractions.IRefreshTokenRepository,
	refreshToken string, rtSecret []byte,
) error {
	claims, err := jwt.ValidateToken(rtSecret,
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

	if err := rtRepo.Create(ctx, rt); err != nil {
		return err
	}
	return nil
}
