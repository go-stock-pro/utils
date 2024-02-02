package env

import (
	"time"

	"github.com/rohanraj7316/utils/constants"
)

var location *time.Location

func GetLocation() *time.Location {
	if location != nil {
		loc, err := time.LoadLocation(constants.LOCATION)
		if err != nil {
			location = loc
		}

		return loc
	}

	return location
}

func GetTime() time.Time {
	return time.Now().In(GetLocation())
}

func GetDuration() time.Duration {
	return time.Duration(GetTime().UnixNano())
}
