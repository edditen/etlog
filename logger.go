package etlog

import (
	"github.com/EdgarTeng/etlog/common/utils"
	"github.com/EdgarTeng/etlog/config"
	"github.com/EdgarTeng/etlog/core"
	"github.com/EdgarTeng/etlog/handler"
	"github.com/EdgarTeng/etlog/opt"
	"time"
)

const (
	DefaultConfigPath = "log.yaml"
	DefaultSkip       = 5
)

var Log Logger

func init() {
	log := newEtLogger()
	Log = log
}

func SetDefaultLog(log Logger) {
	Log = log
}

type LoggerOptionFunc func(logger *EtLogger) error

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

type loggerInternal struct {
	logLevel   core.Level
	handlers   map[string]*Handlers
	sourceFlag int
	err        error
	fields     map[string]interface{}
	markers    []string
	preFns     []opt.LogFunc
	postFns    []opt.LogFunc
}

type EtLogger struct {
	configPath string
	conf       *config.Config
	internal   *loggerInternal
	errLog     opt.Printfer
	infoLog    opt.Printfer
	preFns     []opt.LogFunc
	postFns    []opt.LogFunc
}

func SetConfigPath(configPath string) LoggerOptionFunc {
	return func(logger *EtLogger) error {
		logger.configPath = configPath
		return nil
	}
}

func SetErrorLog(errLog opt.Printfer) LoggerOptionFunc {
	return func(logger *EtLogger) error {
		logger.errLog = errLog
		return nil
	}
}

func SetInfoLog(infoLog opt.Printfer) LoggerOptionFunc {
	return func(logger *EtLogger) error {
		logger.infoLog = infoLog
		return nil
	}
}

func SetPreLog(preLog ...opt.LogFunc) LoggerOptionFunc {
	return func(logger *EtLogger) error {
		logger.preFns = preLog
		return nil
	}
}

func SetPostLog(postLog ...opt.LogFunc) LoggerOptionFunc {
	return func(logger *EtLogger) error {
		logger.postFns = postLog
		return nil
	}
}

func NewEtLogger(options ...LoggerOptionFunc) (*EtLogger, error) {
	logger := &EtLogger{
		configPath: DefaultConfigPath,
	}

	for _, option := range options {
		if err := option(logger); err != nil {
			return nil, err
		}
	}

	logger.setLoggers()

	conf := config.NewConfig(logger.configPath)
	if err := conf.Init(); err != nil {
		return nil, err
	}

	handlers, level, err := getHandlers(conf)
	if err != nil {
		return nil, err
	}

	logger.internal = newLoggerInternal(handlers, level)
	logger.internal.setLogFunc(logger.preFns, logger.postFns)

	return logger, nil
}

func (dl *EtLogger) setLoggers() {
	if dl.errLog != nil {
		opt.SetErrLog(opt.NewInternalLog(dl.errLog, 10, 1000))
	}
	if dl.infoLog != nil {
		opt.SetInfoLog(opt.NewInternalLog(dl.infoLog, 100, 1000))
	}
}

func newEtLogger() *EtLogger {
	logger := &EtLogger{
		configPath: DefaultConfigPath,
	}

	handlers, level, err := getHandlers(config.DefaultConfig)
	if err != nil {
		opt.GetErrLog().Printf("new default logger err:%+v\n", err)
		return nil
	}

	logger.setLoggers()

	logger.internal = newLoggerInternal(handlers, level)

	return logger
}

func getHandlers(conf *config.Config) (map[string]*Handlers, core.Level, error) {
	level := core.NewLevel(conf.LogConf.Level)

	handlers := make(map[string]*Handlers, 0)
	for _, handlerConf := range conf.LogConf.Handlers {
		h := handler.HandlerFactory(&handlerConf)
		if err := h.Init(); err != nil {
			return nil, core.DEBUG, err
		}

		if hs, ok := handlers[handlerConf.Marker]; ok {
			if hs == nil {
				hs = &Handlers{h}
			} else {
				*hs = append(*hs, h)
			}
		} else {
			handlers[handlerConf.Marker] = &Handlers{h}
		}

	}
	return handlers, level, nil

}

func newLoggerInternal(handlers map[string]*Handlers, level core.Level) *loggerInternal {
	return &loggerInternal{
		logLevel:   level,
		handlers:   handlers,
		fields:     make(map[string]interface{}, 0),
		sourceFlag: DefaultSkip,
		markers:    []string{""},
		preFns:     make([]opt.LogFunc, 0),
		postFns:    make([]opt.LogFunc, 0),
	}

}

func (li *loggerInternal) WithField(field string, v interface{}) Logger {
	if li.fields == nil {
		li.fields = make(map[string]interface{}, 0)
	}
	li.fields[field] = v
	return li
}

func (li *loggerInternal) WithMarkers(markers ...string) Logger {
	li.markers = markers
	return li
}

func (li *loggerInternal) WithError(err error) Logger {
	li.err = err
	return li
}

func (li *loggerInternal) Debug(msg string) {
	li.Log(core.DEBUG, msg)
}

func (li *loggerInternal) Info(msg string) {
	li.Log(core.INFO, msg)
}

func (li *loggerInternal) Data(msg string) {
	li.Log(core.DATA, msg)
}

func (li *loggerInternal) Warn(msg string) {
	li.Log(core.WARN, msg)
}

func (li *loggerInternal) Error(msg string) {
	li.Log(core.ERROR, msg)
}

func (li *loggerInternal) Fatal(msg string) {
	li.Log(core.FATAL, msg)
}

func (li *loggerInternal) finalize(level core.Level, msg string) (entry *core.LogEntry) {
	entry = core.NewLogEntry()
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

func (li *loggerInternal) Log(level core.Level, msg string) {
	if !li.Enable(level) {
		return
	}
	entry := li.finalize(level, msg)

	for marker, handlers := range li.handlers {
		if handlers == nil || !li.contains(marker) {
			continue
		}

		e := entry.Copy()
		e.Marker = marker

		li.preHandle(e)

		for _, h := range *handlers {
			if err := h.Handle(e); err != nil {
				opt.GetErrLog().Printf("handle log err: %+v\n", err)
			}
		}

		li.postHandle(e)
	}

	return
}

func (li *loggerInternal) preHandle(entry *core.LogEntry) {
	if len(li.preFns) == 0 {
		return
	}

	e := newLogE(entry)
	for _, fn := range li.preFns {
		handleFunc(fn, e)
	}

}

func (li *loggerInternal) postHandle(entry *core.LogEntry) {
	if len(li.postFns) == 0 {
		return
	}

	e := newLogE(entry)
	for _, fn := range li.postFns {
		handleFunc(fn, e)
	}
}

func newLogE(entry *core.LogEntry) *opt.LogE {
	return &opt.LogE{
		Time:   entry.Time,
		Level:  entry.Level.String(),
		Marker: entry.Marker,
		Msg:    entry.Msg,
		Err:    entry.Err,
		Fields: entry.Fields,
	}
}

func handleFunc(fn opt.LogFunc, e *opt.LogE) {
	defer func() {
		if r := recover(); r != nil {
			opt.GetErrLog().Printf("handle LogFunc panic recover, %v\n", r)
		}
	}()
	fn(e)
}

func (li *loggerInternal) contains(marker string) bool {
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

func (li *loggerInternal) Enable(level core.Level) bool {
	if level < li.logLevel {
		return false
	}
	return true
}

func (li *loggerInternal) setLogFunc(preFns, postFns []opt.LogFunc) {
	li.preFns = preFns
	li.postFns = postFns
}

func (dl *EtLogger) newInternal() *loggerInternal {
	return &loggerInternal{
		logLevel:   dl.internal.logLevel,
		handlers:   dl.internal.handlers,
		sourceFlag: dl.internal.sourceFlag,
		markers:    dl.internal.markers,
		fields:     make(map[string]interface{}),
		preFns:     dl.preFns,
		postFns:    dl.postFns,
	}
}

func (dl *EtLogger) Debug(msg string) {
	dl.internal.Log(core.DEBUG, msg)
}

func (dl *EtLogger) Info(msg string) {
	dl.internal.Log(core.INFO, msg)
}

func (dl *EtLogger) Data(msg string) {
	dl.internal.Log(core.DATA, msg)
}

func (dl *EtLogger) Warn(msg string) {
	dl.internal.Log(core.WARN, msg)
}

func (dl *EtLogger) Error(msg string) {
	dl.internal.Log(core.ERROR, msg)
}

func (dl *EtLogger) Fatal(msg string) {
	dl.internal.Log(core.FATAL, msg)
}

func (dl *EtLogger) WithError(err error) Logger {
	return dl.newInternal().WithError(err)
}

func (dl *EtLogger) WithField(field string, v interface{}) Logger {
	return dl.newInternal().WithField(field, v)
}

func (dl *EtLogger) WithMarkers(markers ...string) Logger {
	return dl.newInternal().WithMarkers(markers...)
}
