package typecasting

import (
	"time"
)

func ConvertStringToDuration(in string, fieldName string) (time.Duration, error) {
	dur, err := time.ParseDuration(in)
	if err != nil {
		return 0, err
	}

	return dur, nil
}
