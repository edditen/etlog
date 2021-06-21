package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type Fields map[string]interface{}

type LogEntry struct {
	Time     time.Time `json:"time"`
	Level    Level     `json:"level"`
	Line     string    `json:"line"`
	FuncName string    `json:"func"`
	Msg      string    `json:"msg"`
	Err      error     `json:"error"`
	Fields   Fields    `json:"fields"`
	SrcValid bool      `json:"-"`
}

func NewLogMeta() *LogEntry {
	return &LogEntry{}
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
		fmt.Println(err.Error())
		return []byte{}
	}
	return trimNewline(builder.Bytes())
}

func trimNewline(bs []byte) []byte {
	if i := len(bs) - 1; i >= 0 {
		if bs[i] == '\n' {
			bs = bs[:i]
		}
	}
	return bs
}
