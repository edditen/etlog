package utils

import (
	"fmt"
	"runtime"
)

func SourceLoc(skip int) (file string, line int, funcName string, ok bool) {
	var pc uintptr
	if skip > 16 {
		skip = 16
	}
	if pc, file, line, ok = runtime.Caller(skip); ok {
		fun := runtime.FuncForPC(pc)
		funcName = fun.Name()
	}
	return
}

func ShortSourceLoc(skip int) (line string, funcName string, ok bool) {
	var _file string
	var _line int
	var _funcName string
	if _file, _line, _funcName, ok = SourceLoc(skip); !ok {
		return "", "", false
	}
	line = LastSubstring(fmt.Sprintf("%s:%d", _file, _line), "/")
	funcName = LastSubstring(_funcName, "/")
	return
}
