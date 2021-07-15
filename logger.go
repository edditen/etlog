package etlog

import (
	"github.com/EdgarTeng/etlog/common/utils"
	"github.com/EdgarTeng/etlog/config"
	"github.com/EdgarTeng/etlog/core"
	"github.com/EdgarTeng/etlog/handler"
	"github.com/EdgarTeng/etlog/opt"
	"log"
	"os"
	"time"
)

const (
	DefaultConfigPath = "log.yaml"
	BaseSkip          = 5
)

var (
	defaultErr = func() error {
		return nil
	}
	defaultFields = func() core.Fields {
		return core.Fields{}
	}
	defaultMarkers = func() []string {
		return []string{""}
	}

	defaultLogger = func() *EtLogger {
		return &EtLogger{
			configPath: DefaultConfigPath,
			logLevel:   core.DEBUG,
			errLog:     log.New(os.Stderr, "error:", log.LstdFlags),
			infoLog:    log.New(os.Stdout, "", log.LstdFlags),
			sourceSkip: getSourceSkip(0),
			preFns:     make([]opt.LogFunc, 0),
			postFns:    make([]opt.LogFunc, 0),
			conf:       config.DefaultConfig,
		}
	}
)

var Log Logger

func init() {
	SetDefaultLog(newEtLogger())
}

func SetDefaultLog(l Logger) {
	Log = l
}

type OptionFunc func(logger *EtLogger) error

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Data(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
	WithField(field string, v interface{}) Logger
	WithFields(fields core.Fields) Logger
	WithError(err error) Logger
	WithMarkers(markers ...string) Logger
	Enable(level core.Level) bool
}

type Handlers = []handler.Handler

type internalLogger struct {
	err      error
	fields   core.Fields
	markers  []string
	etLogger *EtLogger
}

type EtLogger struct {
	configPath string
	logLevel   core.Level
	errLog     opt.Printfer
	infoLog    opt.Printfer
	sourceSkip int
	preFns     []opt.LogFunc
	postFns    []opt.LogFunc
	conf       *config.Config
	handlers   map[string]*Handlers
	internal   *internalLogger
}

func SetConfigPath(configPath string) OptionFunc {
	return func(logger *EtLogger) error {
		logger.configPath = configPath
		return nil
	}
}

func SetErrorLog(errLog opt.Printfer) OptionFunc {
	return func(logger *EtLogger) error {
		logger.errLog = errLog
		return nil
	}
}

func SetInfoLog(infoLog opt.Printfer) OptionFunc {
	return func(logger *EtLogger) error {
		logger.infoLog = infoLog
		return nil
	}
}

func SetPreLog(preLog ...opt.LogFunc) OptionFunc {
	return func(logger *EtLogger) error {
		logger.preFns = preLog
		return nil
	}
}

func SetSourceSkip(skip uint) OptionFunc {
	return func(logger *EtLogger) error {
		logger.sourceSkip = getSourceSkip(int(skip))
		return nil
	}
}

func getSourceSkip(skip int) int {
	return BaseSkip + skip
}

func SetPostLog(postLog ...opt.LogFunc) OptionFunc {
	return func(logger *EtLogger) error {
		logger.postFns = postLog
		return nil
	}
}

func NewEtLogger(options ...OptionFunc) (*EtLogger, error) {
	logger := defaultLogger()

	for _, option := range options {
		if err := option(logger); err != nil {
			return nil, err
		}
	}

	logger.initInternalLogs()
	logger.conf = config.NewConfig(logger.configPath)
	if err := logger.conf.Init(); err != nil {
		return nil, err
	}

	if err := logger.init(); err != nil {
		return nil, err
	}

	return logger, nil
}

func newEtLogger() *EtLogger {
	logger := defaultLogger()
	logger.initInternalLogs()
	if err := logger.init(); err != nil {
		opt.GetErrLog().Printf("init default etlog err: %+v", err)
		return nil
	}

	return logger
}

func (el *EtLogger) init() (err error) {
	if el.handlers, err = initHandlers(el.conf); err != nil {
		return err
	}
	el.logLevel = core.NewLevel(el.conf.LogConf.Level)
	el.internal = newInternalLogger(el)

	return nil
}

func (el *EtLogger) initInternalLogs() {
	if el.errLog != nil {
		opt.SetErrLog(opt.NewInternalLog(el.errLog, 10, 1000))
	}
	if el.infoLog != nil {
		opt.SetInfoLog(opt.NewInternalLog(el.infoLog, 100, 1000))
	}
}

func initHandlers(conf *config.Config) (map[string]*Handlers, error) {
	handlers := make(map[string]*Handlers, 0)
	for _, handlerConf := range conf.LogConf.Handlers {
		h := handler.HandlerFactory(&handlerConf)
		if err := h.Init(); err != nil {
			return nil, err
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
	return handlers, nil

}

func newInternalLogger(etLogger *EtLogger) *internalLogger {
	return &internalLogger{
		err:      defaultErr(),
		fields:   defaultFields(),
		markers:  defaultMarkers(),
		etLogger: etLogger,
	}

}

func (il *internalLogger) WithField(field string, v interface{}) Logger {
	il.fields[field] = v
	return il
}

func (il *internalLogger) WithFields(fields core.Fields) Logger {
	for k, v := range fields {
		il.fields[k] = v
	}
	return il
}

func (il *internalLogger) WithMarkers(markers ...string) Logger {
	if len(markers) == 0 {
		return il
	}
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
	if fname, line, funcName, ok := utils.ShortSourceLoc(il.etLogger.sourceSkip); ok {
		entry.UseLoc = true
		entry.FileName = fname
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

	for marker, handlers := range il.etLogger.handlers {
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

	il.clean()

	return
}

func (il *internalLogger) clean() {
	il.err = defaultErr()
	il.fields = defaultFields()
	il.markers = defaultMarkers()
}

func (il *internalLogger) preHandle(entry *core.LogEntry) {
	if len(il.etLogger.preFns) == 0 {
		return
	}

	e := newLogE(entry)
	for _, fn := range il.etLogger.preFns {
		handleFunc(fn, e)
	}

}

func (il *internalLogger) postHandle(entry *core.LogEntry) {
	if len(il.etLogger.postFns) == 0 {
		return
	}

	e := newLogE(entry)
	for _, fn := range il.etLogger.postFns {
		handleFunc(fn, e)
	}
}

func newLogE(entry *core.LogEntry) *opt.LogE {
	return &opt.LogE{
		Time:     entry.Time,
		Level:    entry.Level.String(),
		FileName: entry.FileName,
		Line:     entry.Line,
		FuncName: entry.FuncName,
		Marker:   entry.Marker,
		Msg:      entry.Msg,
		Err:      entry.Err,
		Fields:   entry.Fields,
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
	for _, m := range il.markers {
		if m == marker {
			return true
		}
	}
	return false
}

func (il *internalLogger) Enable(level core.Level) bool {
	return il.etLogger.Enable(level)
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
	return newInternalLogger(el).WithError(err)
}

func (el *EtLogger) WithField(field string, v interface{}) Logger {
	return newInternalLogger(el).WithField(field, v)
}

func (el *EtLogger) WithFields(fields core.Fields) Logger {
	return newInternalLogger(el).WithFields(fields)
}

func (el *EtLogger) WithMarkers(markers ...string) Logger {
	return newInternalLogger(el).WithMarkers(markers...)
}

func (el *EtLogger) Enable(level core.Level) bool {
	if level < el.logLevel {
		return false
	}
	return true
}
