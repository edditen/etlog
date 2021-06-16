package etlog

import (
	"fmt"
	"testing"
	"time"
)

func TestShortFormatter_Format(t *testing.T) {
	t.Run("when short format then simple output", func(t *testing.T) {
		formatter := NewShortFormatter()
		s := formatter.Format(time.Now(), NewLevel("INFO"), "hello world")
		fmt.Print(s)
	})
}
