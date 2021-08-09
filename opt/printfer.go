package opt

import (
	"github.com/edditen/etlog/common/utils"
	"log"
	"os"
)

var (
	_errLog  Printfer = NewInternalLog(log.New(os.Stderr, "error:", log.LstdFlags), 10, 1000)
	_infoLog Printfer = NewInternalLog(log.New(os.Stdout, "info:", log.LstdFlags), 10, 1000)
)

type Printfer interface {
	Printf(format string, v ...interface{})
}

type internalLog struct {
	rl *utils.RateLimiter
	p  Printfer
}

func NewInternalLog(p Printfer, rate, interval int64) *internalLog {
	return &internalLog{
		p:  p,
		rl: utils.NewRateLimiter(rate, interval),
	}
}

func (il *internalLog) Printf(format string, v ...interface{}) {
	if il.p == nil {
		return
	}
	if il.rl == nil {
		il.p.Printf(format, v...)
		return
	}

	if !il.rl.Allowable() {
		return
	}
	il.p.Printf(format, v...)
}

func SetErrLog(errLog Printfer) {
	_errLog = errLog
}

func SetInfoLog(infoLog Printfer) {
	_infoLog = infoLog
}

func GetErrLog() Printfer {
	return _errLog
}

func GetInfoLog() Printfer {
	return _infoLog
}
