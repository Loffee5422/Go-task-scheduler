package main

import "testing"

func TestShouldSave(t *testing.T) {
	for _, command := range []string{"add", "update", "delete", "complete", "reopen", "sort"} {
		if !shouldSave(command) {
			t.Errorf("expected %q to save", command)
		}
	}
	for _, command := range []string{"list", "find", "filter", "help", "unknown"} {
		if shouldSave(command) {
			t.Errorf("expected %q not to save", command)
		}
	}
}
