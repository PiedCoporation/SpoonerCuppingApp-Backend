package errorcode

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	// 400
	ErrEmailExists                = errors.New("this email already exists")
	ErrPhoneExists                = errors.New("this phone number already exists")
	ErrAccountIsVerified          = errors.New("this account is already verified")
	ErrEventIsNotStartForRegister = errors.New("this event is not start for register")
	ErrEventIsFull                = errors.New("this event is full")
	ErrUserAlreadyRegistered      = errors.New("this user is already registered for this event")
	ErrEventSamplesRequired       = errors.New("event samples are required")
	ErrEventAddressRequired       = errors.New("event address is required")

	// 401
	ErrInvalidToken      = errors.New("invalid token")
	ErrInvalidJWTPurpose = errors.New("invalid jwt purpose")
	ErrInvalidPassword   = errors.New("invalid password")

	// 403
	ErrAccountIsNotVerified = errors.New("this account is not verified")
	ErrAccountIsDeleted     = errors.New("this account is deleted")

	// 404
	ErrNotFound      = errors.New("not found")
	ErrUserNotFound  = errors.New("user not found")
	ErrEventNotFound = errors.New("event not found")

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
	ErrEmailExists:                http.StatusBadRequest,
	ErrPhoneExists:                http.StatusBadRequest,
	ErrAccountIsVerified:          http.StatusBadRequest,
	ErrEventIsNotStartForRegister: http.StatusBadRequest,
	ErrEventIsFull:                http.StatusBadRequest,
	ErrUserAlreadyRegistered:      http.StatusBadRequest,
	ErrEventSamplesRequired:       http.StatusBadRequest,
	ErrEventAddressRequired:       http.StatusBadRequest,

	// 401
	ErrInvalidToken:      http.StatusUnauthorized,
	ErrInvalidJWTPurpose: http.StatusUnauthorized,
	ErrInvalidPassword:   http.StatusUnauthorized,

	// 403
	ErrAccountIsNotVerified: http.StatusForbidden,
	ErrAccountIsDeleted:     http.StatusForbidden,

	// 404
	ErrNotFound:      http.StatusNotFound,
	ErrUserNotFound:  http.StatusNotFound,
	ErrEventNotFound: http.StatusNotFound,

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
