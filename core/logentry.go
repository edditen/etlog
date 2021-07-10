package core

import (
	"bytes"
	"encoding/json"
	"github.com/EdgarTeng/etlog/common/utils"
	"github.com/EdgarTeng/etlog/opt"
	"time"
)

type Fields map[string]interface{}

type LogEntry struct {
	Time     time.Time `json:"time"`
	Level    Level     `json:"level"`
	Line     string    `json:"line"`
	FuncName string    `json:"func"`
	Msg      string    `json:"msg"`
	Marker   string    `json:"marker"`
	Err      error     `json:"error"`
	Fields   Fields    `json:"fields"`
	UseLoc   bool      `json:"-"`
}

func NewLogMeta() *LogEntry {
	return &LogEntry{}
}

func (le *LogEntry) Copy() *LogEntry {
	return &LogEntry{
		Time:     le.Time,
		Level:    le.Level,
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
