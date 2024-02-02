package typecasting

import (
	"strings"

	"github.com/pkg/errors"
)

func ConvertStringToBool(in string, fieldName string) (bool, error) {
	switch strings.ToLower(in) {
	case "true", "yes":
		return true, nil
	case "false", "no":
		return false, nil
	default:
		return false, errors.Errorf("invalid input string conversion to bool, at field: %s, value: %s", fieldName, in)
	}
}
