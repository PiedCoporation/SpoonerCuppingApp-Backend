package validation

import (
	"backend/pkg/utils/stringutils"
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func TranslateValidationError(err error) string {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		msgs := []string{}
		for _, fe := range ve {
			var msg string
			fieldName := stringutils.ToSnakeCase(fe.Field())
			paramName := stringutils.ToSnakeCase(fe.Param())
			switch fe.Tag() {
			case "required":
				msg = fmt.Sprintf("%s is required", fieldName)
			case "email":
				msg = fmt.Sprintf("%s must be an email", fieldName)
			case "eqfield":
				msg = fmt.Sprintf("%s must be equal to %s", fieldName, paramName)
			case "len":
				msg = fmt.Sprintf("%s must be exactly %s characters long", fieldName, paramName)
			case "min":
				msg = fmt.Sprintf("%s must be at least %s characters", fieldName, paramName)
			case "max":
				msg = fmt.Sprintf("%s must be at most %s characters", fieldName, paramName)
			default:
				msg = fmt.Sprintf("%s: %s", fieldName, fe.Tag())
			}
			msgs = append(msgs, msg)
		}
		return strings.Join(msgs, "; ")
	}
	return err.Error()
}
