package handler

import (
	"fmt"
	"github.com/EdgarTeng/etlog/config"
	"github.com/EdgarTeng/etlog/core"
	"github.com/pkg/errors"
)

type StdHandler struct {
	*BaseHandler
}

func NewStdHandler(handlerConf *config.HandlerConfig) *StdHandler {
	return &StdHandler{
		BaseHandler: NewBaseHandler(handlerConf),
	}
}

func (sh *StdHandler) Init() error {
	return sh.BaseHandler.Init()
}

func (sh *StdHandler) Handle(entry *core.LogEntry) error {
	if !sh.BaseHandler.MarkerMatched(entry.Marker) {
		return nil
	}
	if !sh.BaseHandler.Contains(entry.Level) {
		return nil
	}
	msg := sh.BaseHandler.formatter.Format(entry)
	if _, err := fmt.Print(msg); err != nil {
		return errors.Wrap(err, "std print error")
	}
	return nil
}
