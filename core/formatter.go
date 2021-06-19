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

func (sf SimpleFormatter) Format(entry *LogEntry) string {
	fmtTime := fmt.Sprintf(entry.Time.Format("2006-01-02 15:04:05.000000"))
	return fmt.Sprintf(sf.format, fmtTime, entry.Level, entry.Msg)
}

type FullFormatter struct {
	format string
}

func NewFullFormatter() *FullFormatter {
	return &FullFormatter{
		// format: "time|level|line|func|message|error|fields"
		format: "%s|%s|%s|%s|%s|%s|%s\n",
	}
}

func (ff FullFormatter) Format(entry *LogEntry) string {
	fmtTime := fmt.Sprintf(entry.Time.Format("2006-01-02 15:04:05.000000"))
	builder := &strings.Builder{}
	// timestamp
	builder.WriteString(fmtTime)
	builder.WriteString("|")
	// level
	builder.WriteString(fmt.Sprintf("%s", entry.Level))
	builder.WriteString("|")

	// line & func
	if entry.SrcValid {
		builder.WriteString(entry.Line)
		builder.WriteString("|")
		builder.WriteString(entry.FuncName)
	} else {
		builder.WriteString("-|-")
	}
	builder.WriteString("|")

	// msg
	builder.WriteString(entry.Msg)
	builder.WriteString("|")

	// error
	if entry.Err != nil {
		builder.WriteString(fmt.Sprintf("%s", entry.Err))
	} else {
		builder.WriteString("-")
	}
	builder.WriteString("|")

	// fields
	if entry.Fields != nil && len(entry.Fields) > 0 {
		builder.WriteString(fmt.Sprintf("%s", entry.Fields))

	}
	builder.WriteString("\n")

	return builder.String()
}
