package main

import (
	"fmt"
	"github.com/EdgarTeng/etlog"
	"log"
	"time"
)

func main() {
	logger, err := etlog.NewDefaultLogger(etlog.SetConfigPath("example/log.yaml"))
	//logger, err := etlog.NewDefaultLogger()
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
	etlog.Log.WithError(fmt.Errorf("oops")).
		WithField("key", "word").
		WithField("now", time.Now()).
		Error("something wrong happened")
}
