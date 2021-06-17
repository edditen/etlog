package main

import (
	"github.com/EdgarTeng/etlog"
	"log"
)

func main() {
	logger, err := etlog.NewDefaultLogger("example/log.yaml")
	if err != nil {
		log.Fatalf("err: %+v", err)
	}
	etlog.SetDefaultLog(logger)
	etlog.Log.Info("hello")
	etlog.Log.Warn("world")
}
