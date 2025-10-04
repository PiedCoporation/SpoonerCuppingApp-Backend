package common

type Result[T any] struct {
	Data *T `json:"data"`
	IsSuccess bool `json:"is_success"`
	IsFailure bool `json:"is_failure"`
	Error *Error `json:"error"`
}

func Success[T any](data *T) *Result[T] {
	r := &Result[T]{}
	r.Data = data
	r.IsSuccess = true
	r.IsFailure = false
	r.Error = nil
	return r		
}

func Failure[T any](error *Error) *Result[T] {
	r := &Result[T]{}
	r.Data = nil
	r.IsSuccess = false
	r.IsFailure = true
	r.Error = error
	return r
}