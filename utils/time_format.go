package utils

import (
	"fmt"
	"math/bits"
)

func TimeFormat(seconds uint64) (timeStr string) {
	hours, remainder := bits.Div64(0, seconds, 3600)
	minutes, seconds := bits.Div64(0, remainder, 60)
	days, hours := bits.Div64(0, hours, 24)
	timeStr = ""
	if days > 0 {
		if days == 1 {
			timeStr += fmt.Sprintf("%d day, ", days)
		} else {
			timeStr += fmt.Sprintf("%d days, ", days)
		}
	}
	if hours > 0 {
		if hours == 1 {
			timeStr += fmt.Sprintf("%d hour, ", hours)
		} else {
			timeStr += fmt.Sprintf("%d hours, ", hours)
		}
	}
	if minutes > 0 {
		if minutes == 1 {
			timeStr += fmt.Sprintf("%d minute, ", minutes)
		} else {
			timeStr += fmt.Sprintf("%d minutes, ", minutes)
		}
	}
	if seconds > 0 {
		if seconds == 1 {
			timeStr += fmt.Sprintf("%d second", seconds)
		} else {
			timeStr += fmt.Sprintf("%d seconds", seconds)
		}
	}
	return timeStr
}
