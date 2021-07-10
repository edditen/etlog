package cleaner

import (
	"fmt"
	"github.com/EdgarTeng/etlog/common/runnable"
	"github.com/EdgarTeng/etlog/common/utils"
	"github.com/EdgarTeng/etlog/opt"
	"github.com/pkg/errors"
	"math"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	defaultBackupExt      = ".zip"
	defaultBackupCount    = math.MaxInt32
	defaultBackupDuration = 100 * 365 * 24 * time.Hour //100 years
	defaultCheckInterval  = 10 * time.Minute
	defaultTimeFormat     = "2006-01-02.150405"
	defaultTimePattern    = regexp.MustCompile(".*([\\d]{4}-[\\d]{2}-[\\d]{2}\\.[\\d]{6}).*")
)

type FileInfo struct {
	FileDir    string
	Filename   string
	BackupTime time.Time
}

func (fi FileInfo) String() string {
	return fmt.Sprintf("{dir:%s, name:%s, time:%v}",
		fi.FileDir, fi.Filename, fi.BackupTime)
}

type Cleaner interface {
	runnable.Runnable
	Clean() error
}

type Option func(*LogCleaner) error

type LogCleaner struct {
	backupDir      string
	backupBaseName string
	backupExt      string
	backupCount    int
	backupDuration time.Duration
	checkInterval  time.Duration
	mutex          *sync.Mutex
	ticker         *time.Ticker
	exitC          chan interface{}
}

func NewLogCleaner(backupDir, backupBaseName string, options ...Option) (*LogCleaner, error) {
	fc := &LogCleaner{
		backupDir:      backupDir,
		backupBaseName: backupBaseName,
		backupCount:    defaultBackupCount,
		backupDuration: defaultBackupDuration,
		backupExt:      defaultBackupExt,
		checkInterval:  defaultCheckInterval,
		mutex:          new(sync.Mutex),
		exitC:          make(chan interface{}),
	}

	for _, option := range options {
		if err := option(fc); err != nil {
			return nil, err
		}
	}

	return fc, nil

}

func (lc *LogCleaner) Init() error {
	lc.ticker = time.NewTicker(lc.checkInterval)
	go lc.Run()
	return nil
}

func (lc *LogCleaner) Run() error {
	defer lc.ticker.Stop()
	for {
		select {
		case <-lc.ticker.C:
			_ = lc.Clean()
		case <-lc.exitC:
			break
		}
	}
	return nil
}

func (lc *LogCleaner) Shutdown() {
	close(lc.exitC)
}

func (lc *LogCleaner) Clean() error {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	cleanFiles, required := lc.shouldClean()
	if !required {
		return nil
	}

	if err := lc.removeFiles(cleanFiles); err != nil {
		return errors.Wrap(err, "remove file error")
	}
	return nil
}

func (lc *LogCleaner) removeFiles(files []FileInfo) error {
	for _, f := range files {
		if err := os.Remove(path.Join(f.FileDir, f.Filename)); err != nil {
			return err
		}
	}

	return nil
}

func (lc *LogCleaner) shouldClean() (files []FileInfo, required bool) {
	files = make([]FileInfo, 0)
	required = false

	matchedFiles := lc.listBackupFiles()
	if len(matchedFiles) == 0 {
		return
	}

	expiredFiles := lc.expiredFiles(matchedFiles)
	if len(expiredFiles) > 0 {
		files = append(files, expiredFiles...)
		required = true
	}

	remains := lc.deduplicateByFilename(matchedFiles, expiredFiles)

	cleanFiles := lc.gtBackupLimit(remains)

	if len(cleanFiles) > 0 {
		files = append(files, cleanFiles...)
		required = true
	}

	return files, required
}

func (lc *LogCleaner) listBackupFiles() []FileInfo {
	matchedFiles := make([]FileInfo, 0)
	err := filepath.Walk(lc.backupDir,
		func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			filename := info.Name()
			if filename == lc.backupBaseName {
				return nil
			}

			if !strings.HasPrefix(filename, lc.backupBaseName) ||
				!strings.HasSuffix(filename, lc.backupExt) {
				return nil
			}

			var ts string
			var matched bool
			if ts, matched = utils.GetFirstMatchedString(defaultTimePattern, filename); !matched {
				return nil
			}

			var backupTime time.Time
			if backupTime, err = time.Parse(defaultTimeFormat, ts); err != nil {
				return nil
			}

			matchedFiles = append(matchedFiles, FileInfo{
				FileDir:    path.Dir(filePath),
				Filename:   filename,
				BackupTime: backupTime,
			})

			return nil
		})
	if err != nil {
		opt.GetErrLog().Printf("list backup files err: %+v\n", err)
	}

	return matchedFiles
}

func (lc *LogCleaner) expiredFiles(files []FileInfo) []FileInfo {
	expired := make([]FileInfo, 0)
	now := time.Now()
	for _, info := range files {
		if now.Sub(info.BackupTime) > lc.backupDuration {
			expired = append(expired, info)
		}
	}
	return expired
}

// gtBackupLimit greater than backup limit,
// then keep the backup limit files, the oldest files will be returned
func (lc *LogCleaner) gtBackupLimit(files []FileInfo) []FileInfo {
	if len(files) <= lc.backupCount {
		return make([]FileInfo, 0)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].BackupTime.Before(files[j].BackupTime)
	})

	removeCount := len(files) - lc.backupCount
	return files[:removeCount]
}

func (lc *LogCleaner) deduplicateByFilename(fullList, sublist []FileInfo) []FileInfo {
	if len(fullList) == 0 || len(sublist) == 0 {
		return fullList
	}

	keys := make(map[string]bool, 0)
	for _, it := range sublist {
		keys[it.Filename] = true
	}

	resultList := make([]FileInfo, 0)
	for _, it := range fullList {
		if _, ok := keys[it.Filename]; !ok {
			resultList = append(resultList, it)
		}
	}
	return resultList
}

func SetBackupCount(backupCount int) Option {
	return func(cleaner *LogCleaner) error {
		cleaner.backupCount = backupCount
		return nil
	}
}

func SetBackupExt(backupExt string) Option {
	return func(cleaner *LogCleaner) error {
		cleaner.backupExt = backupExt
		return nil
	}
}

func SetCheckInterval(interval time.Duration) Option {
	return func(cleaner *LogCleaner) error {
		cleaner.checkInterval = interval
		return nil
	}
}

func SetBackupDuration(duration time.Duration) Option {
	return func(cleaner *LogCleaner) error {
		cleaner.backupDuration = duration
		return nil
	}
}
