package ini

import (
	"strings"
	"testing"
)

func TestIniUnmarshal(t *testing.T) {
	input := `
# Comment
[database]
server = 192.0.2.62
port = 5432
enabled = true

[owner]
name = John Doe
organization = Acme Widgets
`
	var result map[string]any
	codec := &Codec{}
	err := codec.Unmarshal([]byte(input), &result)
	if err != nil {
		t.Fatalf("unexpected error unmarshaling INI: %v", err)
	}

	db, ok := result["database"].(map[string]any)
	if !ok {
		t.Fatalf("expected database section, got: %v", result["database"])
	}
	if db["server"] != "192.0.2.62" {
		t.Errorf("expected server 192.0.2.62, got: %v", db["server"])
	}
	if db["port"] != 5432 {
		t.Errorf("expected port 5432, got: %v", db["port"])
	}
	if db["enabled"] != true {
		t.Errorf("expected enabled true, got: %v", db["enabled"])
	}
}

func TestIniMarshal(t *testing.T) {
	data := map[string]any{
		"database": map[string]any{
			"port":    5432,
			"enabled": true,
			"server":  "127.0.0.1",
		},
		"scalar": "standalone_value",
	}
	codec := &Codec{}
	out, err := codec.Marshal(data)
	if err != nil {
		t.Fatalf("unexpected error marshaling INI: %v", err)
	}

	s := string(out)
	if !strings.Contains(s, "[database]") {
		t.Errorf("expected [database] section in output, got: %s", s)
	}
	if !strings.Contains(s, "5432") {
		t.Errorf("expected 5432 in output, got: %s", s)
	}
}

func TestIniErrors(t *testing.T) {
	codec := &Codec{}
	if err := codec.Unmarshal([]byte("[invalid"), &map[string]any{}); err == nil {
		t.Errorf("expected error unmarshaling invalid ini, got nil")
	}
	if _, err := codec.Marshal("not a map"); err == nil {
		t.Errorf("expected error marshaling non-map, got nil")
	}
}
