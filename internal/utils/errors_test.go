package utils

import (
	"errors"
	"strings"
	"testing"

	"github.com/gotd/td/tgerr"
)

func TestDescribeFloodWait(t *testing.T) {
	if _, ok := DescribeFloodWait(nil); ok {
		t.Fatal("nil error described as flood wait")
	}
	if _, ok := DescribeFloodWait(errors.New("boom")); ok {
		t.Fatal("generic error described as flood wait")
	}
	hint, ok := DescribeFloodWait(tgerr.New(420, "FLOOD_WAIT_7"))
	if !ok {
		t.Fatal("flood wait not detected")
	}
	if !strings.Contains(hint, "7s") {
		t.Fatalf("hint = %q, want the wait duration", hint)
	}
}
