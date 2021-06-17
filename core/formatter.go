package core

import (
	"fmt"
	"strings"
)

type Format int

const (
	defaultFormat        = SIMPLE
	SIMPLE        Format = iota
	FULL
)

func NewFormat(format string) Format {
	switch strings.ToUpper(format) {
	case "SIMPLE":
		return SIMPLE
	case "FULL":
		return FULL
	}
	return defaultFormat
}

func (f Format) String() string {
	switch f {
	case SIMPLE:
		return "SIMPLE"
	case FULL:
		return "FULL"
	}
	return ""
}

type Formatter interface {
	Format(entry *LogEntry) string
}

type SimpleFormatter struct {
	format string
}

func FormatterFactory(format Format) Formatter {
	switch format {
	case SIMPLE:
		return NewSimpleFormatter()
	case FULL:
		return NewFullFormatter()
	}
	return NewSimpleFormatter()
}

func NewSimpleFormatter() *SimpleFormatter {
	return &SimpleFormatter{
		format: "%s\t%s\t%s\n",
	}
}

func (s SimpleFormatter) Format(entry *LogEntry) string {
	fmtTime := fmt.Sprintf(entry.Time.Format("2006-01-02 15:04:05.000000"))
	return fmt.Sprintf(s.format, fmtTime, entry.Level, entry.Msg)
}

type FullFormatter struct {
	format string
}

func NewFullFormatter() *FullFormatter {
	return &FullFormatter{
		format: "%s|%s|%s|%s|%s\n",
	}
}

func (s FullFormatter) Format(entry *LogEntry) string {
	fmtTime := fmt.Sprintf(entry.Time.Format("2006-01-02 15:04:05.000000"))
	if entry.SrcValid {
		return fmt.Sprintf(s.format, fmtTime, entry.Level, entry.Line, entry.FuncName, entry.Msg)
	}
	return fmt.Sprintf(s.format, fmtTime, entry.Level, "-", "-", entry.Msg)
}
