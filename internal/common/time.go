package common

import "time"

func GetCurrentTime() time.Time {
	return time.Now().UTC().Truncate(time.Microsecond)
}
