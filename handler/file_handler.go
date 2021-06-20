package handler

import (
	"fmt"
	"github.com/EdgarTeng/etlog/common/utils"
	"github.com/EdgarTeng/etlog/config"
	"github.com/EdgarTeng/etlog/core"
	"github.com/pkg/errors"
	"io/fs"
	"math"
	"os"
	"path"
	"sync"
	"time"
)

const (
	fileFlag                        = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	fileMode            fs.FileMode = 0644
	backupTimeFormat                = "2006-01-02.150405"
	defaultLogSize                  = "10G"
	defaultRolloverTime             = "1d"
	defaultBackupTime               = "365d"
	defaultBackupCount              = math.MaxInt32
)

type FileHandler struct {
	*BaseHandler
	fileWriter     *os.File
	filePath       string
	fileDir        string
	fileExt        string
	fileName       string
	rotateSize     int
	rotateInterval int
	backupTime     int
	backupCount    int
	writenSize     int64
	rotateAt       time.Time
	rotateLock     *sync.RWMutex
}

func NewFileHandler(conf *config.HandlerConfig) *FileHandler {
	return &FileHandler{
		BaseHandler: NewBaseHandler(conf),
		rotateLock:  new(sync.RWMutex),
	}
}

func (fh *FileHandler) Init() error {
	if err := fh.BaseHandler.Init(); err != nil {
		return err
	}

	if err := fh.settingFileInfo(); err != nil {
		return err
	}

	if err := fh.settingRolloverSize(); err != nil {
		return err
	}

	if err := fh.settingRolloverInterval(); err != nil {
		return err
	}

	if err := fh.settingBackupTime(); err != nil {
		return err
	}

	if err := fh.settingBackupCount(); err != nil {
		return err
	}

	if err := fh.settingFileWriter(); err != nil {
		return err
	}

	if err := fh.settingWritenSize(); err != nil {
		return err
	}

	fh.rotateAt = fh.nextTimeRotate(fh.rotateInterval)

	return nil
}

func (fh *FileHandler) Handle(entry *core.LogEntry) error {
	if !fh.BaseHandler.Contains(entry.Level) {
		return nil
	}

	if fh.shouldRotate() {
		fh.Rotate()
	}

	msg := fh.BaseHandler.formatter.Format(entry)
	if err := fh.Flush(msg); err != nil {
		return err
	}
	return nil
}

func (fh *FileHandler) Flush(msg string) error {
	fh.rotateLock.RLock()
	defer fh.rotateLock.RUnlock()

	if _, err := fh.fileWriter.WriteString(msg); err != nil {
		return errors.Wrap(err, "write file error")
	}
	fh.writenSize += int64(len([]byte(msg)))
	return nil
}

func (fh *FileHandler) Rotate() error {
	fh.rotateLock.Lock()
	defer fh.rotateLock.Unlock()

	backupName := fh.backupFileName()
	_ = fh.fileWriter.Close()

	if err := os.Rename(fh.filePath, backupName); err != nil {
		return errors.Wrap(err, "rotate file error")
	}
	if err := fh.settingFileWriter(); err != nil {
		return err
	}
	if err := fh.settingWritenSize(); err != nil {
		return err
	}

	fh.rotateAt = fh.nextTimeRotate(fh.rotateInterval)

	return nil
}

func (fh *FileHandler) shouldRotate() bool {
	if time.Now().After(fh.rotateAt) {
		return true
	}
	if fh.writenSize >= int64(fh.rotateSize) {
		return true
	}
	return false
}

func (fh *FileHandler) backupFileName() string {
	filename := fh.fileName[:len(fh.fileName)-len(fh.fileExt)]
	t := time.Now().Format(backupTimeFormat)
	filename = fmt.Sprintf("%s.%s%s", filename, t, fh.fileExt)
	return path.Join(fh.fileDir, filename)
}

func (fh *FileHandler) Shutdown() {
	_ = fh.fileWriter.Close()
}

func (fh *FileHandler) settingFileInfo() (err error) {
	if fh.BaseHandler.handlerConfig.File == "" {
		return errors.New("file is empty")
	}

	fh.filePath = fh.BaseHandler.handlerConfig.File
	fh.fileDir = path.Dir(fh.filePath)
	fh.fileName = path.Base(fh.filePath)
	fh.fileExt = path.Ext(fh.filePath)

	return nil
}

func (fh *FileHandler) settingRolloverSize() (err error) {
	if fh.BaseHandler.handlerConfig.Rollover.RolloverSize == "" {
		fh.BaseHandler.handlerConfig.Rollover.RolloverSize = defaultLogSize
	}

	fh.rotateSize, err = utils.ParseSize(fh.BaseHandler.handlerConfig.Rollover.RolloverSize)
	if err != nil {
		return errors.Wrap(err, "parse rotate size error")
	}
	return nil
}

func (fh *FileHandler) settingRolloverInterval() (err error) {
	if fh.BaseHandler.handlerConfig.Rollover.RolloverInterval == "" {
		fh.BaseHandler.handlerConfig.Rollover.RolloverInterval = defaultRolloverTime
	}
	fh.rotateInterval, err = utils.ParseSeconds(fh.BaseHandler.handlerConfig.Rollover.RolloverInterval)
	if err != nil {
		return errors.Wrap(err, "parse rotate interval error")
	}

	return nil
}

func (fh *FileHandler) nextTimeRotate(interval int) time.Time {
	return time.Now().Add(time.Duration(interval) * time.Second)
}

func (fh *FileHandler) settingBackupTime() (err error) {
	if fh.BaseHandler.handlerConfig.Rollover.BackupTime == "" {
		fh.BaseHandler.handlerConfig.Rollover.BackupTime = defaultBackupTime
	}
	fh.backupTime, err = utils.ParseSeconds(fh.BaseHandler.handlerConfig.Rollover.BackupTime)
	if err != nil {
		return errors.Wrap(err, "parse backup time error")
	}
	return nil

}

func (fh *FileHandler) settingBackupCount() error {
	if fh.BaseHandler.handlerConfig.Rollover.BackupCount <= 0 {
		fh.BaseHandler.handlerConfig.Rollover.BackupCount = defaultBackupCount
	}
	fh.backupCount = fh.BaseHandler.handlerConfig.Rollover.BackupCount
	return nil
}

func (fh *FileHandler) settingFileWriter() (err error) {
	if err = os.MkdirAll(fh.fileDir, os.ModePerm); err != nil {
		return errors.Wrap(err, "create dir error")
	}

	fh.fileWriter, err = os.OpenFile(fh.filePath, fileFlag, fileMode)
	if err != nil {
		return errors.Wrap(err, "open file error")
	}
	return nil
}

func (fh *FileHandler) settingWritenSize() error {
	if fileInfo, err := os.Stat(fh.filePath); err != nil {
		return errors.Wrap(err, "file stat error")
	} else {
		fh.writenSize = fileInfo.Size()
	}

	return nil
}
