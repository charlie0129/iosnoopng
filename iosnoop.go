package main

import (
	"bufio"
	"context"
	_ "embed"
	"io"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

var iosnoopBin = `#pragma D option quiet
io:::start {
    printf("%s|%d|%d|%s\n", execname, args[0]->b_flags & B_READ, args[0]->b_bcount,args[2]->fi_pathname);
}
`

func runIosnoop() {
	cmd := exec.CommandContext(context.Background(), "/usr/sbin/dtrace", "-n", iosnoopBin)

	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		logrus.Fatalf("failed to get stdout pipe: %s", err)
	}

	cmdStderr, err := cmd.StderrPipe()
	if err != nil {
		logrus.Fatalf("failed to get stderr pipe: %s", err)
	}

	if err := cmd.Start(); err != nil {
		logrus.Fatalf("failed to start iosnoop: %s", err)
	}

	go readStderr(cmdStderr)
	readStdout(cmdStdout)

	if err := cmd.Wait(); err != nil {

		logrus.Fatalf("iosnoop failed: %s", err)
	}
}

func readStderr(stderr io.Reader) {
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		logrus.Warn(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		logrus.Fatalf("Failed to scan iosnoop stderr: %s", err)
	}
}

func readStdout(stdout io.Reader) {
	var fWriter io.Writer
	if rawLogOutput != "" {
		f, err := os.OpenFile(rawLogOutput, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logrus.Fatalf("Failed to open raw log output: %s", err)
		}
		defer f.Close()
		fWriter = f
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()

		if fWriter != nil {
			if _, err := fWriter.Write([]byte(line + "\n")); err != nil {
				logrus.Warnf("Failed to write to raw log output: %s", err)
			}
		}

		output, err := NewIosnoopOutput(line)
		if err != nil {
			logrus.Warnf("Failed to parse iosnoop output: %s", err)
			continue
		}

		logrus.WithField("output", output).Debug("Parsed iosnoop output")
		processStat.Add(output)
	}

	if err := scanner.Err(); err != nil {
		logrus.Fatalf("Failed to scan iosnoop output: %s", err)
	}
}
