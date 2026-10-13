package avro

import (
	"testing"
)

func TestAvroCodec(t *testing.T) {
	codec := &Codec{}

	records := []map[string]any{
		{
			"id":     "rec-1",
			"count":  float64(10),
			"active": true,
		},
		{
			"id":     "rec-2",
			"count":  float64(20),
			"active": false,
		},
	}

	data, err := codec.Marshal(records)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded []any
	if err := codec.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(decoded) != 2 {
		t.Fatalf("expected 2 records, got %d", len(decoded))
	}

	// Single record marshal
	single := map[string]any{"key": "value"}
	singleData, err := codec.Marshal(single)
	if err != nil {
		t.Fatalf("Single record Marshal failed: %v", err)
	}
	var singleDecoded []any
	if err := codec.Unmarshal(singleData, &singleDecoded); err != nil {
		t.Fatalf("Single record Unmarshal failed: %v", err)
	}
	if len(singleDecoded) != 1 {
		t.Errorf("expected 1 record, got %d", len(singleDecoded))
	}

	// Empty records marshal error
	if _, err := codec.Marshal([]map[string]any{}); err == nil {
		t.Errorf("expected error marshaling empty slice, got nil")
	}

	// Corrupted data unmarshal error
	var corrupted any
	if err := codec.Unmarshal([]byte("not valid avro"), &corrupted); err == nil {
		t.Errorf("expected error unmarshaling corrupted data, got nil")
	}
}
