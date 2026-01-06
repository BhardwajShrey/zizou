# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Zizou is an AI-powered code review CLI tool that reads git diffs and provides structured feedback using Claude API. It parses git diffs, sends them to Claude for analysis, and returns categorized review comments with severity levels.

## Communication Style

When working in this repository, keep all output **concise yet informative**:

- **Be direct**: Get to the point quickly. State what you're doing, why, and the outcome.
- **Be specific**: Reference exact files and line numbers (e.g., "parser.go:215-235").
- **Be actionable**: Provide concrete next steps or specific code examples when suggesting changes.
- **Avoid verbosity**: Skip unnecessary explanations, redundant confirmations, or over-documentation.
- **Match the project philosophy**: Zizou itself emphasizes "concise, impact-driven feedback" and "optimize for minimal tokens" (see reviewer_client.go:55, 137). Apply the same principle to your communication. Even the git commit messages should be as concise as possible
- **Be curious**: Think hard on each problem and don't hesitate to point out any flaws in my decisions. Always present me with follow-up and food-for-thought questions if there are any.

Example:
- ❌ "I'm going to analyze the parser code to understand how it handles line numbers, which is important because..."
- ✅ "Checking parser line number handling in parser.go:215-235"

## Build & Test Commands

### Building
```bash
# Build the binary
go build

# Build with specific output name
go build -o zizou

# Install globally
go install github.com/BhardwajShrey/zizou@latest
```

### Testing
```bash
# Run all tests (note: examples directory has build conflicts)
go test ./internal/... -v

# Run tests with coverage
go test ./internal/... -cover

# Run specific package tests
go test ./internal/diff -v
go test ./internal/client -v

# Generate coverage report
go test ./internal/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Running
```bash
# Basic usage
git diff | ./zizou

# From file
./zizou --file test_diff.txt

# With output formats
./zizou --file test_diff.txt --output json
./zizou --file test_diff.txt --output markdown

# Set API key
export ANTHROPIC_API_KEY_ZIZOU="sk-ant-..."
./zizou --file test_diff.txt
```

## Architecture

### Core Flow
1. **Input** → Read diff from stdin or file (cmd/root.go)
2. **Parse** → Parse git diff into structured format (internal/diff/parser.go)
3. **Review** → Check cache, send to Claude API if needed (internal/review/reviewer.go)
4. **Output** → Format and display results (internal/output/formatter.go)

### Key Components

#### Diff Parser (internal/diff/)
- **parser.go**: Line-by-line parser for git diff format using regex patterns
- **types.go**: Core data structures (Diff, FileDiff, Hunk, Line)
- **helpers.go**: Utility functions for diff manipulation
- Handles complex edge cases: binary files, renames, mode changes, new/deleted files

The parser uses regex patterns to identify diff sections:
- File headers: `diff --git a/... b/...`
- Hunk headers: `@@ -23,7 +23,8 @@` (old line, old count, new line, new count)
- Line types: `+` (added), `-` (removed), ` ` (context)
- Metadata: renames, mode changes, similarity index

#### Review System (internal/review/)
- **reviewer.go**: Orchestrates review with caching (SHA256 hash of diff content)
- **comment.go**: Defines Comment and Result structures
- **prompt.go**: Legacy prompt builder (newer logic in reviewer_client.go)

Review flow:
1. Generate cache key from diff (SHA256 of JSON-serialized diff)
2. Check file cache (~/.zizou/cache/)
3. If miss, call ReviewerClient
4. Parse JSON response from Claude
5. Cache result for future use

#### Client (internal/client/)
- **reviewer_client.go**: Builds prompts and calls Claude API
- **enhanced_client.go**: Wraps Anthropic SDK with retry/rate limiting
- **rate_limiter.go**: Token bucket rate limiting
- **config.go**: API configuration (reads ANTHROPIC_API_KEY_ZIZOU)

The prompt structure (in buildReviewPrompt):
- Senior reviewer persona with concise, actionable feedback
- Focus areas: security, bugs, performance, maintainability
- Language-specific checks (Go: error handling, races; C++: memory safety; etc.)
- Diff statistics and file context
- Structured JSON output format with severity/category

#### Cache (internal/cache/)
- **cache.go**: Interface + FileCache and NoOpCache implementations
- SHA256-based keys stored as JSON files
- Default location: ~/.zizou/cache/

#### Output (internal/output/)
- **formatter.go**: Three formats (text, json, markdown)
- Groups by severity, adds icons/formatting

### Data Flow
```
stdin/file → Parser → Diff struct → Reviewer (check cache)
                                         ↓
                                   ReviewerClient
                                         ↓
                         buildReviewPrompt (detailed prompt)
                                         ↓
                         EnhancedClient (Anthropic SDK)
                                         ↓
                         parseReviewResponse (extract JSON)
                                         ↓
                                    Result struct
                                         ↓
                                Formatter → stdout
```

### Important Types

**Diff Structure**:
- `Diff`: Contains []FileDiff
- `FileDiff`: OldPath, NewPath, []Hunk, metadata flags (IsNew, IsDeleted, IsRenamed, IsBinary)
- `Hunk`: OldStart/Count, NewStart/Count, []Line
- `Line`: Type (added/removed/context), Content, Number

**Review Structure**:
- `Comment`: File, Line, Severity, Category, Message, Snippet
- `Result`: []Comment, Summary
- Severity: critical, high, medium, low, info
- Category: security, performance, bug_risk, maintainability, style, best_practice, other

## Development Notes

### Testing Edge Cases
The diff parser handles numerous edge cases (see internal/diff/*_test.go):
- Binary files
- File renames with similarity index
- Mode changes (e.g., chmod)
- New/deleted files
- Empty hunks
- Context-only diffs
- Multiple files in single diff

When modifying the parser, run the edge case tests:
```bash
go test ./internal/diff -v -run TestEdgeCases
go test ./internal/diff -v -run TestAdvancedEdgeCases
```

### Prompt Engineering
The review prompt (reviewer_client.go:51-139) is carefully structured:
- Senior reviewer persona for concise, impact-driven feedback
- Explicit severity/category definitions to guide Claude
- Language-specific focus areas (Go, C++, Python, JS)
- Structured JSON output format with examples
- Diff statistics for context
- Token optimization ("Be concise but complete")

When modifying prompts, test with various diff types to ensure JSON parsing remains reliable.

### API Key Management
Zizou expects the API key in environment variable `ANTHROPIC_API_KEY_ZIZOU` (not the standard `ANTHROPIC_API_KEY`). This allows users to have separate keys for different tools.

### Caching Strategy
- Cache key: SHA256 hash of JSON-serialized Diff struct
- Cache invalidation: Manual (delete ~/.zizou/cache/)
- Cache miss is not an error; triggers API call
- NoOpCache used when --no-cache flag is set

### Examples Directory
The examples/ directory contains multiple main() functions for demonstration purposes. These cause build failures in `go test ./...` but are intentionally kept as runnable examples:
- test_parser.go: Demonstrates parser without API calls
- claude_integration_example.go: Full integration example
- test_edge_cases.go: Edge case demonstrations

Run examples individually:
```bash
go run examples/test_parser.go
go run examples/claude_integration_example.go
```

## Git Diff Format Reference

Understanding the diff format is critical for parser work. See docs/DIFF_GUIDE.md for complete reference.

Key patterns:
- File header: `diff --git a/file.go b/file.go`
- Hunk header: `@@ -23,7 +23,8 @@` means:
  - Old file: starting at line 23, showing 7 lines
  - New file: starting at line 23, showing 8 lines
- Line prefixes: `+` added, `-` removed, ` ` context
- Special markers: `--- a/`, `+++ b/`, `--- /dev/null` (new file), `+++ /dev/null` (deleted)

## Common Pitfalls

1. **Parser Line Numbers**: Added lines increment lineNumber, removed lines don't. Context lines increment.
   - See parser.go:215-235 for line number tracking logic

2. **JSON Extraction**: Claude may wrap JSON in markdown code blocks. extractJSON() handles this.
   - See reviewer_client.go:227-255

3. **Cache Keys**: Changing Diff struct breaks cache. Consider this when modifying types.

4. **Empty Diffs**: Parser returns error for empty content. Handle at caller level.

5. **Binary Files**: Parser detects binary files but doesn't include content in review prompt.
   - See reviewer_client.go:148-152

## Documentation

- **docs/DIFF_GUIDE.md**: Complete git diff format explanation
- **docs/DIFF_PARSING.md**: Parser implementation guide
- **docs/EDGE_CASES.md**: Edge cases and how parser handles them
- **RUN_ZIZOU.md**: Complete usage guide with examples
