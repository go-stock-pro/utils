package typecasting

import (
	"strconv"

	"github.com/pkg/errors"
)

func ConvertStringToFloat32(in string, fieldName string) (float32, error) {
	out, err := strconv.ParseFloat(in, 32)
	if err != nil {
		return 0.0, errors.Errorf("%v, at field: %s, value: %s", err, fieldName, in)
	}

	return float32(out), nil
}

func ConvertStringToFloat64(in string, fieldName string) (float64, error) {
	out, err := strconv.ParseFloat(in, 64)
	if err != nil {
		return 0.0, errors.Errorf("%v, at field: %s, value: %s", err, fieldName, in)
	}

	return out, nil
}
