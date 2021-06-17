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
	etlog.Log.Debug("hello")
	etlog.Log.Info("hello")
	etlog.Log.Data("hello")
	etlog.Log.Warn("world")
	etlog.Log.Error("world")
	etlog.Log.Fatal("world")
}
