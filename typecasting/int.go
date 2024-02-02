package typecasting

import (
	"strconv"

	"github.com/pkg/errors"
)

func ConvertStringToInt(in string, fieldName string) (int, error) {
	out, err := strconv.Atoi(in)
	if err != nil {
		return 0, errors.Errorf("%v, at field: %s, value: %s", err, fieldName, in)
	}

	return out, nil
}

func ConvertStringToInt64(in string, fieldName string) (int64, error) {
	out, err := strconv.ParseInt(in, 10, 64)
	if err != nil {
		return 0, errors.Errorf("%v, at field: %s, value: %s", err, fieldName, in)
	}

	return out, nil
}
