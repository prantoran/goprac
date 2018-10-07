package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/sebest/logrusly"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

var logglyToken string = "YOUR_LOGGLY_TOKEN"

func main() {
	log := logrus.New()
	hook := logrusly.NewLogglyHook(logglyToken, "www.hostname.com", logrus.WarnLevel, "tag1", "tag2")
	log.Hooks.Add(hook)

	log.WithFields(logrus.Fields{
		"name": "joe",
		"age":  42,
	}).Error("Hello world!")

	// Flush is automatic for panic/fatal
	// Just make sure to Flush() before exiting or you may loose up to 5 seconds
	// worth of messages.
	hook.Flush()

	glog.Info("Prepare to repel boarders")

	// glog.Fatalf("Initialization failed: %s", errors.New("yo"))
	// zap
	fmt.Println("zap***")
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar0 := logger.Sugar()
	sugar0.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", "prantoran.me",
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar0.Infof("Failed to fetch URL: %s", "prantoran.me")

	fmt.Println("***********")
	sugar := zap.NewExample().Sugar()
	defer sugar.Sync()
	sugar.Infow("failed to fetch URL",
		"url", "http://example.com",
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("failed to fetch URL: %s", "http://example.com")

	fmt.Println("**********")
	rawJSON := []byte(`{
		"level": "debug",
		"encoding": "json",
		"outputPaths": ["stdout", "/tmp/logs"],
		"errorOutputPaths": ["stderr"],
		"initialFields": {"foo": "bar"},
		"encoderConfig": {
		  "messageKey": "message",
		  "levelKey": "level",
		  "levelEncoder": "lowercase"
		}
	  }`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Info("logger construction succeeded")
}
