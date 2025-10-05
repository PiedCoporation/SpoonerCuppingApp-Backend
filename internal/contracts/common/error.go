package common

type Error struct {
	Code string `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface for common.Error.
func (e *Error) Error() string {
    if e == nil {
        return ""
    }
    return e.Message
}