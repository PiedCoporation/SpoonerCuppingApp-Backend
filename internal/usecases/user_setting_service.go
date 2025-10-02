package usecases

import (
	"backend/internal/constants/errorcode"
	"backend/internal/contracts/user"
	repoAbstractions "backend/internal/persistents/abstractions"
	serviceAbstractions "backend/internal/usecases/abstractions"
	"context"
	"errors"

	"github.com/google/uuid"
)

type userSettingService struct {
	userRepo repoAbstractions.IUserRepository
}

func NewUserSettingService(
	userRepo repoAbstractions.IUserRepository,
) serviceAbstractions.IUserSettingService {
	return userSettingService{
		userRepo: userRepo,
	}
}

// GetUserSetting implements abstractions.IUserSettingService.
func (s userSettingService) GetByID(ctx context.Context, userID uuid.UUID) (*user.UserSettingRes, error) {
	// get user from db
	userEntity, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, errorcode.ErrNotFound) {
			return nil, errorcode.ErrUserNotFound
		}
		return nil, err
	}

	return &user.UserSettingRes{
		ID:          userEntity.ID,
		CircleStyle: userEntity.CircleStyle,
	}, nil
}

// UpdateUserSetting implements abstractions.IUserSettingService.
func (s userSettingService) Update(ctx context.Context, userID uuid.UUID, req user.UpdateUserSettingReq) error {
	// get user from db
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, errorcode.ErrNotFound) {
			return errorcode.ErrUserNotFound
		}
		return err
	}

	if req.CircleStyle != nil && !req.CircleStyle.IsValid() {
		return errorcode.ErrInvalidCircleStyle
	}

	fieldMap := make(map[string]any)
	if req.CircleStyle != nil && *req.CircleStyle != user.CircleStyle {
		fieldMap["circle_style"] = *req.CircleStyle
	}

	// update
	if len(fieldMap) == 0 {
		return nil
	}
	if err := s.userRepo.Update(ctx, userID, fieldMap); err != nil {
		return err
	}

	return nil
}
