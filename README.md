# log

A wrapper library based on `go.uber.org/zap`.

## Installation

`go get -u github.com/magicedy/log`

## Quick Start

```go
package main

import (
	"github.com/magicedy/log"
	"go.uber.org/zap"
	"sync"
)

var once sync.Once

func setupLogger() {
	once.Do(func() {
		fileConf := log.FileLogConfig{RootPath: "logs", Filename: "log.jsonl"}
		conf := &log.Config{Level: "debug", Stdout: true, Format: "json", File: fileConf}
		l, p, err := log.InitLogger(conf)
		if err == nil {
			log.ReplaceGlobals(l, p)
		} else {
			log.Fatal("initialize logger error", zap.Error(err))
		}
	})
}

func main() {
	setupLogger()

    log.Print("this is info message")
	log.Println("Hello", "World")
	log.Debug("this is debug message")
	log.Info("this is info message")
	log.Infow("this is info message with fileds", "age", 37, "agender", "man")
	log.Warn("this is warn message")
	log.Error("this is error message")
}
```

```
{"level":"INFO","msg":"this is info message"}
{"level":"INFO","msg":"Hello World"}
{"level":"DEBUG","msg":"this is debug message"}
{"level":"INFO","msg":"this is info message"}
{"level":"INFO","msg":"this is info message with fileds","age":37,"agender":"man"}
{"level":"WARN","msg":"this is warn message"}
{"level":"ERROR","msg":"this is error message"}
```

