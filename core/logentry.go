package core

import (
	"bytes"
	"encoding/json"
	"github.com/EdgarTeng/etlog/common/utils"
	"log"
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
		log.Println(err.Error())
		return []byte{}
	}
	return utils.TrimBytes(builder.Bytes())
}
