package utils

import (
	"time"
)

const ENS_TIME_LAYOR = "2006-01-02 15:04:05"

func GetTime() string {
	return time.Now().Format(ENS_TIME_LAYOR)
}
