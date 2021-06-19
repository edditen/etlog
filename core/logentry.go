package core

import (
	"encoding/json"
	"fmt"
	"strings"
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
	if f == nil || len(f) == 0 {
		return ""
	}
	builder := &strings.Builder{}
	if err := json.NewEncoder(builder).Encode(f); err != nil {
		fmt.Println(err.Error())
		return ""
	}
	return builder.String()
}
