package config

import (
	"testing"
)

func TestLoadBackends(t *testing.T) {
	// Test case 1: Empty args slice should result in an empty BackendConfig.
	args1 := []string{}
	expected1 := BackendConfig{}
	actual1 := LoadBackends(args1)

	if len(expected1) != len(actual1) {
		t.Errorf("Test case 1: Length mismatch. Expected %v, but got %v", expected1, actual1)
	}

	// Test case 2: Valid args should result in a populated BackendConfig.
	args2 := []string{"us=http://localhost:9001", "ru=http://localhost:9002"}
	expected2 := BackendConfig{
		"us": "http://localhost:9001",
		"ru": "http://localhost:9002",
	}
	actual2 := LoadBackends(args2)

	if len(expected2) != len(actual2) {
		t.Errorf("Test case 2: Length mismatch. Expected %v, but got %v", expected2, actual2)
	}
}
