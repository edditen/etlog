package handler

import (
	"github.com/EdgarTeng/etlog/common/utils"
	"github.com/EdgarTeng/etlog/config"
	"github.com/EdgarTeng/etlog/core"
	"github.com/pkg/errors"
	"os"
)

type FileHandler struct {
	*BaseHandler
	file *os.File
}

func NewFileHandler(conf *config.HandlerConfig) *FileHandler {
	return &FileHandler{
		BaseHandler: NewBaseHandler(conf),
	}
}

func (fh *FileHandler) Init() (err error) {
	if err = fh.BaseHandler.Init(); err != nil {
		return err
	}

	path := utils.FirstSubstring(fh.BaseHandler.handlerConfig.File, "/")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			return err
		}

	}

	fh.file, err = os.OpenFile(fh.BaseHandler.handlerConfig.File,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "open file error")
	}
	return nil
}

func (fh *FileHandler) Handle(entry *core.LogEntry) error {
	if !fh.BaseHandler.Contains(entry.Level) {
		return nil
	}
	msg := fh.BaseHandler.formatter.Format(entry)
	if _, err := fh.file.WriteString(msg); err != nil {
		return errors.Wrap(err, "write file error")
	}
	return nil
}

func (fh *FileHandler) Shutdown() {
	_ = fh.file.Close()
}
