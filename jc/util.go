package jc

import (
	"time"
)

// Get proper time format for JumpCloud
func getTime() string {
	location, _ := time.LoadLocation("Etc/GMT")
	time := time.Now().In(location).Format(time.RFC1123)
	return time
}
