package codec

import (
	"testing"
)

func TestProtoCodec(t *testing.T) {
	codec := &Codec{}

	protoContent := `
syntax = "proto3";
package testservice;

// User message
message User {
    string name = 1;
    int32 id = 2;
    string email = 3;
}

enum Status {
    UNKNOWN = 0;
    ACTIVE = 1;
    INACTIVE = 2;
}
`
	var result map[string]any
	err := codec.Unmarshal([]byte(protoContent), &result)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	pkg, ok := result["testservice"].(map[string]any)
	if !ok {
		t.Fatalf("expected testservice package, got %v", result)
	}

	messages, ok := pkg["message"].(map[string]any)
	if !ok {
		t.Fatalf("expected messages, got %v", pkg)
	}
	if _, ok := messages["User"]; !ok {
		t.Errorf("expected User message in proto AST")
	}

	enums, ok := pkg["enum"].(map[string]any)
	if !ok {
		t.Fatalf("expected enums, got %v", pkg)
	}
	if _, ok := enums["Status"]; !ok {
		t.Errorf("expected Status enum in proto AST")
	}
}
