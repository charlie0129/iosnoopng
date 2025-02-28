package main

import (
	"github.com/VictoriaMetrics/metrics"
)

func formatReadCounterName(processName string) string {
	return "iosnoopng_read_bytes_total{process_name=\"" + processName + "\"}"
}

func formatWriteCounterName(processName string) string {
	return "iosnoopng_written_bytes_total{process_name=\"" + processName + "\"}"
}

func CountRead(processName string, reads uint64) {
	c := metrics.GetOrCreateCounter(formatReadCounterName(processName))
	c.AddInt64(int64(reads))
}

func CountWrite(processName string, writes uint64) {
	c := metrics.GetOrCreateCounter(formatWriteCounterName(processName))
	c.AddInt64(int64(writes))
}
