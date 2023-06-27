package zkerrors

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ZkErrorParamToMessage map[string]string

func MessageFromValidation(fe validator.FieldError) string {
	valueString := fmt.Sprintf("%v", fe.Value())
	switch fe.Tag() {
	case "required":
		return "The field '" + fe.Field() + "' is required"
	case "email":
		return "Invalid email address provided " + valueString
	}
	return fe.Error() // default error
}
