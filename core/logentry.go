package core

import (
	"bytes"
	"encoding/json"
	"github.com/edditen/etlog/common/utils"
	"github.com/edditen/etlog/opt"
	"time"
)

type Fields map[string]interface{}

type LogEntry struct {
	Time     time.Time `json:"time,omitempty"`
	Level    Level     `json:"level,omitempty"`
	SrcFile  string    `json:"srcf,omitempty"`
	Line     int       `json:"line,omitempty"`
	FuncName string    `json:"func,omitempty"`
	Msg      string    `json:"msg,omitempty"`
	Marker   string    `json:"marker,omitempty"`
	Err      error     `json:"error,omitempty"`
	Fields   Fields    `json:"fields,omitempty"`
	UseLoc   bool      `json:"-"`
}

func NewLogEntry() *LogEntry {
	return &LogEntry{}
}

func (le *LogEntry) Copy() *LogEntry {
	return &LogEntry{
		Time:     le.Time,
		Level:    le.Level,
		SrcFile:  le.SrcFile,
		Line:     le.Line,
		FuncName: le.FuncName,
		Msg:      le.Msg,
		Marker:   le.Marker,
		Err:      le.Err,
		Fields:   le.Fields,
		UseLoc:   le.UseLoc,
	}
}

func (f Fields) String() string {
	return string(f.Bytes())
}

func (f Fields) Bytes() []byte {
	if f == nil || len(f) == 0 {
		return []byte{}
	}
	builder := &bytes.Buffer{}
	if err := json.NewEncoder(builder).Encode(f); err != nil {
		opt.GetErrLog().Printf("json encode builder err: %+v\n", err)
		return []byte{}
	}
	return utils.TrimBytes(builder.Bytes())
}
