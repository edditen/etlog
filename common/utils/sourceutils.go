package utils

import (
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

func ShortSourceLoc(skip int) (fileName string, line int, funcName string, ok bool) {
	var _file string
	var _funcName string
	if _file, line, _funcName, ok = SourceLoc(skip); !ok {
		return "", 0, "", false
	}
	fileName = LastSubstring(_file, "/")
	funcName = LastSubstring(_funcName, "/")
	return
}
