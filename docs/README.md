# Zizou Documentation

Complete documentation for understanding and using Zizou's diff parser.

## 📚 Documentation Index

### For Beginners

**Start here if you're new to git diffs:**

1. **[Quick Reference Card](QUICK_REFERENCE.md)** ⚡
   - 3-second overview
   - Line prefix cheat sheet
   - Common patterns
   - Memory tricks

2. **[Line-by-Line Explanation](DIFF_FORMAT_EXPLAINED.md)** 📖
   - Complete breakdown of every line
   - Visual examples with annotations
   - Edge cases (new files, deletions, renames)
   - Multiple hunks explained

### For Developers

**Building with Zizou's parser:**

3. **[Diff Parsing Guide](DIFF_PARSING.md)** 🔧
   - Parser architecture
   - Data structures
   - Usage examples
   - Integration guide

## 🎯 Quick Navigation

| I want to... | Go to... |
|-------------|----------|
| Understand diff format basics | [Quick Reference](QUICK_REFERENCE.md) |
| Learn what each line means | [Line-by-Line](DIFF_FORMAT_EXPLAINED.md) |
| Use the parser in my code | [Parsing Guide](DIFF_PARSING.md) |
| See working examples | `../examples/` |
| Run interactive demo | `go run examples/line_by_line_demo.go` |

## 📝 Diff Format Summary

### The Essentials

```diff
diff --git a/file.go b/file.go     ← File header
--- a/file.go                      ← Old version
+++ b/file.go                      ← New version
@@ -23,7 +23,8 @@                  ← Hunk header (line ranges)
 context                            ← Space = unchanged
-removed                            ← Minus = deleted
+added                              ← Plus = new
```

### Key Concepts

1. **File Headers** - Identify which files changed
2. **Hunks** - Sections of changes within a file
3. **Line Prefixes** - Tell you what happened to each line
   - ` ` (space) = Context (unchanged)
   - `-` = Removed (deleted)
   - `+` = Added (new)

## 🛠️ Parser Architecture

```
Input: Git diff string
  ↓
Parser (regex + state machine)
  ↓
Diff structure
  ├── FileDiff (per file)
  │   └── Hunk (per change section)
  │       └── Line (per line)
  ↓
Helper methods
  ├── GetAddedLines()
  ├── GetRemovedLines()
  ├── GetModifiedFiles()
  └── Stats()
```

## 💡 Examples

### Interactive Learning
```bash
# Visual line-by-line explanation
go run examples/line_by_line_demo.go
```

### Programmatic Usage
```go
import "github.com/shreybhardwaj/zizou/internal/diff"

parser := diff.NewParser()
d, _ := parser.Parse(diffContent)

// Get statistics
stats := d.Stats()
fmt.Printf("+%d -%d across %d files\n",
    stats.LinesAdded,
    stats.LinesRemoved,
    stats.Files)

// Find security issues in added code
for _, line := range d.GetAddedLines() {
    if strings.Contains(line.Content, "exec.Command") {
        fmt.Printf("⚠️  %s:%d uses exec.Command\n",
            line.File, line.Line)
    }
}
```

## 📖 Detailed Topics

### Understanding Hunks

A **hunk** is a continuous section of changes. One file can have multiple hunks if changes are in different parts of the file.

```diff
@@ -23,7 +23,8 @@
      ││  │  ││  │
      ││  │  ││  └─ New: 8 lines shown
      ││  │  │└──── New: start at line 23
      ││  │  └───── Old: 7 lines shown
      │└──┴──────── Old: start at line 23
      └───────────── Hunk markers
```

**Rule:** `-old_start,old_lines +new_start,new_lines`

### Line Numbering

- **Context lines** (`  `) keep the same line number
- **Added lines** (`+`) get new line numbers
- **Removed lines** (`-`) reference old line numbers

The parser tracks this automatically!

### Edge Cases

1. **New file:** `--- /dev/null`
2. **Deleted file:** `+++ /dev/null`
3. **Binary files:** `Binary files differ`
4. **Renamed:** `rename from/to`
5. **No newline:** `\ No newline at end of file`

See [DIFF_FORMAT_EXPLAINED.md](DIFF_FORMAT_EXPLAINED.md) for complete coverage.

## 🔍 Real-World Example

**Your example diff explained:**

```diff
# Line 1: File being compared
diff --git a/app/main.go b/app/main.go

# Lines 2-4: Metadata and file markers
index 30dcdee..2cb2e68 100644
--- a/app/main.go
+++ b/app/main.go

# Line 5: This hunk shows lines 23-29 → 23-30
@@ -23,7 +23,8 @@ func main() {

     image := os.Args[2]              # Context (unchanged)

-    EnterNewJail(os.Args[3], image)  # Removed (old code)
+    jailPath := EnterNewJail(...)    # Added (new code)
+    defer os.Remove(jailPath)        # Added (new code)

     cmd := exec.Command(...)         # Context (unchanged)
```

**What changed:**
- Removed: 1 line (function call)
- Added: 2 lines (capture return value + cleanup)
- Net: +1 line

## 🚀 Next Steps

1. **Learn the format:**
   ```bash
   go run examples/line_by_line_demo.go
   ```

2. **Read the guide:**
   - [Quick Reference](QUICK_REFERENCE.md) for basics
   - [Line-by-Line](DIFF_FORMAT_EXPLAINED.md) for deep dive

3. **Try the parser:**
   ```bash
   go run examples/simple_usage.go
   ```

4. **Build something:**
   - See [Diff Parsing Guide](DIFF_PARSING.md)
   - Check `internal/diff/parser.go`

## 📚 Additional Resources

- **Parser Implementation:** `../internal/diff/`
  - `types.go` - Data structures
  - `parser.go` - Parsing logic
  - `helpers.go` - Utility methods

- **Examples:** `../examples/`
  - `line_by_line_demo.go` - Interactive explanation
  - `simple_usage.go` - Common use cases
  - `parse_demo.go` - Comprehensive demo

- **Integration:** See `cmd/root.go` for how Zizou uses the parser

## ❓ FAQ

**Q: Why do some hunks skip line numbers?**
A: Git only shows changed sections plus a few lines of context. Large unchanged sections are omitted.

**Q: What's the difference between `---` and `-`?**
A: Three dashes (`---`) mark the old file. Single dash (`-`) marks a removed line.

**Q: Can one file have multiple hunks?**
A: Yes! If changes are in different parts of the file, each section gets its own hunk.

**Q: How do I test the parser?**
A: Run any example with your own diffs:
```bash
git diff | go run ../main.go --file -
```

## 🤝 Contributing

Found an issue or want to improve the docs? Check `../README.md` for contribution guidelines.

---

**Legend:**
- ⚡ Quick reference
- 📖 In-depth explanation
- 🔧 Technical guide
- 💡 Examples
- ❓ Help
