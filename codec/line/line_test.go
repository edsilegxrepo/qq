package line

import (
	"testing"
)

func TestLineCodec(t *testing.T) {
	codec := &Codec{}

	input := "line1\n123\ntrue\nline4\n"
	var result []any
	err := codec.Unmarshal([]byte(input), &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(result) != 4 {
		t.Fatalf("expected 4 lines, got %d", len(result))
	}
	if result[0] != "line1" {
		t.Errorf("expected line1, got %v", result[0])
	}
	if result[1] != float64(123) {
		t.Errorf("expected 123, got %v", result[1])
	}
	if result[2] != true {
		t.Errorf("expected true, got %v", result[2])
	}

	// Error when v is not a pointer
	var notPtr []any
	if err := codec.Unmarshal([]byte(input), notPtr); err == nil {
		t.Errorf("expected error when destination is not a pointer, got nil")
	}
}
