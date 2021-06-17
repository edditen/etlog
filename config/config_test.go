package config

import (
	"encoding/json"
	stdlog "log"
	"testing"
)

func TestConfig(t *testing.T) {
	conf := NewConfig("../example/log.yaml")
	if err := conf.Init(); err != nil {
		stdlog.Fatalf("%+v", err)
	}
	b, err := json.Marshal(conf.LogConf)
	if err != nil {
		stdlog.Fatalf("%+v", err)
	}
	stdlog.Println(string(b))
}
