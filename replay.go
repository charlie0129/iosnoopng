package main

import (
	"bufio"
	"os"

	"github.com/sirupsen/logrus"
)

func replayLogFile() error {
	if rawLogOutput == "" {
		logrus.Warn("No log file to replay")
		return nil
	}

	f, err := os.Open(rawLogOutput)
	if err != nil {
		logrus.Fatalf("Failed to open log file: %s", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		output, err := NewIosnoopOutput(line)
		if err != nil {
			logrus.Warnf("Failed to parse iosnoop output: %s", err)
			continue
		}

		processStat.Add(output)
	}

	if err := scanner.Err(); err != nil {
		logrus.Fatalf("Failed to scan log file: %s", err)
	}

	return nil
}
