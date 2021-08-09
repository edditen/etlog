package archiver

import (
	"github.com/edditen/etlog/common/queue"
	"github.com/edditen/etlog/common/runnable"
	"github.com/edditen/etlog/common/utils"
	"github.com/pkg/errors"
	"os"
	"path"
	"time"
)

const (
	defaultBackupExt      = ".zip"
	defaultBackupDelay    = 5 * time.Second
	defaultDelayQueueSize = 1000
	defaultTimeout        = 5 * time.Second
)

type Archiver interface {
	runnable.Runnable
	Archive(sourceFile string) error
}

type ArchiverOpt func(*LogArchiver) error

type LogArchiver struct {
	backupDir   string
	backupExt   string
	backupDelay time.Duration
	delayQueue  *queue.DelayQueue
	exitC       chan interface{}
}

func NewLogArchiver(backupDir string, options ...ArchiverOpt) (*LogArchiver, error) {
	a := &LogArchiver{
		backupDir:   backupDir,
		backupExt:   defaultBackupExt,
		backupDelay: defaultBackupDelay,
		delayQueue:  queue.NewDelayQueue(defaultDelayQueueSize),
		exitC:       make(chan interface{}),
	}

	for _, opt := range options {
		if err := opt(a); err != nil {
			return nil, err
		}
	}

	return a, nil
}

func (la *LogArchiver) Init() error {
	if err := os.MkdirAll(la.backupDir, os.ModePerm); err != nil {
		return errors.Wrap(err, "create archive dir error")
	}
	go la.Run()
	return nil
}

func (la *LogArchiver) Run() error {
	for {
		if val, ok := la.delayQueue.Take(defaultTimeout).(string); ok {
			_ = la.archive(val)
		}
		if la.isDown() && la.delayQueue.Len() == 0 {
			// exit util all data handled
			break
		}

	}
	return nil
}

func (la *LogArchiver) Shutdown() {
	close(la.exitC)
}

func (la *LogArchiver) isDown() bool {
	select {
	case <-la.exitC:
		return true
	default:
		return false
	}
}

func (la *LogArchiver) Archive(sourceFile string) error {
	if la.isDown() {
		return errors.New("log cleaner already shutdown")
	}

	expiredAt := time.Now().Add(la.backupDelay)
	if err := la.delayQueue.Offer(sourceFile, expiredAt); err != nil {
		return err
	}
	return nil
}

func (la *LogArchiver) archive(sourceFile string) error {
	archiveName := path.Base(sourceFile) + la.backupExt
	archiveFile := path.Join(la.backupDir, archiveName)
	if err := utils.ZipCompress(sourceFile, archiveFile); err != nil {
		return errors.Wrap(err, "archive error")
	}
	return la.removeSource(sourceFile)
}

func (la *LogArchiver) removeSource(sourceFile string) error {
	if err := os.Remove(sourceFile); err != nil {
		return errors.Wrap(err, "remove source file error")
	}
	return nil
}

func SetBackupExt(ext string) ArchiverOpt {
	return func(archiver *LogArchiver) error {
		archiver.backupExt = ext
		return nil
	}
}

func SetBackupDelay(delay time.Duration) ArchiverOpt {
	return func(archiver *LogArchiver) error {
		archiver.backupDelay = delay
		return nil
	}
}
