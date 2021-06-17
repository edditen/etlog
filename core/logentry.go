package core

import (
	"time"
)

type LogMeta struct {
	Time     time.Time `json:"time"`
	Level    Level     `json:"level"`
	Line     string    `json:"line"`
	FuncName string    `json:"func"`
	Msg      string    `json:"msg"`
	SrcValid bool      `json:"-"`
}

func NewLogMeta() *LogMeta {
	return &LogMeta{
		SrcValid: false,
	}
}
