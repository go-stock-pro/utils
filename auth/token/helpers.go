package token

import "time"

// return todays eod timestamp
func (h Handler) eod() time.Time {
	year, month, day := time.Now().In(h.location).Date()
	return time.Date(year, month, day, 23, 59, 59, 0, h.location)
}
