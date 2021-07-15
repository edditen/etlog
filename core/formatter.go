package core

import (
	"encoding/json"
	"fmt"
	"github.com/EdgarTeng/etlog/common/bufferpool"
	"github.com/EdgarTeng/etlog/opt"
	"strings"
	"time"
)

type Format int

const (
	defaultTimeFormat = "2006-01-02 15:04:05.000000"
)

const (
	defaultFormat        = SIMPLE
	SIMPLE        Format = iota
	FULL
	JSON
)

func NewFormat(format string) Format {
	switch strings.ToUpper(format) {
	case "SIMPLE":
		return SIMPLE
	case "FULL":
		return FULL
	case "JSON":
		return JSON
	}
	return defaultFormat
}

func (f Format) String() string {
	switch f {
	case SIMPLE:
		return "SIMPLE"
	case FULL:
		return "FULL"
	case JSON:
		return "JSON"
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
	case JSON:
		return NewJSONFormatter()
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
	buf.AppendString(time.Unix(0, entry.Time).Format(defaultTimeFormat))
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
}

func NewFullFormatter() *FullFormatter {
	// format: "time|level|src:line|func|message|error|fields"
	return &FullFormatter{}
}

func (ff *FullFormatter) Format(entry *LogEntry) *bufferpool.Buffer {
	buf := bufferpool.Borrow()
	// timestamp
	buf.AppendString(time.Unix(0, entry.Time).Format(defaultTimeFormat))
	buf.AppendByte('|')

	// level
	buf.AppendString(fmt.Sprint(entry.Level))
	buf.AppendByte('|')

	// line & func
	if entry.UseLoc {
		buf.AppendString(entry.SrcFile)
		buf.AppendByte(':')
		buf.AppendInt(int64(entry.Line))
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

type JSONFormatter struct {
}

func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

func (jf *JSONFormatter) Format(entry *LogEntry) *bufferpool.Buffer {
	buf := bufferpool.Borrow()
	b, err := json.Marshal(entry)
	if err != nil {
		opt.GetErrLog().Printf("JSONFormatter marshal entry error: %+v\n", err)
		return buf
	}
	buf.AppendBytes(b)
	buf.AppendNewLine()
	return buf
}
