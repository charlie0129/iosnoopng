package main

import (
	"sync"

	"github.com/dustin/go-humanize"
)

var processStat *ProcessStat

type RWStat struct {
	Writes uint64 `json:"writes"`
	Reads  uint64 `json:"reads"`
}

func (p *RWStat) ReadsHuman() string {
	return humanize.Bytes(p.Reads)
}

func (p *RWStat) WritesHuman() string {
	return humanize.Bytes(p.Writes)
}

type ProcessStat struct {
	mu      *sync.RWMutex
	Meta    map[string]*RWStat            `json:"processStatsMeta,omitempty"`
	Details map[string]map[string]*RWStat `json:"processStats,omitempty"`
	Events  uint64                        `json:"events"`
}

func NewProcessStat() *ProcessStat {
	ret := &ProcessStat{
		mu:      &sync.RWMutex{},
		Details: make(map[string]map[string]*RWStat),
		Meta:    make(map[string]*RWStat),
	}
	return ret
}

func (ps *ProcessStat) Add(in *IosnoopOutput) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, ok := ps.Details[in.ExecName]; !ok {
		ps.Details[in.ExecName] = make(map[string]*RWStat)
	}
	if _, ok := ps.Details[in.ExecName][in.Path]; !ok {
		ps.Details[in.ExecName][in.Path] = &RWStat{}
	}

	detail := ps.Details[in.ExecName][in.Path]

	if _, ok := ps.Meta[in.ExecName]; !ok {
		ps.Meta[in.ExecName] = &RWStat{}
	}
	meta := ps.Meta[in.ExecName]

	if !in.IsRead {
		detail.Writes += in.OperationSize
		meta.Writes += in.OperationSize
	} else {
		detail.Reads += in.OperationSize
		meta.Reads += in.OperationSize
	}

	ps.Events++
}

func (ps *ProcessStat) GetDetailsByProcess(processName string) map[string]*RWStat {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	ret := make(map[string]*RWStat, len(ps.Details[processName]))
	for k, v := range ps.Details[processName] {
		ret[k] = v
	}

	return ret
}

func (ps *ProcessStat) GetAll() map[string]map[string]*RWStat {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	ret := make(map[string]map[string]*RWStat, len(ps.Details))
	for k, v := range ps.Details {
		ret[k] = make(map[string]*RWStat, len(v))
		for k2, v2 := range v {
			ret[k][k2] = v2
		}
	}

	return ret
}

func (ps *ProcessStat) GetMeta() *ProcessStat {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	ret := make(map[string]*RWStat, len(ps.Meta))
	for k, v := range ps.Meta {
		ret[k] = v
	}

	return &ProcessStat{
		Events: ps.Events,
		Meta:   ret,
	}
}
