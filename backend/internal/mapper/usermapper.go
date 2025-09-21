package mapper

import (
	userContract "backend/internal/contracts/user"
	"backend/internal/domains/entities"
)

func MapUserToContractUserResponse(u *entities.User) *userContract.UserReponse {
	return &userContract.UserReponse{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Phone:     u.Phone,
	}
}
