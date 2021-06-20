package main

import (
	"fmt"
	"github.com/EdgarTeng/etlog"
	"log"
	"time"
)

func main() {
	RunAll()
}

func RunAll() {
	logger, err := etlog.NewDefaultLogger(etlog.SetConfigPath("example/log.yaml"))
	//logger, err := etlog.NewDefaultLogger()
	if err != nil {
		log.Fatalf("err: %+v", err)
	}
	etlog.SetDefaultLog(logger)
	etlog.Log.Debug("hello")
	etlog.Log.Info("hello")
	etlog.Log.Info("world")
	etlog.Log.Data("hello")
	etlog.Log.Warn("world")
	etlog.Log.Error("world")
	etlog.Log.Fatal("world")
	etlog.Log.WithError(fmt.Errorf("oops")).
		WithField("key", "word").
		WithField("now", time.Now()).
		Error("something wrong happened")

	for i := 1; i < 10; i++ {
		etlog.Log.WithError(fmt.Errorf("oops")).
			WithField("key", "word").
			WithField("now", time.Now()).
			Info("something wrong happened")
	}
}

func RunRotate() {
	logger, err := etlog.NewDefaultLogger(etlog.SetConfigPath("example/log.yaml"))
	//logger, err := etlog.NewDefaultLogger()
	if err != nil {
		log.Fatalf("err: %+v", err)
	}
	etlog.SetDefaultLog(logger)

	endTime := time.Now().Add(30 * time.Second)
	for {
		etlog.Log.
			WithField("key", "word").
			WithField("now", time.Now()).
			Info("this is message")

		if time.Now().After(endTime) {
			break
		}
	}
	log.Println("done")

}
