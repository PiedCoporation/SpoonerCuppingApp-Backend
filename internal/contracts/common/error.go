package common

type Error struct {
	Code int    `json:"-"`
	Message string `json:"message"`
}

// Error implements the error interface for common.Error.
func (e *Error) Error() string {
    if e == nil {
        return ""
    }
    return e.Message
}