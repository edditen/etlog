package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/EdgarTeng/etlog"
)

func main() {
	//RunAll()
	RunRotate()
	log.Println("done")

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
	// etlog.Log.Data("hello")
	etlog.Log.Warn("world")
	etlog.Log.Error("world")
	// etlog.Log.Fatal("world")
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
	etlog.Log.WithField("beginTime", time.Now()).Info("start test")

	endTime := time.Now().Add(30 * time.Second)
	wg := new(sync.WaitGroup)
	for i := 0; i < 10; i++ {
		index := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				etlog.Log.
					WithMarker("trace").
					WithField("key", "word").
					WithField("index", index).
					WithField("now", time.Now()).
					Data("this is message")

				if time.Now().After(endTime) {
					break
				}
			}
		}()
	}
	wg.Wait()
	etlog.Log.WithField("endTime", time.Now()).Info("complete test")
	time.Sleep(10 * time.Second)

}
