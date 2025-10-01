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
