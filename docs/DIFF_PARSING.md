# Git Diff Parsing Guide

## Understanding Git Diff Format

### Structure Overview

A git diff consists of several sections:

```diff
diff --git a/app/main.go b/app/main.go          ← File Header
index 30dcdee..2cb2e68 100644                   ← Git Index (metadata)
--- a/app/main.go                               ← Old File Marker
+++ b/app/main.go                               ← New File Marker
@@ -23,7 +23,8 @@ func main() {                ← Hunk Header

     image := os.Args[2]                         ← Context Line (space prefix)

-    EnterNewJail(os.Args[3], image)            ← Removed Line (- prefix)
+    jailPath := EnterNewJail(os.Args[3], image) ← Added Line (+ prefix)
+    defer os.Remove(jailPath)                   ← Added Line (+ prefix)
```

### Section Breakdown

#### 1. File Header
```diff
diff --git a/app/main.go b/app/main.go
```
- Identifies which file was modified
- `a/` prefix = old version
- `b/` prefix = new version

#### 2. Index Line
```diff
index 30dcdee..2cb2e68 100644
```
- Git object hashes (old..new)
- File mode (100644 = regular file)

#### 3. File Markers
```diff
--- a/app/main.go
+++ b/app/main.go
```
- `---` marks the old file path
- `+++` marks the new file path

#### 4. Hunk Header
```diff
@@ -23,7 +23,8 @@ func main() {
```
Breaking down: `@@ -23,7 +23,8 @@`
- `-23,7` = In old file: starting at line 23, showing 7 lines
- `+23,8` = In new file: starting at line 23, showing 8 lines
- Optional context after `@@` (e.g., function name)

#### 5. Change Lines

Three types of lines:
- **Context** (space prefix): `  image := os.Args[2]`
- **Removed** (- prefix): `-    EnterNewJail(os.Args[3], image)`
- **Added** (+ prefix): `+    jailPath := EnterNewJail(os.Args[3], image)`

## Parser Architecture

### Data Structures

```go
// Top-level container for the entire diff
type Diff struct {
    Files []FileDiff
}

// Represents changes to a single file
type FileDiff struct {
    OldPath string      // Path in old version
    NewPath string      // Path in new version
    Hunks   []Hunk      // Sections of changes
}

// A continuous section of changes
type Hunk struct {
    OldStart int        // Starting line in old file
    OldCount int        // Number of lines in old file
    NewStart int        // Starting line in new file
    NewCount int        // Number of lines in new file
    Lines    []Line     // Actual line changes
}

// A single line in the diff
type Line struct {
    Type    LineType    // added, removed, context
    Content string      // The actual line content
    Number  int         // Line number in new file
}
```

### Parsing Algorithm

The parser uses a state machine approach:

```go
1. Scan line by line
2. Match regex patterns to identify section types
3. Build nested structure:
   - File headers create new FileDiff
   - Hunk headers create new Hunk
   - +/- lines create Line entries
4. Append completed sections to parent structures
```

### Regular Expressions

```go
fileHeaderRegex = `^diff --git a/(.*) b/(.*)$`
oldFileRegex    = `^--- a/(.*)$`
newFileRegex    = `^\+\+\+ b/(.*)$`
hunkHeaderRegex = `^@@ -(\d+),?(\d*) \+(\d+),?(\d*) @@`
```

## Usage Examples

### Basic Parsing

```go
import "github.com/BhardwajShrey/zizou/internal/diff"

parser := diff.NewParser()
parsedDiff, err := parser.Parse(diffContent)
if err != nil {
    log.Fatal(err)
}

// Access parsed data
for _, file := range parsedDiff.Files {
    fmt.Printf("File: %s\n", file.NewPath)

    for _, hunk := range file.Hunks {
        for _, line := range hunk.Lines {
            switch line.Type {
            case diff.LineAdded:
                fmt.Printf("+ %s\n", line.Content)
            case diff.LineRemoved:
                fmt.Printf("- %s\n", line.Content)
            case diff.LineContext:
                fmt.Printf("  %s\n", line.Content)
            }
        }
    }
}
```

### Using Helper Methods

```go
// Get statistics
stats := parsedDiff.Stats()
fmt.Printf("Files: %d, +%d -%d\n",
    stats.Files, stats.LinesAdded, stats.LinesRemoved)

// Get only added lines
for _, lineInfo := range parsedDiff.GetAddedLines() {
    fmt.Printf("%s:%d: %s\n",
        lineInfo.File, lineInfo.Line, lineInfo.Content)
}

// Get list of modified files
files := parsedDiff.GetModifiedFiles()

// Get context around a change
for _, file := range parsedDiff.Files {
    for _, hunk := range file.Hunks {
        // Get 3 lines of context around line index 5
        context := hunk.GetContextAroundLine(5, 3)
    }
}
```

### Common Patterns

#### Find Security-Sensitive Changes

```go
for _, file := range parsedDiff.Files {
    for _, hunk := range file.Hunks {
        for _, line := range hunk.Lines {
            if line.Type == diff.LineAdded {
                // Check for security issues
                if strings.Contains(line.Content, "exec.Command") {
                    fmt.Printf("Alert: Command execution at %s:%d\n",
                        file.NewPath, line.Number)
                }
            }
        }
    }
}
```

#### Generate Summary

```go
stats := parsedDiff.Stats()
fmt.Printf(`
Diff Summary:
- Files changed: %d
- Lines added: %d
- Lines removed: %d
- Net change: %+d
`, stats.Files, stats.LinesAdded, stats.LinesRemoved,
   stats.LinesAdded - stats.LinesRemoved)
```

## Advanced Topics

### Handling Edge Cases

1. **New files**: `--- /dev/null`
2. **Deleted files**: `+++ /dev/null`
3. **Binary files**: `Binary files differ`
4. **Renamed files**: Detected via `rename from/to`

### Performance Considerations

- Uses `bufio.Scanner` for efficient line-by-line reading
- Regex patterns compiled once at package level
- Memory-efficient for large diffs

### Testing Your Parser

```bash
# Run the demo
cd examples
go run parse_demo.go

# Test with real diffs
git diff | go run ../main.go --file -
```

## Integration with Zizou

The diff parser is the first step in Zizou's review pipeline:

```
Git Diff → Parser → Structured Data → Claude API → Review Comments
```

See `internal/review/prompt.go` for how the parsed diff is formatted for Claude.
