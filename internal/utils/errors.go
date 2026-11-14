package utils

import (
	"fmt"
	"time"

	"github.com/gotd/td/tgerr"
)

// DescribeFloodWait turns a FLOOD_WAIT error into a message a user can act on.
// It returns false when err is not a flood wait.
func DescribeFloodWait(err error) (string, bool) {
	if err == nil {
		return "", false
	}
	wait, ok := tgerr.AsFloodWait(err)
	if !ok {
		return "", false
	}
	if wait <= 0 {
		wait = time.Second
	}
	return fmt.Sprintf(
		"Telegram is rate limiting this bot, please try again in %s.",
		wait.Round(time.Second),
	), true
}
