package main

import (
	_ "embed"
	"flag"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	logLevel                 = "info"
	listenAddr               = "127.0.0.1:8092"
	rawLogOutput             = ""
	replayLog                = false
	dedupLaunchd             = false
	mergeSmallFilesThreshold = 32 * 1024 * 1024 // 32MB
	mergeMinStaleTime        = 60 * 60          // 1h
)

func main() {
	flag.StringVar(&logLevel, "l", logLevel, "Log level")
	flag.StringVar(&listenAddr, "a", listenAddr, "Listen address")
	flag.StringVar(&rawLogOutput, "o", rawLogOutput, "Save raw dtrace output to a file")
	flag.BoolVar(&replayLog, "r", replayLog, "Replay previous log file before collecting new data")
	flag.BoolVar(&dedupLaunchd, "d", dedupLaunchd, "Deduplicate launchd processes. Since launchd writes on behalf of other processes if the same file is written by another process it will be counted twice. This option removes the duplicated entries from launchd.")
	flag.IntVar(&mergeSmallFilesThreshold, "t", mergeSmallFilesThreshold, "Merge R/W smaller than this number of bytes into a single entry. This is useful for processes that write to many small files to save memory and make the output more readable.")
	flag.IntVar(&mergeMinStaleTime, "s", mergeMinStaleTime, "Only merge entries that are not updated than this number of seconds.")
	flag.Parse()

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Fatalf("Failed to parse log level: %s", err)
	}
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	processStat = NewProcessStat()

	logrus.Infof("Dedup launchd: %t", dedupLaunchd)
	logrus.Infof("Merge small files threshold: %d bytes", mergeSmallFilesThreshold)
	logrus.Infof("Merge min stale time: %d seconds", mergeMinStaleTime)

	if replayLog {
		logrus.Infof("Replaying log file: %s", rawLogOutput)
		if err := replayLogFile(); err != nil {
			logrus.Fatalf("Failed to replay log file: %s", err)
		}
	}

	go func() {
		sleepTime := 60
		if mergeMinStaleTime < sleepTime {
			sleepTime = mergeMinStaleTime
		}
		for {
			time.Sleep(time.Duration(sleepTime) * time.Second)
			processStat.MergeSmallEntries()
		}
	}()

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
