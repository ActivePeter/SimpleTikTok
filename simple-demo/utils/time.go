package utils

import "time"

func TimeStamp() int64 {
	return time.Now().Unix()
}
