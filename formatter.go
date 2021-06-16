package etlog

import (
	"fmt"
	"time"
)

type Formatter interface {
	Format(t time.Time, level Level, msg string) string
}

type ShortFormatter struct {
	format string
}

func NewShortFormatter() *ShortFormatter {
	return &ShortFormatter{
		format: "%s|%s|%s\n",
	}
}

func (s ShortFormatter) Format(t time.Time, level Level, msg string) string {
	st := fmt.Sprintf(t.Format("2006-01-02 15:04:05.000"))
	return fmt.Sprintf(s.format, st, level, msg)
}
