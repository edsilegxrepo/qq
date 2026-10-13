package util

import (
	"testing"
)

func TestOrderRegistry(t *testing.T) {
	ptr := uintptr(0x123456)
	keys := []string{"foo", "bar", "baz"}

	SetKeyOrder(ptr, keys)
	got, ok := GetKeyOrder(ptr)
	if !ok {
		t.Fatalf("expected key order to be found")
	}
	if len(got) != 3 || got[0] != "foo" || got[1] != "bar" || got[2] != "baz" {
		t.Errorf("unexpected key order: %v", got)
	}

	ClearKeyOrder()
	_, ok = GetKeyOrder(ptr)
	if ok {
		t.Errorf("expected key order to be cleared")
	}
}

func TestParseValue(t *testing.T) {
	DisableAutoConvert = false
	if v := ParseValue("123"); v != 123 {
		t.Errorf("expected 123, got %v", v)
	}
	if v := ParseValue("true"); v != true {
		t.Errorf("expected true, got %v", v)
	}
	if v := ParseValue("12.34"); v != 12.34 {
		t.Errorf("expected 12.34, got %v", v)
	}
	if v := ParseValue("2026-10-12"); v == "2026-10-12" {
		t.Errorf("expected parsed date, got %v", v)
	}

	DisableAutoConvert = true
	defer func() { DisableAutoConvert = false }()
	if v := ParseValue("123"); v != "123" {
		t.Errorf("expected string '123' when disabled, got %v", v)
	}
}

func TestNormalizeNumbers(t *testing.T) {
	input := map[string]any{
		"int":   int64(42),
		"float": float64(3.14),
		"slice": []any{uint(10), uint64(100)},
	}
	normalized := NormalizeNumbers(input)
	m, ok := normalized.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any")
	}
	if m["int"] != 42 {
		t.Errorf("expected normalized int 42, got %v", m["int"])
	}
}
