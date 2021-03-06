package handler

import (
	"fmt"
	"github.com/edditen/etlog/common/bufferpool"
	"github.com/edditen/etlog/handler/archiver"
	"github.com/edditen/etlog/handler/cleaner"
	"github.com/edditen/etlog/opt"
	"io/fs"
	"math"
	"os"
	"path"
	"sync"
	"sync/atomic"
	"time"

	"github.com/edditen/etlog/common/utils"
	"github.com/edditen/etlog/config"
	"github.com/edditen/etlog/core"
	"github.com/pkg/errors"
)

const (
	fileFlag                         = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	fileMode             fs.FileMode = 0644
	backupTimeFormat                 = "2006-01-02.150405"
	defaultLogSize                   = "10G"
	defaultRolloverTime              = "1d"
	defaultBackupTime                = "365d"
	defaultBackupCount               = math.MaxInt32
	defaultQueueSize                 = 8192
	defaultFlushInterval             = 100
	defaultFlushSize                 = 256
)

type LogEntries = []*core.LogEntry

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
	writtenSize    int64
	rotateAt       time.Time
	rotateLock     *sync.RWMutex
	flushLock      *sync.Mutex
	asyncWrite     bool
	queueSize      int
	flushInterval  int
	flushSize      int
	entryC         chan *core.LogEntry
	entryBuf       []*core.LogEntry
	ticker         *time.Ticker
	asyncMutex     *sync.RWMutex
	queueFull      chan bool
	cleaner        cleaner.Cleaner
	archiver       archiver.Archiver
}

func NewFileHandler(conf *config.HandlerConfig) *FileHandler {
	return &FileHandler{
		BaseHandler: NewBaseHandler(conf),
		rotateLock:  new(sync.RWMutex),
		flushLock:   new(sync.Mutex),
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

	if err := fh.settingSync(); err != nil {
		return err
	}

	if err := fh.settingChan(); err != nil {
		return err
	}

	if err := fh.settingCleaner(); err != nil {
		return err
	}

	if err := fh.settingArchiver(); err != nil {
		return err
	}

	return nil
}

func (fh *FileHandler) Handle(entry *core.LogEntry) (err error) {
	if !fh.BaseHandler.MarkerMatched(entry.Marker) {
		return nil
	}
	if !fh.BaseHandler.Contains(entry.Level) {
		return nil
	}

	if fh.asyncWrite {
		select {
		case fh.entryC <- entry:
		default:
			err = fh.syncHandle(entry)
			fh.notifyFull()
		}
		return err
	}
	return fh.syncHandle(entry)
}

func (fh *FileHandler) syncHandle(entry *core.LogEntry) error {
	buf := fh.BaseHandler.formatter.Format(entry)
	defer buf.Free()

	if err := fh.syncFlush(buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (fh *FileHandler) syncFlush(bs []byte) error {
	if fh.shouldCreateFile() {
		if err := fh.createFileWithLock(); err != nil {
			return err
		}
	}

	if fh.shouldRotate() {
		if err := fh.Rotate(); err != nil {
			return err
		}
	}

	if err := fh.Flush(bs); err != nil {
		return err
	}
	return nil
}

func (fh *FileHandler) shouldCreateFile() bool {
	return fh.fileWriter == nil
}

func (fh *FileHandler) createFileWithLock() error {
	fh.rotateLock.Lock()
	defer fh.rotateLock.Unlock()

	return fh.createFile()
}

func (fh *FileHandler) createFile() error {
	// double check
	if !fh.shouldCreateFile() {
		return nil
	}

	if err := fh.settingFileWriter(); err != nil {
		return err
	}

	if err := fh.settingWrittenSize(); err != nil {
		return err
	}

	fh.rotateAt = fh.nextTimeRotate(fh.rotateInterval)
	return nil
}

func (fh *FileHandler) Flush(bs []byte) error {
	fh.rotateLock.RLock()
	defer fh.rotateLock.RUnlock()

	fh.flushLock.Lock()
	defer fh.flushLock.Unlock()

	if _, err := fh.fileWriter.Write(bs); err != nil {
		return errors.Wrap(err, "write file error")
	}

	atomic.AddInt64(&fh.writtenSize, int64(len(bs)))
	return nil
}

func (fh *FileHandler) Rotate() error {
	fh.rotateLock.Lock()
	defer fh.rotateLock.Unlock()

	if !fh.shouldRotate() {
		return nil
	}

	backupName, err := fh.backup()
	if err != nil {
		return err
	}

	if err := fh.createFile(); err != nil {
		return err
	}

	go fh.postRotate(backupName)

	return nil
}

func (fh *FileHandler) backup() (string, error) {
	fh.closeFileWriter()

	backupName := fh.genBackupFileName()
	if err := os.Rename(fh.filePath, backupName); err != nil {
		return "", errors.Wrap(err, "rotate file error")
	}
	return backupName, nil
}

func (fh *FileHandler) postRotate(backupName string) {
	// archive
	if err := fh.archiver.Archive(backupName); err != nil {
		opt.GetErrLog().Printf("archive files err: %+v\n", err)
	}

	//clean
	if err := fh.cleaner.Clean(); err != nil {
		opt.GetErrLog().Printf("clean backup files err: %+v\n", err)
	}

}

func (fh *FileHandler) shouldRotate() bool {
	if time.Now().After(fh.rotateAt) {
		return true
	}
	if fh.writtenSize > int64(fh.rotateSize) {
		return true
	}
	return false
}

func (fh *FileHandler) genBackupFileName() string {
	filename := fh.fileName[:len(fh.fileName)-len(fh.fileExt)]
	t := time.Now().Format(backupTimeFormat)
	filename = fmt.Sprintf("%s.%s%s", filename, t, fh.fileExt)
	return path.Join(fh.fileDir, filename)
}

func (fh *FileHandler) closeFileWriter() {
	_ = fh.fileWriter.Close()
	fh.fileWriter = nil
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

func (fh *FileHandler) settingWrittenSize() error {
	if fileInfo, err := os.Stat(fh.filePath); err != nil {
		return errors.Wrap(err, "file stat error")
	} else {
		fh.writtenSize = fileInfo.Size()
	}

	return nil
}

func (fh *FileHandler) settingSync() error {

	if fh.BaseHandler.handlerConfig.Sync.QueueSize <= 0 {
		fh.BaseHandler.handlerConfig.Sync.QueueSize = defaultQueueSize
	}
	if fh.BaseHandler.handlerConfig.Sync.FlushInterval <= 0 {
		fh.BaseHandler.handlerConfig.Sync.FlushInterval = defaultFlushInterval
	}
	if fh.BaseHandler.handlerConfig.Sync.FlushSize <= 0 {
		fh.BaseHandler.handlerConfig.Sync.FlushSize = defaultFlushSize
	}

	fh.queueSize = fh.BaseHandler.handlerConfig.Sync.QueueSize
	fh.flushInterval = fh.BaseHandler.handlerConfig.Sync.FlushInterval
	fh.flushSize = fh.BaseHandler.handlerConfig.Sync.FlushSize
	fh.asyncWrite = fh.BaseHandler.handlerConfig.Sync.AsyncWrite

	return nil
}

func (fh *FileHandler) settingChan() error {
	if !fh.asyncWrite {
		return nil
	}

	fh.asyncMutex = new(sync.RWMutex)
	fh.queueFull = make(chan bool)
	fh.entryC = make(chan *core.LogEntry, fh.queueSize)
	fh.entryBuf = make([]*core.LogEntry, 0)
	fh.ticker = time.NewTicker(time.Duration(fh.flushInterval) * time.Millisecond)

	go fh.runChan()

	return nil
}

func (fh *FileHandler) settingCleaner() (err error) {
	duration := time.Duration(fh.backupTime) * time.Second
	baseName := fh.fileName[:len(fh.fileName)-len(fh.fileExt)]

	fh.cleaner, err = cleaner.NewLogCleaner(
		fh.fileDir, baseName,
		cleaner.SetBackupCount(fh.backupCount),
		cleaner.SetBackupDuration(duration),
	)
	if err != nil {
		return errors.Wrap(err, "create log cleaner error")
	}

	if err = fh.cleaner.Init(); err != nil {
		return errors.Wrap(err, "init log cleaner error")
	}

	return nil
}

func (fh *FileHandler) settingArchiver() (err error) {

	fh.archiver, err = archiver.NewLogArchiver(fh.fileDir)
	if err != nil {
		return errors.Wrap(err, "create log archiver error")
	}

	if err = fh.archiver.Init(); err != nil {
		return errors.Wrap(err, "init log archiver error")
	}

	return nil
}

func (fh *FileHandler) notifyFull() {
	select {
	case fh.queueFull <- true:
	default:
	}
}

func (fh *FileHandler) runChan() {
	for {
		select {
		case logEntry := <-fh.entryC:
			fh.appendLogEntry(logEntry)
		case <-fh.ticker.C:
			fh.handleLogEntry()
		case <-fh.queueFull:
			fh.handleLogEntry()
		}
	}
}

func (fh *FileHandler) appendLogEntry(entry *core.LogEntry) {
	fh.asyncMutex.RLock()
	defer fh.asyncMutex.RUnlock()

	fh.entryBuf = append(fh.entryBuf, entry)
}

func (fh *FileHandler) handleLogEntry() {
	fh.asyncMutex.Lock()
	defer fh.asyncMutex.Unlock()

	if len(fh.entryBuf) == 0 {
		return
	}

	blocks := utils.CalculateBlocks(len(fh.entryBuf), fh.flushSize)
	for i := 0; i < blocks; i++ {
		buf := bufferpool.Borrow()

		for j := 0; j < fh.flushSize && i*blocks+j < len(fh.entryBuf); j++ {
			entry := fh.entryBuf[i*blocks+j]
			b := fh.formatter.Format(entry)
			buf.AppendBytes(b.Bytes())
			b.Free()
		}

		if err := fh.syncFlush(buf.Bytes()); err != nil {
			opt.GetErrLog().Printf("sync flush log err: %+v\n", err)
		}
		buf.Free()
	}

	fh.entryBuf = fh.entryBuf[:0]
}
