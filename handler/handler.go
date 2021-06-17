package handler

import (
	"github.com/EdgarTeng/etlog/config"
	"github.com/EdgarTeng/etlog/core"
	"strings"
)

type HandlerType int

const (
	defaultHandleType             = STD
	STD               HandlerType = iota
	FILE
)

func NewHandlerType(handlerType string) HandlerType {
	switch strings.ToUpper(handlerType) {
	case "STD":
		return STD
	case "FILE":
		return FILE
	}
	return defaultHandleType
}

func (h HandlerType) String() string {
	switch h {
	case STD:
		return "STD"
	case FILE:
		return "FILE"
	}
	return ""
}

type Handler interface {
	Init() error
	Handle(meta *core.LogMeta) error
	Shutdown()
}

type BaseHandler struct {
	handlerConfig *config.HandlerConfig
	formatter     core.Formatter
}

func NewBaseHandler(conf *config.HandlerConfig) *BaseHandler {
	return &BaseHandler{handlerConfig: conf}
}

func (bh *BaseHandler) Init() error {
	bh.DefaultSetting()
	format := core.NewFormat(bh.handlerConfig.Message.Format)
	bh.formatter = core.FormatterFactory(format)
	return nil
}

func (bh *BaseHandler) Handle(meta *core.LogMeta) error {
	return nil
}

func (bh *BaseHandler) Shutdown() {
}

func (bh *BaseHandler) DefaultSetting() {
	if bh.handlerConfig == nil {
		bh.handlerConfig = config.NewHandlerConfig()
	}
	if bh.handlerConfig.File == "" {
		bh.handlerConfig.File = defaultLogFile
	}
	if bh.handlerConfig.Message == nil {
		bh.handlerConfig.Message = config.NewMessageConfig()
	}
	if bh.handlerConfig.Sync == nil {
		bh.handlerConfig.Sync = config.NewSyncConfig()
	}
	if bh.handlerConfig.Rollover == nil {
		bh.handlerConfig.Rollover = config.NewRolloverConfig()
	}
}

func HandlerFactory(conf *config.HandlerConfig) Handler {
	htype := NewHandlerType(conf.Type)
	switch htype {
	case STD:
		return NewStdHandler(conf)
	case FILE:
		return NewFileHandler(conf)
	}
	return NewStdHandler(conf)
}
