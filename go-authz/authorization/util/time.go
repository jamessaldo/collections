package util

import "time"

func GetTimestampUTC() time.Time {
	return time.Now().UTC()
}
