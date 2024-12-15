package main

import (
	_ "embed"
	"flag"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	logLevel     = "info"
	listenAddr   = "127.0.0.1:8092"
	rawLogOutput = ""
)

func main() {
	flag.StringVar(&logLevel, "l", logLevel, "Log level")
	flag.StringVar(&listenAddr, "a", listenAddr, "Listen address")
	flag.StringVar(&rawLogOutput, "o", rawLogOutput, "Save raw dtrace output to a file")
	flag.Parse()

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Fatalf("failed to parse log level: %s", err)
	}
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	processStat = NewProcessStat()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		runIosnoop()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		httpServer()
	}()

	// TODO: add prometheus metrics

	wg.Wait()
}
