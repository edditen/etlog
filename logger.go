package etlog

import stdlog "log"

var Log Logger

func init() {
	Log = NewStdLogger()
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
}

type StdLogger struct {
}

func NewStdLogger() *StdLogger {
	return &StdLogger{}
}

func (sl *StdLogger) Debug(msg string) {
	stdlog.Println(msg)
}

func (sl *StdLogger) Info(msg string) {
	stdlog.Println(msg)
}

func (sl *StdLogger) Warn(msg string) {
	stdlog.Println(msg)
}

func (sl *StdLogger) Error(msg string) {
	stdlog.Println(msg)
}

func (sl *StdLogger) Fatal(msg string) {
	stdlog.Println(msg)
}
