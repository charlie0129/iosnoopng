package main

import (
	"fmt"
	"strconv"
	"strings"
)

type IosnoopOutput struct {
	ExecName      string `json:"execName"`
	OperationSize uint64 `json:"operationSize"`
	IsRead        bool   `json:"isRead"`
	Path          string `json:"path"`
}

func NewIosnoopOutput(in string) (*IosnoopOutput, error) {
	fields := strings.Split(strings.TrimSpace(in), "|")
	if len(fields) != 4 {
		return nil, fmt.Errorf("expected 4 fields, got %d", len(fields))
	}

	ret := &IosnoopOutput{}
	var err error

	ret.ExecName = fields[0]

	if fields[1] == "0" {
		ret.IsRead = false
	} else {
		ret.IsRead = true
	}

	ret.OperationSize, err = strconv.ParseUint(fields[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse operation size: %w", err)
	}

	ret.Path = fields[3]
	if strings.HasPrefix(ret.Path, "??") {
		ret.Path = strings.TrimPrefix(ret.Path, "??")
	}

	return ret, nil
}
