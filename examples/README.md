# Diff Parser Examples

This directory contains examples showing how to use Zizou's diff parser.

## Running the Examples

```bash
# Line-by-line explanation (START HERE!)
go run examples/line_by_line_demo.go

# Simple usage example
go run examples/simple_usage.go

# Detailed parsing demo
go run examples/parse_demo.go
```

## Examples Explained

### 1. `line_by_line_demo.go` - **Start Here!**
Interactive explanation showing what each line in a diff means. Perfect for learning the diff format.

**Output:** Annotated diff with explanations for every line.

### 2. `simple_usage.go`
Shows common use cases: stats, finding added/removed lines, file analysis.

**Output:** Summary statistics and change listings.

### 3. `parse_demo.go`
Comprehensive demo showing all parser capabilities and edge cases.

**Output:** Detailed breakdown of hunks, lines, and context.

## Quick Start

```go
import "github.com/BhardwajShrey/zizou/internal/diff"

// Parse a diff
parser := diff.NewParser()
d, err := parser.Parse(diffContent)

// Get stats
stats := d.Stats()
fmt.Printf("Files: %d, +%d -%d\n",
    stats.Files, stats.LinesAdded, stats.LinesRemoved)

// Get added lines
for _, line := range d.GetAddedLines() {
    fmt.Printf("%s:%d: %s\n", line.File, line.Line, line.Content)
}
```

## Understanding the Output

When you run `simple_usage.go`, you'll see:

1. **Statistics**: Summary of changes
2. **Modified Files**: List of all changed files
3. **New Code**: All added lines with file:line references
4. **Removed Code**: All deleted lines
5. **Detailed Analysis**: Section-by-section breakdown

## What the Parser Extracts

From your git diff format, the parser extracts:

- **File paths** (old and new)
- **Line numbers** (accurate position in new file)
- **Change type** (added/removed/context)
- **Content** (the actual code)
- **Hunks** (grouped sections of changes)

## Available Helper Methods

```go
// On Diff
d.Stats()               // DiffStats
d.GetAddedLines()       // []LineInfo
d.GetRemovedLines()     // []LineInfo
d.GetModifiedFiles()    // []string

// On Hunk
hunk.HasChanges()                        // bool
hunk.GetContextAroundLine(index, size)   // []Line
```
