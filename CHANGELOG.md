# Changelog

All notable changes to this project will be documented in this file.

---

## [0.4.0] - 2026-10-12

### Added
- **Dynamic Version Derivation & Dependency Management**:
  - Implemented dynamic version resolution via `runtime/debug.ReadBuildInfo()` when installed using `go install`, falling back gracefully to `"dev"` when `-ldflags` is not supplied.

### Fixed & Improved
- **TOML Array of Tables Direct Querying**:
  - Normalized generic TOML containers and arrays of tables to allow direct `gojq` queries without intermediate JSON serialization or integer precision loss.
- **HCL & Terraform Round-Trip Integrity**:
  - Corrected Terraform block attributes in HCL round-trip transcoding, accurately distinguishing between HCL blocks and object-list attributes.
- **Gron Stability & Deterministic Output**:
  - Resolved panic on malformed array assignments and invalid indices in Gron deserialization.
  - Added deterministic map key sorting during Gron traversal for reproducible outputs and diffs.
- **Large Integer Precision Preservation**:
  - Upgraded JSON, JSONL, and Gron generic deserialization to preserve large integer precision (values > $2^{53}$) without lossy IEEE-754 float64 conversion.
- **Malformed Environment Variable Rejection**:
  - Enforced strict syntax validation on `.env` file key-value pairs, rejecting malformed lines.
- **Deterministic Key Serialization**:
  - Added deterministic key sorting for INI and Parquet schema serialization.

### Dependencies & Maintenance
- Bumped core dependencies

---

## [0.3.5] - 2026-07-10

### Added
- **Key Order Preservation (`--preserve-key-order` / `-k`)**:
  - Custom JSON token parsing (`dec.Token()`) tracks and registers map key sequences in a global thread-safe registry (`sync.RWMutex`).
  - Custom YAML decoding registers AST mappings and normalizes pointers to guarantee key order preservation.
  - CSV/TSV codecs maintain column sequence matching the original headers on both import and export.
  - XML codecs traverse standard token streams to preserve order of sibling tags.
  - Added an unconditional global memory registry cleanup routine (`ClearKeyOrder()`) at query startup to prevent memory leaks during interactive REPL sessions and module integrations.
- **Auto-Coercion Disabling (`--no-auto-convert`)**:
  - Disables automatic string conversion during parsing (e.g. converting `"T"`, `"F"`, `"1"`, or `"0"` to booleans/numbers).
  - Supported across CSV, TSV, XML, INI, Gron, and Line formats.
  - Bypasses string trimming and parse routines, offering a significant parsing performance boost.
- **Interactive Mode Format Decoupling (`--interactive` / `-I`)**:
  - Always initializes the interactive TUI state buffer using JSON to prevent parsing errors when query prototyping with non-JSON target output formats.
  - Correctly serializes and formats the final query result to the target output format (e.g. YAML, XML, TOML) on graceful TUI exit.
- **Cross-Platform Build Detection**:
  - Added OS detection in `Makefile` to automatically build with the `.exe` extension on Windows hosts, ensuring native shell execution.
- **Architecture Documentation**:
  - Authored a comprehensive `ARCHITECTURE.md` file detailing application objectives, design choices, data flow pipelines, compatibility matrices, security audits, and testing structures.


### Testing & Code Quality
- **Integration Test Suite**:
  - Created `cli/integration_test.go` covering full integration pipelines.
  - Added E2E tests for JQ expression support (`.[]`, `length`, `first`, `.[0]`) on CSV inputs.
  - Added E2E tests for key preservation across all formats (JSON, JSONL, CSV, TSV, YAML, XML).
  - Added coercion tests for CSV, XML, INI, and Gron formats under `--no-auto-convert`.
  - Added edge case tests for Unicode characters and header-only empty files.
- **Clean Linter & Security Scan**:
  - Resolved all 24 warnings in the codebase.
  - Verified `golangci-lint run ./... --no-config` passes with **0 issues**.
  - Verified `gosec ./...` security scanner returns **0 issues**.
