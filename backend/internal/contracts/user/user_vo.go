package user

import "github.com/google/uuid"

type RegisterUserVO struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Password  string
}

type LoginUserVO struct {
	Email    string
	Password string
}

type ChangePasswordVO struct {
	UserID   uuid.UUID
	Password string
}
