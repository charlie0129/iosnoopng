package main

import (
	"strings"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/sirupsen/logrus"
)

var processStat *ProcessStat

type RWStat struct {
	Writes      uint64    `json:"writes"`
	Reads       uint64    `json:"reads"`
	LastUpdated time.Time `json:"lastUpdated"`
}

func (p *RWStat) ReadsHuman() string {
	return humanize.Bytes(p.Reads)
}

func (p *RWStat) WritesHuman() string {
	return humanize.Bytes(p.Writes)
}

func (p *RWStat) CopyFrom(in *RWStat) {
	p.Reads = in.Reads
	p.Writes = in.Writes
	p.LastUpdated = in.LastUpdated
}

type MetaStat struct {
	Events       uint64             `json:"events"`
	ProcessStats map[string]*RWStat `json:"processStats"`
}

type ProcessStat struct {
	mu *sync.RWMutex
	// exec <-> path <-> RWStat
	Details map[string]map[string]*RWStat `json:"processStats,omitempty"`
	Events  uint64                        `json:"events"`
}

func NewProcessStat() *ProcessStat {
	ret := &ProcessStat{
		mu:      &sync.RWMutex{},
		Details: make(map[string]map[string]*RWStat),
	}
	return ret
}

func (ps *ProcessStat) MergeSmallEntries() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	pathsToDelete := make(map[string][]string)
	pathsToMerge := make(map[string][]string)

	for exec, pathRW := range ps.Details {
		needDedup := dedupLaunchd && exec == "launchd"

		for path, rw := range pathRW {
			if !strings.HasPrefix(path, "<Smaller Files of") &&
				rw.Reads < uint64(mergeSmallFilesThreshold) &&
				rw.Writes < uint64(mergeSmallFilesThreshold) &&
				time.Since(rw.LastUpdated) > time.Duration(mergeMinStaleTime)*time.Second {
				// If this path is written by another process, we should remove it.
				if needDedup && ps.isPathExistsInNonLaunchd(path) {
					pathsToDelete[exec] = append(pathsToDelete[exec], path)
					continue
				}

				// Merge.
				pathsToMerge[exec] = append(pathsToMerge[exec], path)
				pathsToDelete[exec] = append(pathsToDelete[exec], path)
			}
		}
	}

	mergeSum := 0
	for exec, paths := range pathsToMerge {
		mergedName := "<Smaller Files of " + exec + ">"
		if _, ok := ps.Details[exec][mergedName]; !ok {
			ps.Details[exec][mergedName] = &RWStat{}
		}
		mergedRW := ps.Details[exec][mergedName]

		for _, path := range paths {
			rw := ps.Details[exec][path]
			mergedRW.Reads += rw.Reads
			mergedRW.Writes += rw.Writes
		}

		mergedRW.LastUpdated = time.Now()
		mergeSum += len(paths)
	}

	deleteSum := 0
	for exec, paths := range pathsToDelete {
		for _, path := range paths {
			delete(ps.Details[exec], path)
		}
		deleteSum += len(paths)
	}

	logrus.Infof("Merged %d entries, deleted %d entries", mergeSum, deleteSum)
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

	detail.LastUpdated = time.Now()
	if !in.IsRead {
		detail.Writes += in.OperationSize
	} else {
		detail.Reads += in.OperationSize
	}

	ps.Events++
}

// Must lock before calling this function.
func (ps *ProcessStat) isPathExistsInNonLaunchd(path string) bool {
	for exec, pathRW := range ps.Details {
		if exec == "launchd" {
			continue
		}
		if _, ok := pathRW[path]; ok {
			return true
		}
	}

	return false
}

func (ps *ProcessStat) GetMeta() *MetaStat {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	ret := make(map[string]*RWStat, len(ps.Details))
	for exec, pathRW := range ps.Details {
		needDedup := dedupLaunchd && exec == "launchd"
		ret[exec] = &RWStat{}
		for path, rw := range pathRW {
			if needDedup && ps.isPathExistsInNonLaunchd(path) {
				continue
			}
			ret[exec].Reads += rw.Reads
			ret[exec].Writes += rw.Writes
			if rw.LastUpdated.After(ret[exec].LastUpdated) {
				ret[exec].LastUpdated = rw.LastUpdated
			}
		}
	}

	return &MetaStat{
		Events:       ps.Events,
		ProcessStats: ret,
	}
}

func (ps *ProcessStat) GetDetailsByProcess(processName string) map[string]*RWStat {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	needDedup := dedupLaunchd && processName == "launchd"

	ret := make(map[string]*RWStat, len(ps.Details[processName]))

	for path, rw := range ps.Details[processName] {
		if needDedup && ps.isPathExistsInNonLaunchd(path) {
			continue
		}
		newRW := &RWStat{}
		newRW.CopyFrom(rw)
		ret[path] = newRW
	}

	return ret
}

func (ps *ProcessStat) GetAll() map[string]map[string]*RWStat {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	ret := make(map[string]map[string]*RWStat, len(ps.Details))
	for exec, pathRW := range ps.Details {
		needDedup := dedupLaunchd && exec == "launchd"
		ret[exec] = make(map[string]*RWStat, len(pathRW))
		for path, rw := range pathRW {
			if needDedup && ps.isPathExistsInNonLaunchd(path) {
				continue
			}
			newRW := &RWStat{}
			newRW.CopyFrom(rw)
			ret[exec][path] = newRW
		}
	}

	return ret
}
