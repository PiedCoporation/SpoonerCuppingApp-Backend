package user

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
