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
	Format(meta *LogMeta) string
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

func (s SimpleFormatter) Format(meta *LogMeta) string {
	fmtTime := fmt.Sprintf(meta.Time.Format("2006-01-02 15:04:05.000000"))
	return fmt.Sprintf(s.format, fmtTime, meta.Level, meta.Msg)
}

type FullFormatter struct {
	format string
}

func NewFullFormatter() *FullFormatter {
	return &FullFormatter{
		format: "%s|%s|%s|%s|%s\n",
	}
}

func (s FullFormatter) Format(meta *LogMeta) string {
	fmtTime := fmt.Sprintf(meta.Time.Format("2006-01-02 15:04:05.000000"))
	if meta.SrcValid {
		return fmt.Sprintf(s.format, fmtTime, meta.Level, meta.Line, meta.FuncName, meta.Msg)
	}
	return fmt.Sprintf(s.format, fmtTime, meta.Level, "-", "-", meta.Msg)
}
