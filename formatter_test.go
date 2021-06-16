package etlog

import (
	"testing"
	"time"
)

func TestShortFormatter_Format(t *testing.T) {
	t.Run("when short format then simple output", func(t *testing.T) {
		formatter := NewShortFormatter()
		date := time.Date(2021, 6, 15, 12, 20, 45, 152*1e6, time.UTC)
		got := formatter.Format(date, NewLevel("INFO"), "hello world")
		t.Log(got)
		expect := "2021-06-15 12:20:45.152|INFO|hello world\n"
		if got != expect {
			t.Errorf("got: %s, expected: %s", got, expect)
		}

	})
}
