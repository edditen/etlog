package etlog

import (
	"github.com/EdgarTeng/etlog/common/utils"
	"github.com/EdgarTeng/etlog/config"
	"github.com/EdgarTeng/etlog/core"
	"github.com/EdgarTeng/etlog/handler"
	"log"
	"time"
)

var Log Logger

func SetDefaultLog(log Logger) {
	Log = log
}

type Logger interface {
	WithField(field string, v interface{}) Logger
	Debug(msg string)
	Info(msg string)
	Data(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
}

type LoggerInternal struct {
	logLevel   core.Level
	conf       *config.Config
	handlers   []handler.Handler
	sourceFlag int
}

type DefaultLogger struct {
	conf     *config.Config
	internal *LoggerInternal
}

func NewDefaultLogger(configPath string) (*DefaultLogger, error) {
	conf := config.NewConfig(configPath)
	if err := conf.Init(); err != nil {
		return nil, err
	}
	internal := NewLoggerInternal(conf)
	if err := internal.Init(); err != nil {
		return nil, err
	}
	return &DefaultLogger{
		conf:     conf,
		internal: internal,
	}, nil
}

func NewLoggerInternal(conf *config.Config) *LoggerInternal {
	return &LoggerInternal{
		conf:       conf,
		logLevel:   core.DEBUG,
		handlers:   make([]handler.Handler, 0),
		sourceFlag: 5,
	}
}

func (li *LoggerInternal) Init() error {
	li.logLevel = core.NewLevel(li.conf.LogConf.Level)
	log.Println("[Init] log level:", li.logLevel)
	for _, handlerConf := range li.conf.LogConf.Handlers {
		handler := handler.HandlerFactory(&handlerConf)
		if err := handler.Init(); err != nil {
			return err
		}
		li.handlers = append(li.handlers, handler)
	}
	return nil
}

func (li *LoggerInternal) finalize(level core.Level, msg string) (entry *core.LogEntry) {
	entry = core.NewLogMeta()
	entry.Time = time.Now()
	entry.Level = level
	entry.Msg = msg
	if line, funcName, ok := utils.ShortSourceLoc(li.sourceFlag); ok {
		entry.SrcValid = true
		entry.Line = line
		entry.FuncName = funcName
	}

	return entry
}

func (li *LoggerInternal) Log(level core.Level, msg string) {
	if !li.Enable(level) {
		return
	}
	entry := li.finalize(level, msg)
	for _, handler := range li.handlers {
		handler.Handle(entry)
	}
	return
}

func (li *LoggerInternal) Enable(level core.Level) bool {
	if level < li.logLevel {
		return false
	}
	return true
}

func (dl *DefaultLogger) Debug(msg string) {
	dl.internal.Log(core.DEBUG, msg)
}

func (dl *DefaultLogger) Info(msg string) {
	dl.internal.Log(core.INFO, msg)
}

func (dl *DefaultLogger) Data(msg string) {
	dl.internal.Log(core.DATA, msg)
}

func (dl *DefaultLogger) Warn(msg string) {
	dl.internal.Log(core.WARN, msg)
}

func (dl *DefaultLogger) Error(msg string) {
	dl.internal.Log(core.ERROR, msg)
}

func (dl *DefaultLogger) Fatal(msg string) {
	dl.internal.Log(core.FATAL, msg)
}

func (dl *DefaultLogger) WithField(field string, v interface{}) Logger {
	return dl
}
