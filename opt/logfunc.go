package opt

import "time"

type LogFunc func(e *LogE)

type LogE struct {
	Time   time.Time              `json:"time"`
	Level  string                 `json:"level"`
	Msg    string                 `json:"msg"`
	Marker string                 `json:"marker"`
	Err    error                  `json:"error"`
	Fields map[string]interface{} `json:"fields"`
}
