package diff

import (
	"testing"
)

func TestParser_Parse_BinaryFile(t *testing.T) {
	diffContent := `diff --git a/image.png b/image.png
index abc123..def456 100644
Binary files a/image.png and b/image.png differ`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(result.Files) != 1 {
		t.Fatalf("Files count = %d, want 1", len(result.Files))
	}

	file := result.Files[0]

	if !file.IsBinary {
		t.Error("IsBinary = false, want true")
	}

	if file.NewPath != "image.png" {
		t.Errorf("NewPath = %q, want %q", file.NewPath, "image.png")
	}

	// Binary files shouldn't have hunks
	if len(file.Hunks) != 0 {
		t.Errorf("Hunks count = %d, want 0 for binary file", len(file.Hunks))
	}
}

func TestParser_Parse_RenamedFileWithSimilarity(t *testing.T) {
	diffContent := `diff --git a/oldname.go b/newname.go
similarity index 95%
rename from oldname.go
rename to newname.go
index abc123..def456 100644
--- a/oldname.go
+++ b/newname.go
@@ -1,2 +1,3 @@
 package main
+// renamed file
 func main() {}`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	file := result.Files[0]

	if !file.IsRenamed {
		t.Error("IsRenamed = false, want true")
	}

	if file.OldPath != "oldname.go" {
		t.Errorf("OldPath = %q, want %q", file.OldPath, "oldname.go")
	}

	if file.NewPath != "newname.go" {
		t.Errorf("NewPath = %q, want %q", file.NewPath, "newname.go")
	}

	if file.Similarity != 95 {
		t.Errorf("Similarity = %d, want 95", file.Similarity)
	}
}

func TestParser_Parse_CopiedFile(t *testing.T) {
	diffContent := `diff --git a/original.go b/copy.go
similarity index 100%
copy from original.go
copy to copy.go`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	file := result.Files[0]

	if !file.IsCopied {
		t.Error("IsCopied = false, want true")
	}

	if file.OldPath != "original.go" {
		t.Errorf("OldPath = %q, want %q", file.OldPath, "original.go")
	}

	if file.NewPath != "copy.go" {
		t.Errorf("NewPath = %q, want %q", file.NewPath, "copy.go")
	}

	if file.Similarity != 100 {
		t.Errorf("Similarity = %d, want 100", file.Similarity)
	}
}

func TestParser_Parse_ModeChange(t *testing.T) {
	diffContent := `diff --git a/script.sh b/script.sh
old mode 100644
new mode 100755
index abc123..abc123 100755
--- a/script.sh
+++ b/script.sh
@@ -1,2 +1,2 @@
 #!/bin/bash
-echo "hello"
+echo "Hello World"`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	file := result.Files[0]

	if file.OldMode != "100644" {
		t.Errorf("OldMode = %q, want %q", file.OldMode, "100644")
	}

	if file.NewMode != "100755" {
		t.Errorf("NewMode = %q, want %q", file.NewMode, "100755")
	}
}

func TestParser_Parse_NewFileWithMode(t *testing.T) {
	diffContent := `diff --git a/newscript.sh b/newscript.sh
new file mode 100755
index 0000000..abc123
--- /dev/null
+++ b/newscript.sh
@@ -0,0 +1,3 @@
+#!/bin/bash
+echo "New script"
+exit 0`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	file := result.Files[0]

	if !file.IsNew {
		t.Error("IsNew = false, want true")
	}

	if file.NewMode != "100755" {
		t.Errorf("NewMode = %q, want %q", file.NewMode, "100755")
	}
}

func TestParser_Parse_DeletedFileWithMode(t *testing.T) {
	diffContent := `diff --git a/oldscript.sh b/oldscript.sh
deleted file mode 100755
index abc123..0000000
--- a/oldscript.sh
+++ /dev/null
@@ -1,3 +0,0 @@
-#!/bin/bash
-echo "Old script"
-exit 0`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	file := result.Files[0]

	if !file.IsDeleted {
		t.Error("IsDeleted = false, want true")
	}

	if file.OldMode != "100755" {
		t.Errorf("OldMode = %q, want %q", file.OldMode, "100755")
	}
}

func TestParser_Parse_RenamedAndModified(t *testing.T) {
	diffContent := `diff --git a/old.go b/new.go
similarity index 87%
rename from old.go
rename to new.go
index abc123..def456 100644
--- a/old.go
+++ b/new.go
@@ -1,5 +1,6 @@
 package main

+import "fmt"
 func main() {
-    println("old")
+    fmt.Println("new")
 }`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	file := result.Files[0]

	if !file.IsRenamed {
		t.Error("IsRenamed = false, want true")
	}

	if file.Similarity != 87 {
		t.Errorf("Similarity = %d, want 87", file.Similarity)
	}

	// Should also have hunks showing the modifications
	if len(file.Hunks) == 0 {
		t.Error("Expected hunks for renamed and modified file")
	}
}

func TestParser_Parse_BinaryFileRename(t *testing.T) {
	diffContent := `diff --git a/old-image.png b/new-image.png
similarity index 100%
rename from old-image.png
rename to new-image.png`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	file := result.Files[0]

	if !file.IsRenamed {
		t.Error("IsRenamed = false, want true")
	}

	if file.Similarity != 100 {
		t.Errorf("Similarity = %d, want 100", file.Similarity)
	}
}

func TestParser_Parse_ComplexMultipleEdgeCases(t *testing.T) {
	diffContent := `diff --git a/old.sh b/new.sh
similarity index 80%
rename from old.sh
rename to new.sh
old mode 100644
new mode 100755
index abc123..def456
--- a/old.sh
+++ b/new.sh
@@ -1,3 +1,4 @@
 #!/bin/bash
+# Modified and renamed
 echo "test"
 exit 0`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	file := result.Files[0]

	// Check all edge cases are detected
	if !file.IsRenamed {
		t.Error("IsRenamed = false, want true")
	}

	if file.Similarity != 80 {
		t.Errorf("Similarity = %d, want 80", file.Similarity)
	}

	if file.OldMode != "100644" {
		t.Errorf("OldMode = %q, want %q", file.OldMode, "100644")
	}

	if file.NewMode != "100755" {
		t.Errorf("NewMode = %q, want %q", file.NewMode, "100755")
	}

	if file.OldPath != "old.sh" {
		t.Errorf("OldPath = %q, want %q", file.OldPath, "old.sh")
	}

	if file.NewPath != "new.sh" {
		t.Errorf("NewPath = %q, want %q", file.NewPath, "new.sh")
	}
}

func TestParser_Parse_SymlinkChange(t *testing.T) {
	// Symlinks appear similar to regular file diffs
	diffContent := `diff --git a/link.txt b/link.txt
index abc123..def456 120000
--- a/link.txt
+++ b/link.txt
@@ -1 +1 @@
-target1.txt
+target2.txt`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Should parse as regular diff - mode 120000 indicates symlink
	if len(result.Files) != 1 {
		t.Fatalf("Files count = %d, want 1", len(result.Files))
	}
}

func TestParser_Parse_MixedBinaryAndTextFiles(t *testing.T) {
	diffContent := `diff --git a/file.txt b/file.txt
index abc123..def456 100644
--- a/file.txt
+++ b/file.txt
@@ -1,2 +1,3 @@
 line1
+line2
 line3
diff --git a/image.png b/image.png
index 111222..333444 100644
Binary files a/image.png and b/image.png differ
diff --git a/another.txt b/another.txt
index 555666..777888 100644
--- a/another.txt
+++ b/another.txt
@@ -1,1 +1,2 @@
 text
+more text`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(result.Files) != 3 {
		t.Fatalf("Files count = %d, want 3", len(result.Files))
	}

	// First file: text
	if result.Files[0].IsBinary {
		t.Error("First file should not be binary")
	}

	// Second file: binary
	if !result.Files[1].IsBinary {
		t.Error("Second file should be binary")
	}

	if result.Files[1].NewPath != "image.png" {
		t.Errorf("Binary file path = %q, want %q", result.Files[1].NewPath, "image.png")
	}

	// Third file: text
	if result.Files[2].IsBinary {
		t.Error("Third file should not be binary")
	}
}

func TestParser_Parse_EmptyNewFile(t *testing.T) {
	diffContent := `diff --git a/empty.txt b/empty.txt
new file mode 100644
index 0000000..e69de29`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	file := result.Files[0]

	if !file.IsNew {
		t.Error("IsNew = false, want true")
	}

	// Empty file has no hunks
	if len(file.Hunks) != 0 {
		t.Errorf("Hunks count = %d, want 0 for empty file", len(file.Hunks))
	}
}
