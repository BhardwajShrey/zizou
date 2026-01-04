# Complete Git Diff Format Guide

## 🎯 Your Question: Explain the diff format, line by line

Here's the complete answer with your exact example.

---

## Your Diff - Annotated

```diff
Line 1:  diff --git a/app/main.go b/app/main.go
         ─────┬──── ────────┬──────── ────────┬────────
              │             │                  │
         File header    Old path          New path

         📁 Meaning: This diff compares the file "app/main.go"

Line 2:  index 30dcdee..2cb2e68 100644
         ───── ───┬─── ───┬─── ──┬───
               Old     New    File
               hash    hash   mode

         🔖 Meaning: Git metadata (hash: 30dcdee→2cb2e68, mode: 644)

Line 3:  --- a/app/main.go
         ─┬─ ─┬─────────────
          │   │
         Old  Path prefix
         marker

         📄 Meaning: This marks the "before" version of the file

Line 4:  +++ b/app/main.go
         ─┬─ ─┬─────────────
          │   │
         New  Path prefix
         marker

         📄 Meaning: This marks the "after" version of the file

Line 5:  @@ -23,7 +23,8 @@ func main() {
         ─┬─ ──┬─── ──┬─── ─────┬──────
          │    │      │          │
         Hunk  Old    New     Context
         marker range range  (optional)

         📍 Meaning: Change section starts here
            • Old file: line 23, showing 7 lines
            • New file: line 23, showing 8 lines
            • Location: inside "func main()"

Line 6:
         ┬
         │
      Space prefix (context line)

         ⚪ Meaning: Blank line that exists in both versions (unchanged)

Line 7:      image := os.Args[2]
         ┬───────────────────────
         │
      Space prefix (context line)

         ⚪ Meaning: This line is unchanged (shown for context)

Line 8:
         ┬
         │
      Space prefix

         ⚪ Meaning: Another blank line (unchanged)

Line 9:  -    EnterNewJail(os.Args[3], image)
         ┬───────────────────────────────────
         │
      Minus prefix (removed line)

         ❌ Meaning: This line was DELETED from the code

Line 10: +    jailPath := EnterNewJail(os.Args[3], image)
         ┬──────────────────────────────────────────────
         │
      Plus prefix (added line)

         ✅ Meaning: This line was ADDED (replaces line 9)

Line 11: +    defer os.Remove(jailPath)
         ┬────────────────────────────
         │
      Plus prefix (added line)

         ✅ Meaning: This is another ADDED line (brand new)

Line 12:
         ┬
         │
      Space prefix

         ⚪ Meaning: Blank line (unchanged)

Line 13:     cmd := exec.Command(command, args...)
         ┬──────────────────────────────────────
         │
      Space prefix

         ⚪ Meaning: This line is unchanged (shown for context)
```

---

## The 3 Line Types

Every line in the change section (after `@@`) has ONE of these prefixes:

| Prefix | Type | Visual | Meaning |
|--------|------|--------|---------|
| ` ` (space) | **Context** | ⚪ | Line exists in both old and new (unchanged) |
| `-` | **Removed** | ❌ | Line existed in old, deleted in new |
| `+` | **Added** | ✅ | Line doesn't exist in old, added in new |

---

## What Changed in Your Example?

### Before (Old Code)
```go
// Line 23
    image := os.Args[2]

    EnterNewJail(os.Args[3], image)  // ← This line removed

    cmd := exec.Command(command, args...)
```

### After (New Code)
```go
// Line 23
    image := os.Args[2]

    jailPath := EnterNewJail(os.Args[3], image)  // ← Added (modified)
    defer os.Remove(jailPath)                     // ← Added (new)

    cmd := exec.Command(command, args...)
```

### Summary
- **Removed:** 1 line (old function call)
- **Added:** 2 lines (capture return value + cleanup)
- **Net change:** +1 line
- **Purpose:** Changed void function call to capture return value and add cleanup

---

## Understanding the Hunk Header

This is the most important line to understand:

```
@@ -23,7 +23,8 @@ func main() {
```

Breaking it down:

```
@@ -23,7 +23,8 @@ func main() {
   │  │  │  │  │  └─────── Optional: shows where in code (function name)
   │  │  │  │  └────────── New file: display 8 lines
   │  │  │  └───────────── New file: starting at line 23
   │  │  └──────────────── Old file: display 7 lines
   │  └─────────────────── Old file: starting at line 23
   └────────────────────── Hunk marker (always @@)
```

**Simple formula:** `@@ -old_start,old_count +new_start,new_count @@`

**Why 7 vs 8 lines?**
- Old version shows: 3 context + 1 removed + 3 context = **7 lines**
- New version shows: 3 context + 2 added + 3 context = **8 lines**

---

## Different Sections Explained

### 1. File Header (Line 1)
```diff
diff --git a/app/main.go b/app/main.go
```
- Identifies which file was modified
- `a/` prefix = old version
- `b/` prefix = new version
- If paths differ, file was renamed

### 2. Index Line (Line 2)
```diff
index 30dcdee..2cb2e68 100644
```
- Git-specific metadata
- Object hashes: `old_hash..new_hash`
- File mode: `100644` = regular file (644 permissions)
- Safe to ignore when parsing

### 3. File Markers (Lines 3-4)
```diff
--- a/app/main.go
+++ b/app/main.go
```
- Traditional unified diff format
- `---` marks the old file
- `+++` marks the new file
- These are standard in all diffs

### 4. Hunk Header (Line 5)
```diff
@@ -23,7 +23,8 @@ func main() {
```
- Marks the start of a change section (hunk)
- Shows line ranges in old and new files
- Optional context (function/class name)

### 5. Change Lines (Lines 6-13)
```diff
 context line       ← Space prefix
-removed line       ← Minus prefix
+added line         ← Plus prefix
```
- Three types: context, removed, added
- Context lines provide surrounding code
- Changes show what was modified

---

## How to Parse This in Go

### Already Implemented!

The parser at `internal/diff/parser.go` handles all of this:

```go
import "github.com/shreybhardwaj/zizou/internal/diff"

// Parse the diff
parser := diff.NewParser()
parsedDiff, _ := parser.Parse(diffContent)

// Extract information
for _, file := range parsedDiff.Files {
    fmt.Printf("File: %s\n", file.NewPath)

    for _, hunk := range file.Hunks {
        fmt.Printf("  Lines %d-%d\n", hunk.NewStart, hunk.NewStart+hunk.NewCount)

        for _, line := range hunk.Lines {
            switch line.Type {
            case diff.LineAdded:
                fmt.Printf("    + %s\n", line.Content)
            case diff.LineRemoved:
                fmt.Printf("    - %s\n", line.Content)
            case diff.LineContext:
                fmt.Printf("      %s\n", line.Content)
            }
        }
    }
}
```

### Quick Utilities

```go
// Get statistics
stats := parsedDiff.Stats()
fmt.Printf("Files: %d, +%d -%d\n",
    stats.Files, stats.LinesAdded, stats.LinesRemoved)

// Get only added lines
for _, line := range parsedDiff.GetAddedLines() {
    fmt.Printf("%s:%d: %s\n", line.File, line.Line, line.Content)
}

// Get only removed lines
for _, line := range parsedDiff.GetRemovedLines() {
    fmt.Printf("%s: %s\n", line.File, line.Content)
}

// Get modified files
files := parsedDiff.GetModifiedFiles()
```

---

## Try It Yourself

### Run the Interactive Demo
```bash
go run examples/line_by_line_demo.go
```

This will show your exact diff with explanations for each line!

### Output Preview
```
═══════════════════════════════════════════════
     GIT DIFF LINE-BY-LINE EXPLANATION
═══════════════════════════════════════════════

Line  1: diff --git a/app/main.go b/app/main.go
         📁 FILE HEADER: Comparing 'app/main.go' (old) with 'app/main.go' (new)

Line  5: @@ -23,7 +23,8 @@ func main() {
         📍 HUNK HEADER: Old file: line 23, 7 lines shown |
                         New file: line 23, 8 lines shown

Line  9: -    EnterNewJail(os.Args[3], image)
         ❌ REMOVED: '    EnterNewJail(os.Args[3], image)'

Line 10: +    jailPath := EnterNewJail(os.Args[3], image)
         ✅ ADDED: '    jailPath := EnterNewJail(os.Args[3], image)'
...
```

---

## Additional Resources

📖 **Detailed Guides:**
- [Quick Reference Card](docs/QUICK_REFERENCE.md) - Fast lookup
- [Line-by-Line Explanation](docs/DIFF_FORMAT_EXPLAINED.md) - Complete breakdown
- [Diff Parsing Guide](docs/DIFF_PARSING.md) - Technical implementation

🎯 **Examples:**
- `examples/line_by_line_demo.go` - Interactive explanation ⭐
- `examples/simple_usage.go` - Common use cases
- `examples/parse_demo.go` - Comprehensive demo

🔧 **Parser Code:**
- `internal/diff/types.go` - Data structures
- `internal/diff/parser.go` - Parsing logic
- `internal/diff/helpers.go` - Utility methods

---

## Quick Summary

Your diff shows a **modification** to `app/main.go`:

1. **File changed:** `app/main.go` (same file, modified in place)
2. **Location:** Around line 23, inside `func main()`
3. **Old code:** Simple function call (1 line)
4. **New code:** Capture return value + cleanup (2 lines)
5. **Result:** Net +1 line added

**The change:** Refactored `EnterNewJail()` to return a path that gets cleaned up with `defer`.

---

**Need more help?** Check `docs/README.md` for the complete documentation index!
