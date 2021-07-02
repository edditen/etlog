package config

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	stdlog "log"
)

var (
	DefaultConfig = NewConfig("")
)

func init() {
	handlerConfig := NewHandlerConfig()
	handlerConfig.Type = "std"
	handlerConfig.Levels = []string{"debug", "info", "data", "warn", "error", "fatal"}
	handlerConfig.Message.Format = "simple"
	DefaultConfig.LogConf.Handlers = append(DefaultConfig.LogConf.Handlers, *handlerConfig)
	DefaultConfig.LogConf.Level = "debug"
	DefaultConfig.LogConf.Name = "default"
}

type Config struct {
	configPath string
	LogConf    *LogConfig
}

func NewConfig(configPath string) *Config {
	return &Config{
		configPath: configPath,
		LogConf:    NewLogConfig(),
	}
}

func (c *Config) Init() error {
	stdlog.Println("[Init] init log config, log config:", c.configPath)
	yamlFile, err := ioutil.ReadFile(c.configPath)
	if err != nil {
		stdlog.Printf("[Init] init log config read config file error: %+v\n", err)
		return errors.Wrap(err, fmt.Sprintf("read file %s error", c.configPath))
	}

	err = yaml.Unmarshal(yamlFile, c.LogConf)
	if err != nil {
		stdlog.Printf("[Init] init log config unmarshal config file error: %+v\n", err)
		return errors.Wrap(err, fmt.Sprintf("unmarshal file error"))
	}
	return nil
}

type LogConfig struct {
	Handlers []HandlerConfig `yaml:"handlers"`
	Level    string          `yaml:"level"`
	Name     string          `yaml:"name"`
}

func NewLogConfig() *LogConfig {
	return &LogConfig{
		Handlers: make([]HandlerConfig, 0),
	}
}

type HandlerConfig struct {
	Type     string          `yaml:"type"`
	Marker   string          `yaml:"marker"`
	Levels   []string        `yaml:"levels"`
	File     string          `yaml:"file"`
	Rollover *RolloverConfig `yaml:"rollover"`
	Sync     *SyncConfig     `yaml:"sync"`
	Message  *MessageConfig  `yaml:"message"`
}

func NewHandlerConfig() *HandlerConfig {
	return &HandlerConfig{
		Rollover: NewRolloverConfig(),
		Sync:     NewSyncConfig(),
		Message:  NewMessageConfig(),
	}
}

type RolloverConfig struct {
	RolloverInterval string `yaml:"rollover_interval"`
	RolloverSize     string `yaml:"rollover_size"`
	BackupCount      int    `yaml:"backup_count"`
	BackupTime       string `yaml:"backup_time"`
}

func NewRolloverConfig() *RolloverConfig {
	return &RolloverConfig{}
}

type SyncConfig struct {
	AsyncWrite    bool `yaml:"async_write"`
	FlushInterval int  `yaml:"flush_interval"`
	FlushSize     int  `yaml:"flush_size"`
	QueueSize     int  `yaml:"queue_size"`
}

func NewSyncConfig() *SyncConfig {
	return &SyncConfig{}
}

type MessageConfig struct {
	Format       string `yaml:"format"`
	FieldsFormat string `yaml:"fields_format"`
	MaxBytes     string `yaml:"max_bytes"`
	MetaOption   string `yaml:"meta_option"`
}

func NewMessageConfig() *MessageConfig {
	return &MessageConfig{}
}
