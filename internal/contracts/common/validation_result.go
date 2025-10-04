package common

type ValidationResult[T any] struct {
	Result *Result[T] `json:"result"`
	Errors []Error `json:"errors"`
}

func CreateValidationResult[T any](errors []Error) *ValidationResult[T] {
	r := &ValidationResult[T]{}
	r.Result = Failure[T](&Error{Code: "VALIDATION_ERROR", Message: "Validation error"})
	r.Errors = errors
	return r
}

func AddError[T any](r *ValidationResult[T], error Error) *ValidationResult[T] {
	r.Result = Failure[T](&Error{Code: "VALIDATION_ERROR", Message: "Validation error"})
	r.Errors = append(r.Errors, error)
	return r
}