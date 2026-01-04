# Git Diff Quick Reference Card

## The 3-Second Overview

```diff
diff --git a/file.go b/file.go     ← Which file
--- a/file.go                      ← Old version
+++ b/file.go                      ← New version
@@ -23,7 +23,8 @@                  ← Line ranges: old(23-29) → new(23-30)
 context line                       ← Space = unchanged
-removed line                       ← Minus = deleted
+added line                         ← Plus = new code
```

## Line Prefixes Cheat Sheet

| First Character | Meaning | Visual |
|----------------|---------|--------|
| `d` (diff) | File header | 📁 Which file changed |
| `i` (index) | Git metadata | 🔖 Hashes & permissions |
| `-` (3 times) | Old file marker | 📄 Before state |
| `+` (3 times) | New file marker | 📄 After state |
| `@` (twice) | Hunk header | 📍 Line ranges |
| ` ` (space) | Context line | ⚪ Unchanged |
| `-` (once) | Removed | ❌ Deleted |
| `+` (once) | Added | ✅ New |

## Hunk Header Decoder

```
@@ -23,7 +23,8 @@ func main() {
   │  │  │  │  │  └─ Optional: where in code
   │  │  │  │  └──── New: show 8 lines
   │  │  │  └─────── New: start line 23
   │  │  └────────── Old: show 7 lines
   │  └───────────── Old: start line 23
   └──────────────── Markers
```

**Simple rule:** `-old_start,old_count +new_start,new_count`

## Common Patterns

### Simple Addition
```diff
 existing code
+new line
 existing code
```

### Simple Deletion
```diff
 existing code
-deleted line
 existing code
```

### Line Modification
```diff
 existing code
-old version
+new version
 existing code
```

### Block Addition
```diff
 existing code
+line 1
+line 2
+line 3
 existing code
```

## File States

### New File Created
```diff
--- /dev/null          ← Didn't exist
+++ b/new_file.go
@@ -0,0 +1,5 @@       ← 0 lines → 5 lines
```

### File Deleted
```diff
--- a/old_file.go
+++ /dev/null          ← Doesn't exist anymore
@@ -1,10 +0,0 @@      ← 10 lines → 0 lines
```

### File Renamed
```diff
diff --git a/old.go b/new.go
rename from old.go
rename to new.go
```

## Real Example Annotated

```diff
Line 1:  diff --git a/app/main.go b/app/main.go
         ↑ File being compared

Line 2:  index 30dcdee..2cb2e68 100644
         ↑ Git hashes (old→new) + file mode

Line 3:  --- a/app/main.go
         ↑ Old file marker

Line 4:  +++ b/app/main.go
         ↑ New file marker

Line 5:  @@ -23,7 +23,8 @@ func main() {
         ↑ Change section: lines 23-29 → 23-30

Line 6:
         ↑ (space prefix) Context: blank line

Line 7:      image := os.Args[2]
         ↑ (space prefix) Context: unchanged

Line 8:
         ↑ (space prefix) Context: blank line

Line 9:  -    EnterNewJail(os.Args[3], image)
         ↑ (minus prefix) REMOVED this line

Line 10: +    jailPath := EnterNewJail(os.Args[3], image)
         ↑ (plus prefix) ADDED new line

Line 11: +    defer os.Remove(jailPath)
         ↑ (plus prefix) ADDED another line

Line 12:
         ↑ (space prefix) Context: blank line

Line 13:     cmd := exec.Command(command, args...)
         ↑ (space prefix) Context: unchanged
```

## What Changed?

**Before (1 line):**
```go
EnterNewJail(os.Args[3], image)
```

**After (2 lines):**
```go
jailPath := EnterNewJail(os.Args[3], image)
defer os.Remove(jailPath)
```

**Change type:** Modified call to capture return value + added cleanup

## Key Insights

1. **Context lines** (space prefix) show unchanged code around changes
2. **Hunk header** tells you exactly where in the file the change is
3. **Line counts** in hunk header can differ if lines were added/removed
4. **Each hunk** is a separate section of changes (files can have multiple hunks)

## Try It Yourself

```bash
# Run the interactive demo
go run examples/line_by_line_demo.go

# Test with real diffs
git diff HEAD~1 HEAD
git diff --staged
git show <commit-hash>
```

## Memory Tricks

- **3 dashes (`---`)** = OLD (shorter word, goes first)
- **3 pluses (`+++`)** = NEW (longer word, comes second)
- **Space** = Stay the same (neutral)
- **Minus (-)** = Remove/subtract
- **Plus (+)** = Add/include
- **@@** = At this location

## Parsing Pseudocode

```
for each line in diff:
    if starts with "diff --git":
        → new file section
    elif starts with "---":
        → old file path
    elif starts with "+++":
        → new file path
    elif starts with "@@":
        → new hunk (change section)
    elif starts with "-":
        → line was removed
    elif starts with "+":
        → line was added
    elif starts with " ":
        → context (unchanged)
```

## File Mode Reference

| Mode | Meaning |
|------|---------|
| 100644 | Regular file (rw-r--r--) |
| 100755 | Executable (rwxr-xr-x) |
| 120000 | Symbolic link |
| 040000 | Directory |

## Edge Cases to Watch

1. **Empty lines** still have space prefix
2. **Lines starting with +/-** in actual code need escaping
3. **Binary files** show "Binary files differ" instead of content
4. **No newline at end** marked with "\ No newline at end of file"
5. **Tabs vs spaces** matters for exact matching

## Resources

- Full explanation: `docs/DIFF_FORMAT_EXPLAINED.md`
- Parser code: `internal/diff/parser.go`
- Examples: `examples/`
- Run demo: `go run examples/line_by_line_demo.go`
