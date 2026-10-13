package tui

import (
	"strings"
	"testing"
)

func TestGenerateJqOptions(t *testing.T) {
	input := `{"user": {"name": "alice", "id": 1}, "tags": ["admin", "dev"]}`
	options := generateJqOptions(input)

	if len(options) == 0 {
		t.Fatalf("expected extracted paths, got 0")
	}

	foundUser := false
	for _, opt := range options {
		if strings.Contains(opt, ".user") {
			foundUser = true
			break
		}
	}
	if !foundUser {
		t.Errorf("expected .user in jq options, got %v", options)
	}

	// Invalid JSON should return fallback "."
	invalidOpts := generateJqOptions("not json")
	if len(invalidOpts) != 1 || invalidOpts[0] != "." {
		t.Errorf("expected fallback [.], got %v", invalidOpts)
	}
}

func TestJsonStrToInterface(t *testing.T) {
	valid := `{"foo":"bar"}`
	res, err := jsonStrToInterface(valid)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := res.(map[string]any)
	if !ok || m["foo"] != "bar" {
		t.Errorf("expected map with foo=bar, got %v", res)
	}

	invalid := `{"foo":`
	_, err = jsonStrToInterface(invalid)
	if err == nil {
		t.Errorf("expected error on invalid JSON, got nil")
	}
}

func TestModelQueryLogic(t *testing.T) {
	m := newModel(`{"numbers":[1,2,3]}`)

	// Empty query is invalid
	m.textArea.SetValue("")
	if m.isValidQuery() {
		t.Errorf("empty query should be invalid")
	}

	// Valid query
	m.textArea.SetValue(".numbers[]")
	if !m.isValidQuery() {
		t.Errorf("valid query .numbers[] should be valid")
	}

	// Run filter
	m.runJqFilter()
	if !strings.Contains(m.jqOutput, "1") {
		t.Errorf("expected 1 in output, got: %s", m.jqOutput)
	}

	// Invalid query should report error
	m.textArea.SetValue(".[invalid")
	m.runJqFilter()
	if !strings.Contains(m.jqOutput, "Invalid jq query") {
		t.Errorf("expected Invalid jq query in output, got: %s", m.jqOutput)
	}
}
