package constants

import (
	"time"

	"github.com/aptible/supercronic/cronexpr"
)

type EXPIRE_AT string

const (
	TIME_FORMAT = "2006-01-02T15:04:05.000Z"

	LOCATION = "Asia/Kolkata"
)

var (
	CONSTANTS_ORDER_LIMIT_EXPIRE_TIME  EXPIRE_AT = "59 23 * * *"
	CONSTANTS_ORDER_MARKET_EXPIRE_TIME EXPIRE_AT = "59 23 * * *"
)

func ExpireAt(et EXPIRE_AT) time.Time {
	return cronexpr.MustParse(string(et)).Next(time.Now())
}
