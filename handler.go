package etlog

import (
	"github.com/pkg/errors"
	"os"
)

type Handler interface {
	Handle(msg string) error
}

type FileHandler struct {
	filePath string
	file     *os.File
}

func NewFileHandler(file string) *FileHandler {
	return &FileHandler{filePath: file}
}

func (fh *FileHandler) Init() (err error) {
	fh.file, err = os.OpenFile(fh.filePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "open file error")
	}
	return nil
}

func (fh *FileHandler) Shutdown() {
	_ = fh.file.Close()
}

func (fh *FileHandler) Handle(msg string) error {
	if _, err := fh.file.WriteString(msg); err != nil {
		return errors.Wrap(err, "write file error")
	}
	return nil
}
