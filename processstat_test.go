package main

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestDedupOff(t *testing.T) {
	dedupLaunchd = false

	processStat := NewProcessStat()

	processStat.Add(&IosnoopOutput{
		ExecName:      "launchd",
		OperationSize: 1024,
		IsRead:        false,
		Path:          "/1",
	})
	processStat.Add(&IosnoopOutput{
		ExecName:      "another",
		OperationSize: 1024,
		IsRead:        false,
		Path:          "/1",
	})

	out := processStat.GetAll()
	expected := map[string]map[string]*RWStat{
		"launchd": {
			"/1": {
				Reads:       0,
				Writes:      1024,
				LastUpdated: getLastUpdated(out["launchd"]["/1"]),
			},
		},
		"another": {
			"/1": {
				Reads:       0,
				Writes:      1024,
				LastUpdated: getLastUpdated(out["another"]["/1"]),
			},
		},
	}

	if !reflect.DeepEqual(out, expected) {
		t.Errorf("expected %s, got %s", toJSON(expected), toJSON(out))
	}
}

func TestDedupOn(t *testing.T) {
	dedupLaunchd = true

	processStat := NewProcessStat()

	processStat.Add(&IosnoopOutput{
		ExecName:      "launchd",
		OperationSize: 1024,
		IsRead:        false,
		Path:          "/1",
	})
	processStat.Add(&IosnoopOutput{
		ExecName:      "another",
		OperationSize: 1024,
		IsRead:        false,
		Path:          "/1",
	})

	out := processStat.GetAll()
	expected := map[string]map[string]*RWStat{
		"launchd": {},
		"another": {
			"/1": {
				Reads:       0,
				Writes:      1024,
				LastUpdated: getLastUpdated(out["another"]["/1"]),
			},
		},
	}

	if !reflect.DeepEqual(out, expected) {
		t.Errorf("expected %s, got %s", toJSON(expected), toJSON(out))
	}
}

func TestMergeDedupOff(t *testing.T) {
	mergeSmallFilesThreshold = 1024
	mergeMinStaleTime = 1
	dedupLaunchd = false

	processStat := NewProcessStat()

	processStat.Add(&IosnoopOutput{
		ExecName:      "launchd",
		OperationSize: 1024,
		IsRead:        false,
		Path:          "/1",
	})
	processStat.Add(&IosnoopOutput{
		ExecName:      "launchd",
		OperationSize: 512,
		IsRead:        false,
		Path:          "/2",
	})
	processStat.Add(&IosnoopOutput{
		ExecName:      "launchd",
		OperationSize: 512,
		IsRead:        false,
		Path:          "/3",
	})
	processStat.Add(&IosnoopOutput{
		ExecName:      "another",
		OperationSize: 1024,
		IsRead:        false,
		Path:          "/1",
	})
	processStat.Add(&IosnoopOutput{
		ExecName:      "another",
		OperationSize: 512,
		IsRead:        false,
		Path:          "/2",
	})

	out := processStat.GetAll()
	expected := map[string]map[string]*RWStat{
		"launchd": {
			"/1": {
				Reads:       0,
				Writes:      1024,
				LastUpdated: getLastUpdated(out["launchd"]["/1"]),
			},
			"/2": {
				Reads:       0,
				Writes:      512,
				LastUpdated: getLastUpdated(out["launchd"]["/2"]),
			},
			"/3": {
				Reads:       0,
				Writes:      512,
				LastUpdated: getLastUpdated(out["launchd"]["/3"]),
			},
		},
		"another": {
			"/1": {
				Reads:       0,
				Writes:      1024,
				LastUpdated: getLastUpdated(out["another"]["/1"]),
			},
			"/2": {
				Reads:       0,
				Writes:      512,
				LastUpdated: getLastUpdated(out["another"]["/2"]),
			},
		},
	}

	if !reflect.DeepEqual(out, expected) {
		t.Errorf("expected %s, got %s", toJSON(expected), toJSON(out))
	}

	time.Sleep(1100 * time.Millisecond)

	processStat.MergeSmallEntries()
	out = processStat.GetAll()
	expected = map[string]map[string]*RWStat{
		"launchd": {
			"/1": {
				Reads:       0,
				Writes:      1024,
				LastUpdated: getLastUpdated(out["launchd"]["/1"]),
			},
			"<Smaller Files of launchd>": {
				Reads:       0,
				Writes:      1024,
				LastUpdated: getLastUpdated(out["launchd"]["<Smaller Files of launchd>"]),
			},
		},
		"another": {
			"/1": {
				Reads:       0,
				Writes:      1024,
				LastUpdated: getLastUpdated(out["another"]["/1"]),
			},
			"<Smaller Files of another>": {
				Reads:       0,
				Writes:      512,
				LastUpdated: getLastUpdated(out["another"]["<Smaller Files of another>"]),
			},
		},
	}

	if !reflect.DeepEqual(out, expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", toJSON(expected), toJSON(out))
	}

	// Test that the merged entries are not merged again.
	time.Sleep(1100 * time.Millisecond)

	processStat.MergeSmallEntries()
	out = processStat.GetAll()
	expected = map[string]map[string]*RWStat{
		"launchd": {
			"/1": {
				Reads:       0,
				Writes:      1024,
				LastUpdated: getLastUpdated(out["launchd"]["/1"]),
			},
			"<Smaller Files of launchd>": {
				Reads:       0,
				Writes:      1024,
				LastUpdated: getLastUpdated(out["launchd"]["<Smaller Files of launchd>"]),
			},
		},
		"another": {
			"/1": {
				Reads:       0,
				Writes:      1024,
				LastUpdated: getLastUpdated(out["another"]["/1"]),
			},
			"<Smaller Files of another>": {
				Reads:       0,
				Writes:      512,
				LastUpdated: getLastUpdated(out["another"]["<Smaller Files of another>"]),
			},
		},
	}

	if !reflect.DeepEqual(out, expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", toJSON(expected), toJSON(out))
	}
}

func TestMergeDedupOn(t *testing.T) {
	mergeSmallFilesThreshold = 1024
	mergeMinStaleTime = 1
	dedupLaunchd = true

	processStat := NewProcessStat()

	processStat.Add(&IosnoopOutput{
		ExecName:      "launchd",
		OperationSize: 1024,
		IsRead:        false,
		Path:          "/1",
	})
	processStat.Add(&IosnoopOutput{
		ExecName:      "launchd",
		OperationSize: 512,
		IsRead:        false,
		Path:          "/2",
	})
	processStat.Add(&IosnoopOutput{
		ExecName:      "launchd",
		OperationSize: 512,
		IsRead:        false,
		Path:          "/3",
	})
	processStat.Add(&IosnoopOutput{
		ExecName:      "another",
		OperationSize: 1024,
		IsRead:        false,
		Path:          "/1",
	})
	processStat.Add(&IosnoopOutput{
		ExecName:      "another",
		OperationSize: 512,
		IsRead:        false,
		Path:          "/2",
	})

	out := processStat.GetAll()
	expected := map[string]map[string]*RWStat{
		"launchd": {
			"/3": {
				Reads:       0,
				Writes:      512,
				LastUpdated: getLastUpdated(out["launchd"]["/3"]),
			},
		},
		"another": {
			"/1": {
				Reads:       0,
				Writes:      1024,
				LastUpdated: getLastUpdated(out["another"]["/1"]),
			},
			"/2": {
				Reads:       0,
				Writes:      512,
				LastUpdated: getLastUpdated(out["another"]["/2"]),
			},
		},
	}

	if !reflect.DeepEqual(out, expected) {
		t.Errorf("expected %s, got %s", toJSON(expected), toJSON(out))
	}

	time.Sleep(1100 * time.Millisecond)

	processStat.MergeSmallEntries()
	out = processStat.GetAll()
	expected = map[string]map[string]*RWStat{
		"launchd": {
			"<Smaller Files of launchd>": {
				Reads:       0,
				Writes:      512,
				LastUpdated: getLastUpdated(out["launchd"]["<Smaller Files of launchd>"]),
			},
		},
		"another": {
			"/1": {
				Reads:       0,
				Writes:      1024,
				LastUpdated: getLastUpdated(out["another"]["/1"]),
			},
			"<Smaller Files of another>": {
				Reads:       0,
				Writes:      512,
				LastUpdated: getLastUpdated(out["another"]["<Smaller Files of another>"]),
			},
		},
	}

	if !reflect.DeepEqual(out, expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", toJSON(expected), toJSON(out))
	}

	// Test that the merged entries are not merged again.
	time.Sleep(1100 * time.Millisecond)

	processStat.MergeSmallEntries()
	out = processStat.GetAll()
	expected = map[string]map[string]*RWStat{
		"launchd": {
			"<Smaller Files of launchd>": {
				Reads:       0,
				Writes:      512,
				LastUpdated: getLastUpdated(out["launchd"]["<Smaller Files of launchd>"]),
			},
		},
		"another": {
			"/1": {
				Reads:       0,
				Writes:      1024,
				LastUpdated: getLastUpdated(out["another"]["/1"]),
			},
			"<Smaller Files of another>": {
				Reads:       0,
				Writes:      512,
				LastUpdated: getLastUpdated(out["another"]["<Smaller Files of another>"]),
			},
		},
	}

	if !reflect.DeepEqual(out, expected) {
		t.Errorf("expected:\n%s\ngot:\n%s", toJSON(expected), toJSON(out))
	}
}

func toJSON(in any) string {
	b, _ := json.MarshalIndent(in, "", "  ")
	return string(b)
}

func getLastUpdated(in *RWStat) time.Time {
	if in == nil {
		return time.Time{}
	}
	return in.LastUpdated
}
