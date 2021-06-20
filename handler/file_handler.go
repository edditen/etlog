package handler

import (
	"github.com/EdgarTeng/etlog/config"
	"github.com/EdgarTeng/etlog/core"
	"github.com/pkg/errors"
	"io/fs"
	"os"
	"path"
)

const (
	fileFlag                     = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	fileMode         fs.FileMode = 0644
	backupTimeFormat             = "2006-01-02.150405"
)

type FileHandler struct {
	*BaseHandler
	fileWriter *os.File
	filePath   string
	fileDir    string
	fileExt    string
	fileName   string
	rotateSize int64
	size       int64
}

func NewFileHandler(conf *config.HandlerConfig) *FileHandler {
	return &FileHandler{
		BaseHandler: NewBaseHandler(conf),
		filePath:    conf.File,
		fileDir:     path.Dir(conf.File),
		fileName:    path.Base(conf.File),
		fileExt:     path.Ext(conf.File),
	}
}

func (fh *FileHandler) Init() (err error) {
	if err = fh.BaseHandler.Init(); err != nil {
		return err
	}

	if err = os.MkdirAll(fh.fileDir, os.ModePerm); err != nil {
		return errors.Wrap(err, "create dir error")
	}

	fh.fileWriter, err = os.OpenFile(fh.filePath, fileFlag, fileMode)
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
	if _, err := fh.fileWriter.WriteString(msg); err != nil {
		return errors.Wrap(err, "write file error")
	}
	return nil
}

func (fh *FileHandler) Shutdown() {
	_ = fh.fileWriter.Close()
}
