package etlog

import "strings"

type Level int

const (
	defaultLevel       = DEBUG
	DEBUG        Level = iota
	INFO
	DATA
	WARN
	ERROR
	FATAL
)

func NewLevel(level string) Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "DATA":
		return DATA
	case "WARN":
		return WARN
	case "ERROR":
		return ERROR
	case "FATAL":
		return FATAL
	}
	return defaultLevel

}

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case DATA:
		return "DATA"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	}
	return ""
}
