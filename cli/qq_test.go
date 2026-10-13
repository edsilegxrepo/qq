package cli

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/JFryy/qq/codec"
	"github.com/itchyny/gojq"
)

func TestInferFileType(t *testing.T) {
	tests := []struct {
		filename string
		expected codec.EncodingType
	}{
		{"test.json", codec.JSON},
		{"test.yaml", codec.YAML},
		{"test.yml", codec.YAML}, // yml maps to YAML
		{"test.toml", codec.TOML},
		{"test.xml", codec.XML},
		{"test.csv", codec.CSV},
		{"test.hcl", codec.HCL},
		{"test.tf", codec.HCL}, // tf maps to HCL
		{"test.ini", codec.INI},
		{"test.unknown", codec.JSON}, // defaults to JSON
		{"/path/to/file.json", codec.JSON},
		{"FILE.JSON", codec.JSON}, // case insensitive
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := inferFileType(tt.filename)
			if result != tt.expected {
				t.Errorf("inferFileType(%q) = %v, expected %v", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestIsFile(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "test*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpfile.Name()) }()
	_ = tmpfile.Close()

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"existing file", tmpfile.Name(), true},
		{"non-existing file", "/nonexistent/file.txt", false},
		{"directory", "/tmp", false},
		{"empty path", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isFile(tt.path)
			if result != tt.expected {
				t.Errorf("isFile(%q) = %v, expected %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestIsTerminal(t *testing.T) {
	// Test with os.Stdout - this might be a terminal or not depending on test environment
	result := isTerminal(os.Stdout)
	t.Logf("isTerminal(os.Stdout) = %v", result)

	// Create a pipe to test non-terminal
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()
	defer func() { _ = w.Close() }()

	if isTerminal(w) {
		t.Error("Pipe should not be detected as terminal")
	}
}

func TestCreateRootCmd(t *testing.T) {
	cmd := CreateRootCmd()

	if cmd == nil {
		t.Fatal("CreateRootCmd returned nil")
	}

	if cmd.Use == "" {
		t.Error("Command Use should not be empty")
	}

	// Test flags exist
	flags := cmd.Flags()
	if !flags.HasFlags() {
		t.Error("Command should have flags")
	}

	requiredFlags := []string{"input", "output", "raw-output", "monochrome-output", "interactive"}
	for _, flag := range requiredFlags {
		if flags.Lookup(flag) == nil {
			t.Errorf("Missing required flag: %s", flag)
		}
	}
}

func TestRootCmdFlagDefaults(t *testing.T) {
	cmd := CreateRootCmd()

	tests := []struct {
		flagName string
		expected string
	}{
		{"input", "json"},
		{"output", "json"},
	}

	for _, tt := range tests {
		t.Run(tt.flagName, func(t *testing.T) {
			value, err := cmd.Flags().GetString(tt.flagName)
			if err != nil {
				t.Fatalf("Failed to get flag %s: %v", tt.flagName, err)
			}
			if value != tt.expected {
				t.Errorf("Default value for --%s = %q, expected %q", tt.flagName, value, tt.expected)
			}
		})
	}

	// Test boolean flag defaults
	boolTests := []struct {
		flagName string
		expected bool
	}{
		{"raw-output", false},
		{"monochrome-output", false},
		{"interactive", false},
		{"preserve-key-order", false},
	}

	for _, tt := range boolTests {
		t.Run(tt.flagName, func(t *testing.T) {
			value, err := cmd.Flags().GetBool(tt.flagName)
			if err != nil {
				t.Fatalf("Failed to get flag %s: %v", tt.flagName, err)
			}
			if value != tt.expected {
				t.Errorf("Default value for --%s = %v, expected %v", tt.flagName, value, tt.expected)
			}
		})
	}
}

func TestExecuteQueryBasic(t *testing.T) {
	// This is a basic integration test that doesn't require external dependencies
	// More comprehensive tests should be in E2E tests

	tests := []struct {
		name     string
		input    string
		query    string
		fileType codec.EncodingType
		wantErr  bool
	}{
		{
			name:     "simple identity query",
			input:    `{"key": "value"}`,
			query:    ".",
			fileType: codec.JSON,
			wantErr:  false,
		},
		{
			name:     "key extraction",
			input:    `{"key": "value"}`,
			query:    ".key",
			fileType: codec.JSON,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the input
			var data any
			err := codec.Unmarshal([]byte(tt.input), tt.fileType, &data)
			if err != nil {
				t.Fatalf("Failed to unmarshal input: %v", err)
			}

			// Note: Full query execution testing should be in E2E tests
			// This just verifies the data can be unmarshaled
			if data == nil && !tt.wantErr {
				t.Error("Expected data, got nil")
			}
		})
	}
}

func TestVersionFlag(t *testing.T) {
	cmd := CreateRootCmd()
	cmd.SetArgs([]string{"--version"})

	// Capture output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run command in goroutine and capture exit
	done := make(chan bool)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// Command might call os.Exit, which we can't easily test
				done <- true
			}
		}()
		_ = cmd.Execute()
		done <- true
	}()

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	<-done

	// Version output should contain "qq version"
	if !strings.Contains(output, "qq version") && output != "" {
		t.Logf("Version output: %q", output)
	}
}

func TestHelpFlag(t *testing.T) {
	cmd := CreateRootCmd()

	// Get help text
	helpText := cmd.UsageString()

	if helpText == "" {
		t.Error("Help text should not be empty")
	}

	// Should mention key features
	if !strings.Contains(strings.ToLower(helpText), "json") {
		t.Error("Help text should mention json")
	}

	// Should have the basic flags
	if !strings.Contains(helpText, "input") || !strings.Contains(helpText, "output") {
		t.Error("Help text should mention input/output flags")
	}
}

func TestExecuteQuery(t *testing.T) {
	q, err := gojq.Parse(".")
	if err != nil {
		t.Fatalf("gojq.Parse failed: %v", err)
	}

	// Normal JSON
	code := executeQuery(q, map[string]any{"test": 123}, codec.JSON, false, true, false)
	if code != 0 {
		t.Errorf("expected 0, got %d", code)
	}

	// Raw output with string
	qStr, _ := gojq.Parse(".name")
	code = executeQuery(qStr, map[string]any{"name": "hello"}, codec.JSON, true, true, false)
	if code != 0 {
		t.Errorf("expected 0, got %d", code)
	}

	// Binary format
	code = executeQuery(q, map[string]any{"bin": 1}, codec.CBOR, false, true, false)
	if code != 0 {
		t.Errorf("expected 0, got %d", code)
	}

	// Exit status: true -> 0
	qTrue, _ := gojq.Parse(".ok")
	code = executeQuery(qTrue, map[string]any{"ok": true}, codec.JSON, false, true, true)
	if code != 0 {
		t.Errorf("expected 0, got %d", code)
	}

	// Exit status: false -> 1
	code = executeQuery(qTrue, map[string]any{"ok": false}, codec.JSON, false, true, true)
	if code != 1 {
		t.Errorf("expected 1, got %d", code)
	}

	// Exit status: nil -> 1
	code = executeQuery(qTrue, map[string]any{"ok": nil}, codec.JSON, false, true, true)
	if code != 1 {
		t.Errorf("expected 1, got %d", code)
	}

	// Exit status: empty -> 4
	qEmpty, _ := gojq.Parse("empty")
	code = executeQuery(qEmpty, map[string]any{"ok": 1}, codec.JSON, false, true, true)
	if code != 4 {
		t.Errorf("expected 4, got %d", code)
	}
}

func TestSlurpInputs(t *testing.T) {
	// JSON
	jsonIn := []byte(`{"a": 1} {"b": 2}`)
	res, err := slurpInputs(jsonIn, codec.JSON)
	if err != nil {
		t.Fatalf("slurpInputs JSON failed: %v", err)
	}
	arr, ok := res.([]any)
	if !ok || len(arr) != 2 {
		t.Errorf("expected array of 2 elements, got %v", res)
	}

	// JSONL
	jsonlIn := []byte("{\"a\": 1}\n{\"b\": 2}\n")
	res, err = slurpInputs(jsonlIn, codec.JSONL)
	if err != nil {
		t.Fatalf("slurpInputs JSONL failed: %v", err)
	}
	arr, ok = res.([]any)
	if !ok || len(arr) != 2 {
		t.Errorf("expected array of 2 elements, got %v", res)
	}

	// YAML multi-doc
	yamlIn := []byte("a: 1\n---\nb: 2\n")
	res, err = slurpInputs(yamlIn, codec.YAML)
	if err != nil {
		t.Fatalf("slurpInputs YAML failed: %v", err)
	}
	arr, ok = res.([]any)
	if !ok || len(arr) != 2 {
		t.Errorf("expected array of 2 elements, got %v", res)
	}

	// YAML single doc
	yamlSingle := []byte("a: 1\n")
	res, err = slurpInputs(yamlSingle, codec.YAML)
	if err != nil {
		t.Fatalf("slurpInputs YAML single failed: %v", err)
	}
	arr, ok = res.([]any)
	if !ok || len(arr) != 1 {
		t.Errorf("expected array of 1 element, got %v", res)
	}

	// LINE
	lineIn := []byte("foo\nbar\n")
	res, err = slurpInputs(lineIn, codec.LINE)
	if err != nil {
		t.Fatalf("slurpInputs LINE failed: %v", err)
	}
	arr, ok = res.([]any)
	if !ok || len(arr) != 2 {
		t.Errorf("expected array of 2 elements, got %v", res)
	}

	// Default (TOML)
	tomlIn := []byte("title = \"qq\"\n")
	res, err = slurpInputs(tomlIn, codec.TOML)
	if err != nil {
		t.Fatalf("slurpInputs TOML failed: %v", err)
	}
	arr, ok = res.([]any)
	if !ok || len(arr) != 1 {
		t.Errorf("expected array of 1 element, got %v", res)
	}
}

func TestResolveVersion(t *testing.T) {
	v := resolveVersion()
	if v == "" {
		t.Errorf("expected non-empty version")
	}
}

type exitPanic int

func withMockExit(fn func()) (code int, exited bool) {
	origExit := osExit
	defer func() {
		osExit = origExit
	}()
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(exitPanic); ok {
				code = int(e)
				exited = true
			} else {
				panic(r)
			}
		}
	}()
	osExit = func(c int) {
		panic(exitPanic(c))
	}
	fn()
	return 0, false
}

func TestHandleCommand(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_handle_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	_, _ = tmpFile.WriteString(`{"name":"test","age":25}`)
	_ = tmpFile.Close()

	cmd := CreateRootCmd()

	// 1. Help flag
	code, exited := withMockExit(func() {
		handleCommand(cmd, nil, "json", "json", false, true, false, true, false, false, false, false, false)
	})
	if !exited || code != 0 {
		t.Errorf("expected exit 0 for help, got %d (exited: %v)", code, exited)
	}

	// 2. Stream + Interactive mutual exclusion
	code, exited = withMockExit(func() {
		handleCommand(cmd, nil, "json", "json", false, false, true, true, true, false, false, false, false)
	})
	if !exited || code != 1 {
		t.Errorf("expected exit 1 for stream+interactive, got %d", code)
	}

	// 3. Slurp + Stream mutual exclusion
	code, exited = withMockExit(func() {
		handleCommand(cmd, nil, "json", "json", false, false, false, true, true, true, false, false, false)
	})
	if !exited || code != 1 {
		t.Errorf("expected exit 1 for slurp+stream, got %d", code)
	}

	// 4. Slurp + Interactive mutual exclusion
	code, exited = withMockExit(func() {
		handleCommand(cmd, nil, "json", "json", false, false, true, true, false, true, false, false, false)
	})
	if !exited || code != 1 {
		t.Errorf("expected exit 1 for slurp+interactive, got %d", code)
	}

	// 5. Valid file input with query
	code, exited = withMockExit(func() {
		handleCommand(cmd, []string{".name", tmpFile.Name()}, "json", "json", false, false, false, true, false, false, false, false, false)
	})
	if !exited || code != 0 {
		t.Errorf("expected exit 0 for valid file input, got %d", code)
	}

	// 6. Valid file input with single arg (infer query ".")
	code, exited = withMockExit(func() {
		handleCommand(cmd, []string{tmpFile.Name()}, "json", "json", false, false, false, true, false, false, false, false, false)
	})
	if !exited || code != 0 {
		t.Errorf("expected exit 0 for file argument without query, got %d", code)
	}

	// 7. Streaming mode with file
	code, exited = withMockExit(func() {
		handleCommand(cmd, []string{".", tmpFile.Name()}, "json", "json", false, false, false, true, true, false, false, false, false)
	})
	if !exited || code != 0 {
		t.Errorf("expected exit 0 for streaming mode with file, got %d", code)
	}

	// 8. Slurp mode with file
	code, exited = withMockExit(func() {
		handleCommand(cmd, []string{".", tmpFile.Name()}, "json", "json", false, false, false, true, false, true, false, false, false)
	})
	if !exited || code != 0 {
		t.Errorf("expected exit 0 for slurp mode with file, got %d", code)
	}

	// 9. Invalid query syntax
	code, exited = withMockExit(func() {
		handleCommand(cmd, []string{".[invalid", tmpFile.Name()}, "json", "json", false, false, false, true, false, false, false, false, false)
	})
	if !exited || code != 1 {
		t.Errorf("expected exit 1 for invalid query, got %d", code)
	}

	// 10. Nonexistent file
	code, exited = withMockExit(func() {
		handleCommand(cmd, []string{".", "nonexistent_file_12345.json"}, "json", "json", false, false, false, true, false, false, false, false, false)
	})
	if !exited || code != 1 {
		t.Errorf("expected exit 1 for nonexistent file, got %d", code)
	}

	// 11. Stdin pipe input with expression
	oldStdin := os.Stdin
	rPipe, wPipe, _ := os.Pipe()
	_, _ = wPipe.WriteString(`{"piped": 100}`)
	_ = wPipe.Close()
	os.Stdin = rPipe
	code, exited = withMockExit(func() {
		handleCommand(cmd, []string{".piped"}, "json", "json", false, false, false, true, false, false, false, false, false)
	})
	_ = rPipe.Close()
	os.Stdin = oldStdin
	if !exited || code != 0 {
		t.Errorf("expected exit 0 for stdin pipe, got %d", code)
	}

	// 12. Stdin pipe input with 0 args
	rPipe, wPipe, _ = os.Pipe()
	_, _ = wPipe.WriteString(`{"empty_args": true}`)
	_ = wPipe.Close()
	os.Stdin = rPipe
	code, exited = withMockExit(func() {
		handleCommand(cmd, []string{}, "json", "json", false, false, false, true, false, false, false, false, false)
	})
	_ = rPipe.Close()
	os.Stdin = oldStdin
	if !exited || code != 0 {
		t.Errorf("expected exit 0 for stdin pipe with 0 args, got %d", code)
	}

	// 13. Explicit input/output flags and flags changed
	cmdFlags := CreateRootCmd()
	_ = cmdFlags.Flags().Set("input", "json")
	_ = cmdFlags.Flags().Set("output", "yaml")
	code, exited = withMockExit(func() {
		handleCommand(cmdFlags, []string{".", tmpFile.Name()}, "json", "yaml", false, false, false, true, false, false, false, true, true)
	})
	if !exited || code != 0 {
		t.Errorf("expected exit 0 for explicit flags, got %d", code)
	}

	// 14. Invalid input codec
	code, exited = withMockExit(func() {
		handleCommand(cmdFlags, []string{".", tmpFile.Name()}, "invalid_format_xyz", "json", false, false, false, true, false, false, false, false, false)
	})
	if !exited || code != 1 {
		t.Errorf("expected exit 1 for invalid input format, got %d", code)
	}

	// 15. Invalid output codec
	cmdFlags2 := CreateRootCmd()
	_ = cmdFlags2.Flags().Set("input", "json")
	code, exited = withMockExit(func() {
		handleCommand(cmdFlags2, []string{".", tmpFile.Name()}, "json", "invalid_format_xyz", false, false, false, true, false, false, false, false, false)
	})
	if !exited || code != 1 {
		t.Errorf("expected exit 1 for invalid output format, got %d", code)
	}
}
