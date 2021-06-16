package example

import (
	"encoding/json"
	"github.com/EdgarTeng/etlog"
	stdlog "log"
	"testing"
)

func TestConfig(t *testing.T) {
	conf := etlog.NewConfig("log.yaml")
	if err := conf.Init(); err != nil {
		stdlog.Fatalf("%+v", err)
	}
	b, err := json.Marshal(conf.LogConf)
	if err != nil {
		stdlog.Fatalf("%+v", err)
	}
	stdlog.Println(string(b))
}
