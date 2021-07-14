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
	Enable(level core.Level) bool
}

type Handlers = []handler.Handler

type internalLogger struct {
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
	internal   *internalLogger
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

func (el *EtLogger) setLoggers() {
	if el.errLog != nil {
		opt.SetErrLog(opt.NewInternalLog(el.errLog, 10, 1000))
	}
	if el.infoLog != nil {
		opt.SetInfoLog(opt.NewInternalLog(el.infoLog, 100, 1000))
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

func newLoggerInternal(handlers map[string]*Handlers, level core.Level) *internalLogger {
	return &internalLogger{
		logLevel:   level,
		handlers:   handlers,
		fields:     make(map[string]interface{}, 0),
		sourceFlag: DefaultSkip,
		markers:    []string{""},
		preFns:     make([]opt.LogFunc, 0),
		postFns:    make([]opt.LogFunc, 0),
	}

}

func (il *internalLogger) WithField(field string, v interface{}) Logger {
	if il.fields == nil {
		il.fields = make(map[string]interface{}, 0)
	}
	il.fields[field] = v
	return il
}

func (il *internalLogger) WithMarkers(markers ...string) Logger {
	il.markers = markers
	return il
}

func (il *internalLogger) WithError(err error) Logger {
	il.err = err
	return il
}

func (il *internalLogger) Debug(msg string) {
	il.Log(core.DEBUG, msg)
}

func (il *internalLogger) Info(msg string) {
	il.Log(core.INFO, msg)
}

func (il *internalLogger) Data(msg string) {
	il.Log(core.DATA, msg)
}

func (il *internalLogger) Warn(msg string) {
	il.Log(core.WARN, msg)
}

func (il *internalLogger) Error(msg string) {
	il.Log(core.ERROR, msg)
}

func (il *internalLogger) Fatal(msg string) {
	il.Log(core.FATAL, msg)
}

func (il *internalLogger) finalize(level core.Level, msg string) (entry *core.LogEntry) {
	entry = core.NewLogEntry()
	entry.Time = time.Now()
	entry.Level = level
	entry.Msg = msg
	entry.Err = il.err
	entry.Fields = il.fields
	if line, funcName, ok := utils.ShortSourceLoc(il.sourceFlag); ok {
		entry.UseLoc = true
		entry.Line = line
		entry.FuncName = funcName
	}

	return entry
}

func (il *internalLogger) Log(level core.Level, msg string) {
	if !il.Enable(level) {
		return
	}
	entry := il.finalize(level, msg)

	for marker, handlers := range il.handlers {
		if handlers == nil || !il.contains(marker) {
			continue
		}

		e := entry.Copy()
		e.Marker = marker

		il.preHandle(e)

		for _, h := range *handlers {
			if err := h.Handle(e); err != nil {
				opt.GetErrLog().Printf("handle log err: %+v\n", err)
			}
		}

		il.postHandle(e)
	}

	return
}

func (il *internalLogger) preHandle(entry *core.LogEntry) {
	if len(il.preFns) == 0 {
		return
	}

	e := newLogE(entry)
	for _, fn := range il.preFns {
		handleFunc(fn, e)
	}

}

func (il *internalLogger) postHandle(entry *core.LogEntry) {
	if len(il.postFns) == 0 {
		return
	}

	e := newLogE(entry)
	for _, fn := range il.postFns {
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

func (il *internalLogger) contains(marker string) bool {
	if len(il.markers) == 0 {
		return false
	}
	for _, m := range il.markers {
		if m == marker {
			return true
		}
	}
	return false
}

func (il *internalLogger) Enable(level core.Level) bool {
	if level < il.logLevel {
		return false
	}
	return true
}

func (il *internalLogger) setLogFunc(preFns, postFns []opt.LogFunc) {
	il.preFns = preFns
	il.postFns = postFns
}

func (el *EtLogger) newInternal() *internalLogger {
	return &internalLogger{
		logLevel:   el.internal.logLevel,
		handlers:   el.internal.handlers,
		sourceFlag: el.internal.sourceFlag,
		markers:    el.internal.markers,
		fields:     make(map[string]interface{}),
		preFns:     el.preFns,
		postFns:    el.postFns,
	}
}

func (el *EtLogger) Debug(msg string) {
	el.internal.Log(core.DEBUG, msg)
}

func (el *EtLogger) Info(msg string) {
	el.internal.Log(core.INFO, msg)
}

func (el *EtLogger) Data(msg string) {
	el.internal.Log(core.DATA, msg)
}

func (el *EtLogger) Warn(msg string) {
	el.internal.Log(core.WARN, msg)
}

func (el *EtLogger) Error(msg string) {
	el.internal.Log(core.ERROR, msg)
}

func (el *EtLogger) Fatal(msg string) {
	el.internal.Log(core.FATAL, msg)
}

func (el *EtLogger) WithError(err error) Logger {
	return el.newInternal().WithError(err)
}

func (el *EtLogger) WithField(field string, v interface{}) Logger {
	return el.newInternal().WithField(field, v)
}

func (el *EtLogger) WithMarkers(markers ...string) Logger {
	return el.newInternal().WithMarkers(markers...)
}

func (el *EtLogger) Enable(level core.Level) bool {
	return el.internal.Enable(level)
}
