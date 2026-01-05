# Edge Cases Support

Zizou's diff parser now handles all common edge cases found in git diffs.

## Supported Edge Cases

### ✅ 1. Binary Files

**What it is:** Files that contain non-text data (images, executables, etc.)

**Diff format:**
```diff
diff --git a/image.png b/image.png
index abc123..def456 100644
Binary files a/image.png and b/image.png differ
```

**Detected fields:**
- `IsBinary = true`
- No hunks (binary files don't show line changes)

**Use case:** Detect when binary assets have been modified

---

### ✅ 2. File Renames

**What it is:** File moved to a new location or renamed

**Diff format:**
```diff
diff --git a/oldname.go b/newname.go
similarity index 95%
rename from oldname.go
rename to newname.go
index abc123..def456 100644
--- a/oldname.go
+++ b/newname.go
```

**Detected fields:**
- `IsRenamed = true`
- `OldPath = "oldname.go"`
- `NewPath = "newname.go"`
- `Similarity = 95` (percentage)

**Use case:** Track file organization changes, detect refactoring

---

### ✅ 3. File Copies

**What it is:** File duplicated to a new location

**Diff format:**
```diff
diff --git a/original.go b/copy.go
similarity index 100%
copy from original.go
copy to copy.go
```

**Detected fields:**
- `IsCopied = true`
- `OldPath = "original.go"`
- `NewPath = "copy.go"`
- `Similarity = 100`

**Use case:** Detect code duplication, template usage

---

### ✅ 4. File Mode Changes

**What it is:** Permission changes (e.g., making a file executable)

**Diff format:**
```diff
diff --git a/script.sh b/script.sh
old mode 100644
new mode 100755
index abc123..abc123 100755
--- a/script.sh
+++ b/script.sh
```

**Detected fields:**
- `OldMode = "100644"` (rw-r--r--)
- `NewMode = "100755"` (rwxr-xr-x)

**Common modes:**
- `100644` - Regular file
- `100755` - Executable file
- `120000` - Symbolic link
- `040000` - Directory

**Use case:** Detect security changes, script permissions

---

### ✅ 5. New Files

**What it is:** File created in the diff

**Diff format:**
```diff
diff --git a/newfile.go b/newfile.go
new file mode 100644
index 0000000..abc123
--- /dev/null
+++ b/newfile.go
@@ -0,0 +1,3 @@
+package main
+func main() {
+}
```

**Detected fields:**
- `IsNew = true`
- `NewMode = "100644"`
- Old path markers show `/dev/null`

**Use case:** Track new functionality, code additions

---

### ✅ 6. Deleted Files

**What it is:** File removed in the diff

**Diff format:**
```diff
diff --git a/oldfile.go b/oldfile.go
deleted file mode 100644
index abc123..0000000
--- a/oldfile.go
+++ /dev/null
@@ -1,3 +0,0 @@
-package main
-func main() {
-}
```

**Detected fields:**
- `IsDeleted = true`
- `OldMode = "100644"`
- New path markers show `/dev/null`

**Use case:** Track removed functionality, cleanup operations

---

### ✅ 7. Renamed AND Modified

**What it is:** File renamed and also has content changes

**Diff format:**
```diff
diff --git a/old.go b/new.go
similarity index 87%
rename from old.go
rename to new.go
index abc123..def456 100644
--- a/old.go
+++ b/new.go
@@ -1,3 +1,4 @@
 package main
+import "fmt"
 func main() {
```

**Detected fields:**
- `IsRenamed = true`
- `Similarity = 87`
- Has hunks showing changes

**Use case:** Detect refactoring with modifications

---

### ✅ 8. Mode Change + Modification

**What it is:** File permissions changed and content modified

**Diff format:**
```diff
diff --git a/script.sh b/script.sh
old mode 100644
new mode 100755
index abc123..def456
--- a/script.sh
+++ b/script.sh
@@ -1,2 +1,3 @@
 #!/bin/bash
+echo "Now executable"
```

**Detected fields:**
- `OldMode = "100644"`
- `NewMode = "100755"`
- Has hunks showing changes

---

### ✅ 9. Binary File Rename

**What it is:** Binary file moved/renamed without modification

**Diff format:**
```diff
diff --git a/old.png b/new.png
similarity index 100%
rename from old.png
rename to new.png
```

**Detected fields:**
- `IsRenamed = true`
- `Similarity = 100`
- No hunks (binary file)

---

### ✅ 10. Symbolic Links

**What it is:** Symlink target changes

**Diff format:**
```diff
diff --git a/link.txt b/link.txt
index abc123..def456 120000
--- a/link.txt
+++ b/link.txt
@@ -1 +1 @@
-target1.txt
+target2.txt
```

**Notes:**
- Mode `120000` indicates symlink
- Parsed as regular diff (shows target change)

---

### ✅ 11. Empty Files

**What it is:** File with no content

**Diff format:**
```diff
diff --git a/empty.txt b/empty.txt
new file mode 100644
index 0000000..e69de29
```

**Detected fields:**
- `IsNew = true`
- No hunks (empty file)

---

### ✅ 12. Multiple Files with Mixed Types

**What it is:** Diff containing text, binary, and special files

**Example:**
- Text file modifications
- Binary file changes
- Renamed files
- Mode changes

**All detected correctly** in a single diff

---

## API Usage

### Accessing Edge Case Information

```go
import "github.com/BhardwajShrey/zizou/internal/diff"

parser := diff.NewParser()
d, _ := parser.Parse(diffContent)

for _, file := range d.Files {
    // Check file status
    if file.IsBinary {
        fmt.Printf("Binary file: %s\n", file.NewPath)
    }

    if file.IsRenamed {
        fmt.Printf("Renamed: %s -> %s (%.0f%% similar)\n",
            file.OldPath, file.NewPath, file.Similarity)
    }

    if file.IsNew {
        fmt.Printf("New file: %s (mode: %s)\n",
            file.NewPath, file.NewMode)
    }

    if file.IsDeleted {
        fmt.Printf("Deleted: %s\n", file.OldPath)
    }

    if file.IsCopied {
        fmt.Printf("Copied: %s -> %s\n",
            file.OldPath, file.NewPath)
    }

    if file.OldMode != file.NewMode && file.NewMode != "" {
        fmt.Printf("Mode changed: %s -> %s\n",
            file.OldMode, file.NewMode)
    }
}
```

### Filtering by Type

```go
// Find all binary files
var binaryFiles []diff.FileDiff
for _, file := range d.Files {
    if file.IsBinary {
        binaryFiles = append(binaryFiles, file)
    }
}

// Find all renames
var renames []diff.FileDiff
for _, file := range d.Files {
    if file.IsRenamed {
        renames = append(renames, file)
    }
}

// Find executable files
var executables []diff.FileDiff
for _, file := range d.Files {
    if file.NewMode == "100755" {
        executables = append(executables, file)
    }
}
```

## Testing

All edge cases have comprehensive test coverage:

```bash
# Run edge case tests
go test ./internal/diff -run Advanced

# Run specific edge case test
go test ./internal/diff -run TestParser_Parse_BinaryFile -v
```

**Test files:**
- `edge_cases_test.go` - Basic edge cases
- `advanced_edge_cases_test.go` - Complex edge cases

**Test coverage:** 98.9% of statements

---

## Real-World Examples

### Example 1: Security Review
```go
// Flag executable files added
for _, file := range d.Files {
    if file.IsNew && file.NewMode == "100755" {
        fmt.Printf("⚠️  New executable: %s\n", file.NewPath)
    }
}
```

### Example 2: Refactoring Detection
```go
// Detect major refactorings
for _, file := range d.Files {
    if file.IsRenamed && file.Similarity < 70 {
        fmt.Printf("🔄 Major refactor: %s -> %s (%d%% similar)\n",
            file.OldPath, file.NewPath, file.Similarity)
    }
}
```

### Example 3: Asset Changes
```go
// Track binary asset updates
for _, file := range d.Files {
    if file.IsBinary && !file.IsNew && !file.IsDeleted {
        fmt.Printf("📦 Binary modified: %s\n", file.NewPath)
    }
}
```

### Example 4: Permission Auditing
```go
// Audit permission changes
for _, file := range d.Files {
    if file.OldMode == "100644" && file.NewMode == "100755" {
        fmt.Printf("🔐 File made executable: %s\n", file.NewPath)
    }
}
```

---

## Not Yet Supported

The following are NOT currently supported (but rarely needed):

- ❌ Submodule changes
- ❌ Merge conflict markers
- ❌ Combined diffs (merge commits)
- ❌ Git attributes (`.gitattributes` custom diff drivers)

These can be added if needed for your use case.

---

## Summary

✅ **Fully supported edge cases:**
1. Binary files
2. File renames
3. File copies
4. Mode changes (chmod)
5. New files
6. Deleted files
7. Renamed + modified
8. Mode + modified
9. Binary renames
10. Symbolic links
11. Empty files
12. Mixed file types

**All edge cases are:**
- ✅ Properly parsed
- ✅ Fully tested
- ✅ Accessible via API
- ✅ Documented with examples

The parser is now **production-ready** for all common git diff scenarios!
