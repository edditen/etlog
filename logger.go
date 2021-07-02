package etlog

import (
	"github.com/EdgarTeng/etlog/common/utils"
	"github.com/EdgarTeng/etlog/config"
	"github.com/EdgarTeng/etlog/core"
	"github.com/EdgarTeng/etlog/handler"
	"log"
	"time"
)

const (
	DefaultConfigPath = "log.yaml"
	DefaultSkip       = 5
)

var Log Logger

func init() {
	log := newDefaultLogger()
	Log = log
}

func SetDefaultLog(log Logger) {
	Log = log
}

type LoggerOptionFunc func(logger *DefaultLogger) error

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Data(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
	WithField(field string, v interface{}) Logger
	WithError(err error) Logger
	WithMarkers(markers ...string) Logger
}

type Handlers = []handler.Handler

type LoggerInternal struct {
	logLevel   core.Level
	handlers   map[string]*Handlers
	sourceFlag int
	err        error
	fields     map[string]interface{}
	markers    []string
}

type DefaultLogger struct {
	configPath string
	conf       *config.Config
	internal   *LoggerInternal
}

func SetConfigPath(configPath string) LoggerOptionFunc {
	return func(logger *DefaultLogger) error {
		logger.configPath = configPath
		return nil
	}
}

func NewDefaultLogger(options ...LoggerOptionFunc) (*DefaultLogger, error) {
	logger := &DefaultLogger{
		configPath: DefaultConfigPath,
	}

	for _, option := range options {
		if err := option(logger); err != nil {
			return nil, err
		}
	}

	conf := config.NewConfig(logger.configPath)
	if err := conf.Init(); err != nil {
		return nil, err
	}

	handlers, level, err := getHandlers(conf)
	if err != nil {
		return nil, err
	}

	logger.internal = NewLoggerInternal(handlers, level)

	return logger, nil
}

func newDefaultLogger() *DefaultLogger {
	logger := &DefaultLogger{
		configPath: DefaultConfigPath,
	}

	handlers, level, err := getHandlers(config.DefaultConfig)
	if err != nil {
		log.Println("new default logger error", err)
		return nil
	}

	logger.internal = NewLoggerInternal(handlers, level)

	return logger
}

func getHandlers(conf *config.Config) (map[string]*Handlers, core.Level, error) {
	level := core.NewLevel(conf.LogConf.Level)

	handlers := make(map[string]*Handlers, 0)
	for _, handlerConf := range conf.LogConf.Handlers {
		handler := handler.HandlerFactory(&handlerConf)
		if err := handler.Init(); err != nil {
			return nil, core.DEBUG, err
		}

		if hs, ok := handlers[handlerConf.Marker]; ok {
			if hs == nil {
				hs = &Handlers{handler}
			} else {
				*hs = append(*hs, handler)
			}
		} else {
			handlers[handlerConf.Marker] = &Handlers{handler}
		}

	}
	return handlers, level, nil

}

func NewLoggerInternal(handlers map[string]*Handlers, level core.Level) *LoggerInternal {
	return &LoggerInternal{
		logLevel:   level,
		handlers:   handlers,
		fields:     make(map[string]interface{}, 0),
		sourceFlag: DefaultSkip,
		markers:    []string{""},
	}

}

func (li *LoggerInternal) WithField(field string, v interface{}) Logger {
	if li.fields == nil {
		li.fields = make(map[string]interface{}, 0)
	}
	li.fields[field] = v
	return li
}

func (li *LoggerInternal) WithMarkers(markers ...string) Logger {
	li.markers = markers
	return li
}

func (li *LoggerInternal) WithError(err error) Logger {
	li.err = err
	return li
}

func (li *LoggerInternal) Debug(msg string) {
	li.Log(core.DEBUG, msg)
}

func (li *LoggerInternal) Info(msg string) {
	li.Log(core.INFO, msg)
}

func (li *LoggerInternal) Data(msg string) {
	li.Log(core.DATA, msg)
}

func (li *LoggerInternal) Warn(msg string) {
	li.Log(core.WARN, msg)
}

func (li *LoggerInternal) Error(msg string) {
	li.Log(core.ERROR, msg)
}

func (li *LoggerInternal) Fatal(msg string) {
	li.Log(core.FATAL, msg)
}

func (li *LoggerInternal) finalize(level core.Level, msg string) (entry *core.LogEntry) {
	entry = core.NewLogMeta()
	entry.Time = time.Now()
	entry.Level = level
	entry.Msg = msg
	entry.Err = li.err
	entry.Fields = li.fields
	if line, funcName, ok := utils.ShortSourceLoc(li.sourceFlag); ok {
		entry.UseLoc = true
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

	for marker, handlers := range li.handlers {
		if handlers == nil {
			continue
		}
		if li.contains(marker) {
			e := entry.Copy()
			e.Marker = marker

			for _, handler := range *handlers {
				handler.Handle(e)
			}
		}
	}

	return
}

func (li *LoggerInternal) contains(marker string) bool {
	if len(li.markers) == 0 {
		return false
	}
	for _, m := range li.markers {
		if m == marker {
			return true
		}
	}
	return false
}

func (li *LoggerInternal) Enable(level core.Level) bool {
	if level < li.logLevel {
		return false
	}
	return true
}

func (dl *DefaultLogger) newInternal() *LoggerInternal {
	return &LoggerInternal{
		logLevel:   dl.internal.logLevel,
		handlers:   dl.internal.handlers,
		sourceFlag: dl.internal.sourceFlag,
		markers:    dl.internal.markers,
		fields:     make(map[string]interface{}),
	}
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

func (dl *DefaultLogger) WithError(err error) Logger {
	return dl.newInternal().WithError(err)
}

func (dl *DefaultLogger) WithField(field string, v interface{}) Logger {
	return dl.newInternal().WithField(field, v)
}

func (dl *DefaultLogger) WithMarkers(markers ...string) Logger {
	return dl.newInternal().WithMarkers(markers...)
}
