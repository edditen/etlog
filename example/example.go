package main

import (
	"fmt"
	"github.com/EdgarTeng/etlog/core"
	"github.com/EdgarTeng/etlog/opt"
	"log"
	"sync"
	"time"

	"github.com/EdgarTeng/etlog"
)

func main() {
	etlog.Log.Info("start")
	logger, err := etlog.NewEtLogger(
		etlog.SetConfigPath("example/log.yaml"),
		etlog.SetPreLog(preLog()...),
		etlog.SetPostLog(postLog()...),
	)
	if err != nil {
		log.Fatalf("err: %+v", err)
	}
	etlog.SetDefaultLog(logger)

	RunAll()
	RunRotate()
	time.Sleep(10 * time.Second)
	log.Println("done")

}

func preLog() []opt.LogFunc {
	fns := make([]opt.LogFunc, 0)
	fns = append(fns, func(e *opt.LogE) {
		if e.Level == "ERROR" {
			log.Printf("received an error at %v, msg: %s, err: %+v, fields: %v\n",
				e.Time, e.Msg, e.Err, e.Fields)
		}
	})
	return fns
}

func postLog() []opt.LogFunc {
	fns := make([]opt.LogFunc, 0)
	fns = append(fns, func(e *opt.LogE) {
		if e.Level == "ERROR" {
			log.Printf("handled an error at %v, msg: %s, err: %+v, fields: %v\n",
				e.Time, e.Msg, e.Err, e.Fields)
		}
	})
	fns = append(fns, func(e *opt.LogE) {
		if e.Level == "ERROR" {
			panic("I'm a panic")
		}
	})

	return fns
}

func RunAll() {
	etlog.Log.Debug("hello")
	etlog.Log.Info("world")
	etlog.Log.Data("hello")
	etlog.Log.Warn("world")
	etlog.Log.Error("world")
	etlog.Log.Fatal("world")
	fmt.Println("DEBUG enable", etlog.Log.Enable(core.DEBUG))
	fmt.Println("INFO enable", etlog.Log.Enable(core.INFO))
	fmt.Println("DATA enable", etlog.Log.Enable(core.DATA))
	fmt.Println("WARN enable", etlog.Log.Enable(core.WARN))
	fmt.Println("ERROR enable", etlog.Log.Enable(core.ERROR))
	fmt.Println("FATAL enable", etlog.Log.Enable(core.FATAL))

	etlog.Log.WithError(fmt.Errorf("oops")).
		WithField("key", "word").
		WithField("now", time.Now()).
		Error("something wrong happened")
	etlog.Log.WithError(fmt.Errorf("oops")).
		WithField("key", "word").
		WithField("now", time.Now()).
		Error("something wrong happened")
	etlog.Log.WithFields(core.Fields{
		"abc": 123,
		"xyz": "hello world",
		"now": time.Now(),
	}).Info("test fields")

	for i := 1; i < 10; i++ {
		etlog.Log.WithError(fmt.Errorf("oops")).
			WithField("key", "word").
			WithField("now", time.Now()).
			Info("something wrong happened")
	}
}

func RunRotate() {

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
					WithMarkers("trace", "data").
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

}
