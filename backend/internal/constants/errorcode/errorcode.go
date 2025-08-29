package errorcode

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	// 400
	ErrEmailExists          = errors.New("this email already exists")
	ErrPhoneExists          = errors.New("this phone number already exists")
	ErrAccountIsVerified    = errors.New("this account is already verified")
	ErrAccountIsNotVerified = errors.New("this account is not verified")

	// 401
	ErrInvalidToken      = errors.New("invalid token")
	ErrInvalidJWTPurpose = errors.New("invalid jwt purpose")

	// 403
	ErrInactiveAccount = errors.New("this account is inactive")
	ErrDeletedAccount  = errors.New("this account is deleted")

	// 404
	ErrUserNotFound = errors.New("user not found")

	// 409
	ErrEmailBelongsToDeletedAccount = errors.New("email belongs to deleted account")
	ErrPhoneBelongsToDeletedAccount = errors.New("phone number belongs to deleted account")

	// 500
	ErrUnexpectedSigningToken = errors.New("unexpected signing token")
	ErrUnexpectedCreatingUser = errors.New("unexpected creating user")
)

// Map code -> http code
var errorStatusMap = map[error]int{
	// 400
	ErrEmailExists:       http.StatusBadRequest,
	ErrPhoneExists:       http.StatusBadRequest,
	ErrAccountIsVerified: http.StatusBadRequest,

	// 401
	ErrInvalidToken:         http.StatusUnauthorized,
	ErrInvalidJWTPurpose:    http.StatusUnauthorized,
	ErrAccountIsNotVerified: http.StatusUnauthorized,

	// 403
	ErrInactiveAccount: http.StatusForbidden,
	ErrDeletedAccount:  http.StatusForbidden,

	// 404
	ErrUserNotFound: http.StatusNotFound,

	// 409
	ErrEmailBelongsToDeletedAccount: http.StatusConflict,
	ErrPhoneBelongsToDeletedAccount: http.StatusConflict,

	// 500
	ErrUnexpectedSigningToken: http.StatusInternalServerError,
	ErrUnexpectedCreatingUser: http.StatusInternalServerError,
}

// utils write error
func JSONError(c *gin.Context, err error) {
	status, ok := errorStatusMap[err]
	if !ok {
		status = http.StatusInternalServerError
	}
	c.JSON(status, gin.H{
		"error": err.Error(),
	})
}
