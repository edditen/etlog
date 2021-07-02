package core

import (
	"fmt"
	"github.com/EdgarTeng/etlog/common/bufferpool"
	"strings"
)

type Format int

const (
	defaultTimeFormat = "2006-01-02 15:04:05.000000"
)

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
	Format(entry *LogEntry) *bufferpool.Buffer
}

type SimpleFormatter struct {
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
	//format: "time  level  msg"
	return &SimpleFormatter{}
}

func (sf *SimpleFormatter) Format(entry *LogEntry) *bufferpool.Buffer {
	buf := bufferpool.Borrow()
	// timestamp
	buf.AppendString(entry.Time.Format(defaultTimeFormat))
	buf.AppendByte(' ')

	// level
	buf.AppendByte('[')
	buf.AppendString(fmt.Sprint(entry.Level))
	buf.AppendByte(']')
	buf.AppendByte('\t')

	// msg
	buf.AppendString("|msg:=")
	buf.AppendValue(entry.Msg)

	if entry.Err != nil {
		buf.AppendString("|err:=")
		buf.AppendString(fmt.Sprint(entry.Err))
	}

	if entry.Fields != nil && len(entry.Fields) > 0 {
		buf.AppendString("|fields:=")
		buf.AppendBytes(entry.Fields.Bytes())
	}

	buf.AppendNewLine()

	return buf
}

type FullFormatter struct {
	format string
}

func NewFullFormatter() *FullFormatter {
	// format: "time|level|line|func|message|error|fields"
	return &FullFormatter{}
}

func (ff *FullFormatter) Format(entry *LogEntry) *bufferpool.Buffer {
	buf := bufferpool.Borrow()
	// timestamp
	buf.AppendString(entry.Time.Format(defaultTimeFormat))
	buf.AppendByte('|')

	// level
	buf.AppendString(fmt.Sprint(entry.Level))
	buf.AppendByte('|')

	// line & func
	if entry.UseLoc {
		buf.AppendString(entry.Line)
		buf.AppendByte('|')
		buf.AppendString(entry.FuncName)
	} else {
		buf.AppendString("-|-")
	}
	buf.AppendByte('|')

	// msg
	buf.AppendString(entry.Msg)
	buf.AppendByte('|')

	// error
	if entry.Err != nil {
		buf.AppendString(fmt.Sprint(entry.Err))
	} else {
		buf.AppendByte('-')
	}
	buf.AppendByte('|')

	// fields
	if entry.Fields != nil && len(entry.Fields) > 0 {
		buf.AppendBytes(entry.Fields.Bytes())
	}
	buf.AppendNewLine()

	return buf
}
