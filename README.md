# etlog

**etlog** is a log component for go. 



## Features

- Support log level
- Stdout and file appender
- log markers support



## Quick start



1. Simple usage

   ```go
   etlog.Log.Info("Hello World")
   ```

It will output the log into stdout device(console log).

2. Using log config

```go
logger, err := etlog.NewEtLogger(etlog.SetConfigPath("log.yaml"))
if err != nil {
	panic(err)
}
logger.Debug("hello")
logger.Info("world")
```

the log config file please refer to:  [example/log.yaml](./example/log.yaml)

3. Using fieds

```go
etlog.Log.WithError(fmt.Errorf("oops")).
			WithField("key", "word").
			WithField("now", time.Now()).
			Error("something wrong happened")
```

the `WithField` method will help you print K-V fields into log.

4. Using markers

```go
etlog.Log.WithMarkers("trace").Data("hello world")
```

Because we support different log handler, to determine which handler the content will be output, we use `marker` to route it. Such   as the example, when we use `trace` as marker of log, then the content will be processed by handler marked as `trace`.

