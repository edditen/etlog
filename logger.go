package etlog

import (
	"github.com/EdgarTeng/etlog/common/utils"
	"github.com/EdgarTeng/etlog/config"
	"github.com/EdgarTeng/etlog/core"
	"github.com/EdgarTeng/etlog/handler"
	"time"
)

var Log Logger

func SetDefaultLog(log Logger) {
	Log = log
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Data(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
}

type LoggerInternal struct {
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
		handlers:   make([]handler.Handler, 0),
		sourceFlag: 5,
	}
}

func (li *LoggerInternal) Init() error {

	for _, handlerConf := range li.conf.LogConf.Handlers {
		handler := handler.HandlerFactory(&handlerConf)
		if err := handler.Init(); err != nil {
			return err
		}
		li.handlers = append(li.handlers, handler)
	}
	return nil
}

func (li *LoggerInternal) Finalize(level core.Level, msg string) (meta *core.LogMeta) {
	meta = core.NewLogMeta()
	meta.Time = time.Now()
	meta.Level = level
	meta.Msg = msg
	if line, funcName, ok := utils.ShortSourceLoc(li.sourceFlag); ok {
		meta.SrcValid = true
		meta.Line = line
		meta.FuncName = funcName
	}

	return meta
}

func (li *LoggerInternal) Log(level core.Level, msg string) *LoggerInternal {
	meta := li.Finalize(level, msg)
	for _, handler := range li.handlers {
		handler.Handle(meta)
	}
	return li
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
