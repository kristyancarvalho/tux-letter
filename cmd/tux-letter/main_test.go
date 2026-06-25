package main

import (
	"reflect"
	"testing"
)

func TestSplitCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		command string
		rest    []string
	}{
		{"empty", nil, "", nil},
		{"once", []string{"once"}, "once", []string{}},
		{"flag before command", []string{"--config", "c.toml", "once"}, "once", []string{"--config", "c.toml"}},
		{"sources test", []string{"sources", "test"}, "sources test", []string{}},
		{"sources test with flag", []string{"sources", "test", "--config", "c.toml"}, "sources test", []string{"--config", "c.toml"}},
		{"only flags", []string{"--version"}, "", []string{"--version"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, rest := splitCommand(tt.args)
			if cmd != tt.command {
				t.Errorf("command = %q, want %q", cmd, tt.command)
			}
			if len(rest) != 0 || len(tt.rest) != 0 {
				if !reflect.DeepEqual(rest, tt.rest) {
					t.Errorf("rest = %v, want %v", rest, tt.rest)
				}
			}
		})
	}
}

func TestRunVersion(t *testing.T) {
	if code := run([]string{"--version"}); code != 0 {
		t.Errorf("run --version = %d, want 0", code)
	}
}

func TestRunUnknownCommand(t *testing.T) {
	if code := run([]string{"bogus"}); code != 2 {
		t.Errorf("run bogus = %d, want 2", code)
	}
}
