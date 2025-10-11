package mapper

import (
	userContract "backend/internal/contracts/user"
	"backend/internal/domains/entities"
)

func MapUserToContractUserResponse(u *entities.User) *userContract.UserViewRes {
	return &userContract.UserViewRes{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}

func MapUserToContractUserLoginResponse(u *entities.User) *userContract.UserRes {
	return &userContract.UserRes{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Phone:     u.Phone,
		Role:      u.Role.Name,
		CircleStyle: u.CircleStyle,
	}
}