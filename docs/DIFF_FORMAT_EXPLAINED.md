# Git Diff Format - Line by Line Explanation

Let's break down your exact diff, line by line:

## Complete Example with Line Numbers

```diff
1  | diff --git a/app/main.go b/app/main.go
2  | index 30dcdee..2cb2e68 100644
3  | --- a/app/main.go
4  | +++ b/app/main.go
5  | @@ -23,7 +23,8 @@ func main() {
6  |
7  |      image := os.Args[2]
8  |
9  | -    EnterNewJail(os.Args[3], image)
10 | +    jailPath := EnterNewJail(os.Args[3], image)
11 | +    defer os.Remove(jailPath)
12 |
13 |      cmd := exec.Command(command, args...)
```

---

## Line-by-Line Breakdown

### Line 1: File Header
```diff
diff --git a/app/main.go b/app/main.go
```

**What it means:**
- `diff --git` - This is a git-format diff
- `a/app/main.go` - Path in the "old" version (prefix `a/`)
- `b/app/main.go` - Path in the "new" version (prefix `b/`)

**Purpose:** Identifies which file changed. If the file was renamed, these paths would differ.

**Example of renamed file:**
```diff
diff --git a/old_name.go b/new_name.go
```

---

### Line 2: Index Line (Git Metadata)
```diff
index 30dcdee..2cb2e68 100644
```

**What it means:**
- `index` - Git index/staging area metadata
- `30dcdee` - Git object hash of old file (abbreviated)
- `2cb2e68` - Git object hash of new file (abbreviated)
- `100644` - Unix file permissions (644 = rw-r--r--)

**File mode meanings:**
- `100644` - Regular file
- `100755` - Executable file
- `040000` - Directory
- `120000` - Symbolic link

**Purpose:** Git-specific metadata. Used for internal Git operations. Safe to ignore when parsing diffs for code review.

---

### Line 3: Old File Marker
```diff
--- a/app/main.go
```

**What it means:**
- `---` - Marks the old/original version
- `a/` - Conventional prefix for "before" state
- `app/main.go` - Path to the file

**Special cases:**
```diff
--- /dev/null         ← File didn't exist (new file)
--- a/deleted.go      ← File will be deleted
```

**Purpose:** Traditional unified diff format marker. Tells diff tools which file is being compared.

---

### Line 4: New File Marker
```diff
+++ b/app/main.go
```

**What it means:**
- `+++` - Marks the new/modified version
- `b/` - Conventional prefix for "after" state
- `app/main.go` - Path to the file

**Special cases:**
```diff
+++ /dev/null         ← File deleted (doesn't exist after)
+++ b/new_file.go     ← New file created
```

**Purpose:** Companion to `---` line. Completes the file identification.

---

### Line 5: Hunk Header
```diff
@@ -23,7 +23,8 @@ func main() {
```

**CRITICAL LINE - Most important to understand!**

Let's break it down:

```
@@ -23,7 +23,8 @@ func main() {
   │  │  │  │  │  └─── Optional: context (function/class name)
   │  │  │  │  └─────── New version: 8 lines shown
   │  │  │  └────────── New version: starts at line 23
   │  │  └───────────── Old version: 7 lines shown
   │  └──────────────── Old version: starts at line 23
   └─────────────────── Hunk marker (always @@)
```

**Detailed breakdown:**
- `@@` - Hunk boundary markers (always come in pairs)
- `-23,7` - In the OLD file:
  - Start at line 23
  - Show 7 lines total
- `+23,8` - In the NEW file:
  - Start at line 23
  - Show 8 lines total
- `@@ func main() {` - Optional context showing where in the code this is

**Why different line counts?**
- Old: 7 lines
- New: 8 lines
- Difference: +1 line (we added more than we removed)

**Understanding the math:**
```
Old file lines shown:
  3 context lines (unchanged)
  1 removed line (-)
  3 more context lines
  = 7 total

New file lines shown:
  3 context lines (unchanged)
  2 added lines (+)
  3 more context lines
  = 8 total
```

---

### Line 6: Empty Context Line
```diff

```

**What it means:**
- Single space ` ` - This is a context line
- Empty content - This was a blank line in the code

**Purpose:** Shows unchanged code around the changes. Context helps you understand where the change occurred.

**Important:** Context lines start with a SPACE character, not nothing!

---

### Line 7: Context Line
```diff
     image := os.Args[2]
```

**What it means:**
- ` ` (space prefix) - Unchanged line
- `    image := os.Args[2]` - The actual code

**Purpose:** Context. This line existed before and still exists after. Helps orient you to the location of changes.

---

### Line 8: Empty Context Line
```diff

```

**What it means:**
- Space prefix - Context line
- Blank line in the code

---

### Line 9: Removed Line
```diff
-    EnterNewJail(os.Args[3], image)
```

**What it means:**
- `-` (minus prefix) - This line was REMOVED
- `    EnterNewJail(os.Args[3], image)` - The code that was deleted

**Purpose:** Shows what code was removed from the old version.

**Visual indicator:**
- Many tools show this in RED
- Indicates deletion

---

### Line 10: Added Line
```diff
+    jailPath := EnterNewJail(os.Args[3], image)
```

**What it means:**
- `+` (plus prefix) - This line was ADDED
- `    jailPath := EnterNewJail(os.Args[3], image)` - The new code

**Purpose:** Shows what code was added in the new version.

**Visual indicator:**
- Many tools show this in GREEN
- Indicates addition

---

### Line 11: Added Line
```diff
+    defer os.Remove(jailPath)
```

**What it means:**
- `+` prefix - Another added line
- `    defer os.Remove(jailPath)` - More new code

**Purpose:** Second line added in this change.

---

### Line 12: Empty Context Line
```diff

```

**What it means:**
- Space prefix - Context line
- Blank line

---

### Line 13: Context Line
```diff
     cmd := exec.Command(command, args...)
```

**What it means:**
- Space prefix - Unchanged line
- Shows context after the changes

**Purpose:** Provides context showing what comes after the modification.

---

## Visual Summary

Here's how to read the prefixes:

```diff
 line    ← Space = CONTEXT (unchanged, shown for reference)
-line    ← Minus = REMOVED (existed before, deleted now)
+line    ← Plus  = ADDED (new code, didn't exist before)
```

## Full Example with Annotations

```diff
diff --git a/app/main.go b/app/main.go     ← Which file
index 30dcdee..2cb2e68 100644              ← Git metadata
--- a/app/main.go                          ← Old version marker
+++ b/app/main.go                          ← New version marker
@@ -23,7 +23,8 @@ func main() {           ← Hunk: lines 23-29 (old) → 23-30 (new)

     image := os.Args[2]                   ← CONTEXT (unchanged)

-    EnterNewJail(os.Args[3], image)       ← REMOVED (1 line deleted)
+    jailPath := EnterNewJail(os.Args[3], image)  ← ADDED (replacement)
+    defer os.Remove(jailPath)             ← ADDED (new line)

     cmd := exec.Command(command, args...) ← CONTEXT (unchanged)
```

## What Changed?

**Old code (1 line):**
```go
EnterNewJail(os.Args[3], image)
```

**New code (2 lines):**
```go
jailPath := EnterNewJail(os.Args[3], image)
defer os.Remove(jailPath)
```

**Summary:** Changed function call to capture return value + added cleanup with defer.

---

## Common Patterns

### Adding a New Line
```diff
 existing line
+new line
 existing line
```

### Removing a Line
```diff
 existing line
-deleted line
 existing line
```

### Modifying a Line
```diff
 existing line
-old version of line
+new version of line
 existing line
```

### Adding Multiple Lines
```diff
 existing line
+first new line
+second new line
+third new line
 existing line
```

---

## Multiple Hunks Example

A single file can have multiple separate changes (hunks):

```diff
diff --git a/app/main.go b/app/main.go
--- a/app/main.go
+++ b/app/main.go
@@ -23,7 +23,8 @@ func main() {          ← FIRST HUNK (lines 23-30)
     image := os.Args[2]
-    EnterNewJail(os.Args[3], image)
+    jailPath := EnterNewJail(os.Args[3], image)
+    defer os.Remove(jailPath)
     cmd := exec.Command(command, args...)
@@ -37,14 +38,12 @@ func main() {          ← SECOND HUNK (lines 37-50)
     os.Exit(cmd.ProcessState.ExitCode())
 }
-func EnterNewJail(filepath string, image string) {
+func EnterNewJail(filepath string, image string) string {
```

**Why multiple hunks?**
- Changes are in different parts of the file
- Too much unchanged code between changes
- Git groups nearby changes together

---

## Edge Cases

### New File
```diff
diff --git a/new_file.go b/new_file.go
new file mode 100644                ← Special marker
index 0000000..a1b2c3d
--- /dev/null                       ← Didn't exist before
+++ b/new_file.go
@@ -0,0 +1,5 @@                     ← Old: 0 lines, New: 5 lines
+package main
+
+func NewFunction() {
+    // code
+}
```

### Deleted File
```diff
diff --git a/old_file.go b/old_file.go
deleted file mode 100644            ← Special marker
index a1b2c3d..0000000
--- a/old_file.go
+++ /dev/null                       ← Doesn't exist after
@@ -1,10 +0,0 @@                    ← Old: 10 lines, New: 0 lines
-package main
-
-func OldFunction() {
-    // code that's being removed
-}
```

### Binary File
```diff
diff --git a/image.png b/image.png
index abc123..def456 100644
Binary files a/image.png and b/image.png differ
```

### Renamed File
```diff
diff --git a/old_name.go b/new_name.go
similarity index 95%               ← 95% similar
rename from old_name.go
rename to new_name.go
index abc123..def456 100644
--- a/old_name.go
+++ b/new_name.go
@@ -1,3 +1,3 @@
-package oldpackage
+package newpackage
```

---

## Quick Reference Card

| Prefix | Meaning | Description |
|--------|---------|-------------|
| `diff --git` | File header | Which file changed |
| `index` | Git metadata | Object hashes and permissions |
| `---` | Old file | "Before" version |
| `+++` | New file | "After" version |
| `@@` | Hunk header | Line ranges for this section |
| ` ` (space) | Context | Unchanged line |
| `-` | Removed | Deleted line |
| `+` | Added | New line |

---

## How Zizou Parses This

1. **Regex matching** - Identifies each line type
2. **State machine** - Tracks current file/hunk context
3. **Data structure** - Builds nested Diff → FileDiff → Hunk → Line
4. **Line numbers** - Calculates accurate positions

See `internal/diff/parser.go:41` for the implementation!
