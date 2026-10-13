package cbor

import (
	"reflect"
	"testing"
)

func TestCborCodec(t *testing.T) {
	codec := &Codec{}

	original := map[string]any{
		"name": "alice",
		"id":   int64(42),
		"tags": []any{"go", "cbor"},
	}

	data, err := codec.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded any
	if err := codec.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	m, ok := decoded.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", decoded)
	}
	if m["name"] != "alice" {
		t.Errorf("expected alice, got %v", m["name"])
	}

	// Test Unmarshal error with corrupted data
	var bad any
	if err := codec.Unmarshal([]byte{0xff, 0xff}, &bad); err == nil {
		t.Errorf("expected error unmarshaling invalid CBOR, got nil")
	}

	// Test typed destination
	var typed map[string]any
	if err := codec.Unmarshal(data, &typed); err != nil {
		t.Fatalf("typed Unmarshal failed: %v", err)
	}
	if !reflect.DeepEqual(typed["name"], "alice") {
		t.Errorf("expected alice, got %v", typed["name"])
	}
}
