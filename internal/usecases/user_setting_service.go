package usecases

import (
	"backend/internal/constants/errorcode"
	"backend/internal/contracts/common"
	"backend/internal/contracts/user"
	"backend/internal/mapper"
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

// UpdateUser implements abstractions.IUserSettingService.
func (s userSettingService) Update(ctx context.Context, userID uuid.UUID, req user.UpdateUserReq) (*common.Result[user.UserRes]) {
	// get user from db
	userEntity, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, errorcode.ErrNotFound) {
			return common.Failure[user.UserRes](&common.Error{Code: "404", Message: "User not found"})
		}
		return common.Failure[user.UserRes](&common.Error{Code: "500", Message: "Failed to get user"})
	}

	if req.CircleStyle != nil && !req.CircleStyle.IsValid() {
		return common.Failure[user.UserRes](&common.Error{Code: "400", Message: "Invalid circle style"})
	}

	fieldMap := make(map[string]any)
	if req.CircleStyle != nil {
		fieldMap["circle_style"] = *req.CircleStyle
		userEntity.CircleStyle = *req.CircleStyle
	}

	if req.FirstName != nil {
		fieldMap["first_name"] = *req.FirstName
		userEntity.FirstName = *req.FirstName
	}

	if req.LastName != nil {
		fieldMap["last_name"] = *req.LastName
		userEntity.LastName = *req.LastName
	}

	if req.Phone != nil {
		fieldMap["phone"] = *req.Phone
		userEntity.Phone = *req.Phone
	}
	
	if len(fieldMap) == 0 {
		userRes := mapper.MapUserToContractUserLoginResponse(userEntity)
		return common.Success(userRes)
	}

	if err := s.userRepo.Update(ctx, userID, fieldMap); err != nil {
		return common.Failure[user.UserRes](&common.Error{Code: "500", Message: "Failed to update user"})
	}

	userRes := mapper.MapUserToContractUserLoginResponse(userEntity)
	return common.Success(userRes)
}
