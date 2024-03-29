package constants

import "time"

const (
	KrosmozAlmanaxDateFormat = "2006-01-02"
)

func GetAlmanaxFirstDate() time.Time {
	return time.Date(2012, 9, 18, 0, 0, 0, 0, time.UTC)
}

func GetAlmanaxLastDate() time.Time {
	return time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC)
}
