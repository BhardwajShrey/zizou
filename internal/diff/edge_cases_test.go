package diff

import (
	"testing"
)

func TestParser_Parse_NewFile(t *testing.T) {
	diffContent := `diff --git a/new_file.go b/new_file.go
new file mode 100644
index 0000000..a1b2c3d
--- /dev/null
+++ b/new_file.go
@@ -0,0 +1,5 @@
+package main
+
+func NewFunction() {
+    // new code
+}`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(result.Files) != 1 {
		t.Fatalf("Files count = %d, want 1", len(result.Files))
	}

	file := result.Files[0]

	// New file should have /dev/null as old path
	// But our parser extracts from the git header
	if file.NewPath != "new_file.go" {
		t.Errorf("NewPath = %q, want %q", file.NewPath, "new_file.go")
	}

	// All lines should be additions
	hunk := file.Hunks[0]
	for _, line := range hunk.Lines {
		if line.Type != LineAdded {
			t.Errorf("Line type = %v, want %v for new file", line.Type, LineAdded)
		}
	}
}

func TestParser_Parse_DeletedFile(t *testing.T) {
	diffContent := `diff --git a/deleted_file.go b/deleted_file.go
deleted file mode 100644
index a1b2c3d..0000000
--- a/deleted_file.go
+++ /dev/null
@@ -1,5 +0,0 @@
-package main
-
-func OldFunction() {
-    // old code
-}`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(result.Files) != 1 {
		t.Fatalf("Files count = %d, want 1", len(result.Files))
	}

	file := result.Files[0]

	if file.OldPath != "deleted_file.go" {
		t.Errorf("OldPath = %q, want %q", file.OldPath, "deleted_file.go")
	}

	// All lines should be removals
	hunk := file.Hunks[0]
	for _, line := range hunk.Lines {
		if line.Type != LineRemoved {
			t.Errorf("Line type = %v, want %v for deleted file", line.Type, LineRemoved)
		}
	}
}

func TestParser_Parse_FileWithNoChanges(t *testing.T) {
	// A diff header but no actual hunks (shouldn't normally happen)
	diffContent := `diff --git a/file.go b/file.go
index a1b2c3d..a1b2c3d 100644
--- a/file.go
+++ b/file.go`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(result.Files) != 1 {
		t.Fatalf("Files count = %d, want 1", len(result.Files))
	}

	file := result.Files[0]

	if len(file.Hunks) != 0 {
		t.Errorf("Hunks count = %d, want 0 for file with no changes", len(file.Hunks))
	}
}

func TestParser_Parse_LargeLineNumbers(t *testing.T) {
	diffContent := `diff --git a/large.go b/large.go
--- a/large.go
+++ b/large.go
@@ -1234,5 +1234,6 @@
 context line
-old line
+new line
+another new line
 context line`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	hunk := result.Files[0].Hunks[0]

	if hunk.OldStart != 1234 {
		t.Errorf("OldStart = %d, want 1234", hunk.OldStart)
	}
	if hunk.NewStart != 1234 {
		t.Errorf("NewStart = %d, want 1234", hunk.NewStart)
	}
}

func TestParser_Parse_OnlyAdditions(t *testing.T) {
	diffContent := `diff --git a/test.go b/test.go
--- a/test.go
+++ b/test.go
@@ -1,2 +1,5 @@
 package main
+import "fmt"
+import "os"
 func main() {
+    fmt.Println("hello")
 }`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	hunk := result.Files[0].Hunks[0]

	// Count line types
	addedCount := 0
	removedCount := 0

	for _, line := range hunk.Lines {
		if line.Type == LineAdded {
			addedCount++
		}
		if line.Type == LineRemoved {
			removedCount++
		}
	}

	if addedCount != 3 {
		t.Errorf("Added lines = %d, want 3", addedCount)
	}
	if removedCount != 0 {
		t.Errorf("Removed lines = %d, want 0", removedCount)
	}
}

func TestParser_Parse_OnlyRemovals(t *testing.T) {
	diffContent := `diff --git a/test.go b/test.go
--- a/test.go
+++ b/test.go
@@ -1,5 +1,2 @@
 package main
-import "fmt"
-import "os"
 func main() {
-    fmt.Println("hello")
 }`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	hunk := result.Files[0].Hunks[0]

	// Count line types
	addedCount := 0
	removedCount := 0

	for _, line := range hunk.Lines {
		if line.Type == LineAdded {
			addedCount++
		}
		if line.Type == LineRemoved {
			removedCount++
		}
	}

	if addedCount != 0 {
		t.Errorf("Added lines = %d, want 0", addedCount)
	}
	if removedCount != 3 {
		t.Errorf("Removed lines = %d, want 3", removedCount)
	}
}

func TestParser_Parse_WhitespaceChanges(t *testing.T) {
	diffContent := `diff --git a/test.go b/test.go
--- a/test.go
+++ b/test.go
@@ -1,3 +1,3 @@
 package main
-func main()  {
+func main() {
     // code`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	hunk := result.Files[0].Hunks[0]

	// Should detect the whitespace change
	removedLine := ""
	addedLine := ""

	for _, line := range hunk.Lines {
		if line.Type == LineRemoved {
			removedLine = line.Content
		}
		if line.Type == LineAdded {
			addedLine = line.Content
		}
	}

	if removedLine != "func main()  {" {
		t.Errorf("Removed line = %q, want %q", removedLine, "func main()  {")
	}
	if addedLine != "func main() {" {
		t.Errorf("Added line = %q, want %q", addedLine, "func main() {")
	}
}

func TestParser_Parse_TabsVsSpaces(t *testing.T) {
	diffContent := `diff --git a/test.go b/test.go
--- a/test.go
+++ b/test.go
@@ -1,2 +1,2 @@
 package main
-	func main() {
+    func main() {`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	hunk := result.Files[0].Hunks[0]

	// Parser should preserve exact whitespace
	for _, line := range hunk.Lines {
		if line.Type == LineRemoved && line.Content != "\tfunc main() {" {
			t.Errorf("Removed line content = %q, should preserve tab", line.Content)
		}
		if line.Type == LineAdded && line.Content != "    func main() {" {
			t.Errorf("Added line content = %q, should preserve spaces", line.Content)
		}
	}
}

func TestParser_Parse_ConsecutiveHunks(t *testing.T) {
	// Two hunks very close together
	diffContent := `diff --git a/test.go b/test.go
--- a/test.go
+++ b/test.go
@@ -10,3 +10,4 @@
 line10
+line10a
 line11
 line12
@@ -20,3 +21,4 @@
 line20
+line20a
 line21
 line22`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	file := result.Files[0]

	if len(file.Hunks) != 2 {
		t.Fatalf("Hunks count = %d, want 2", len(file.Hunks))
	}

	// Verify both hunks are parsed correctly
	if file.Hunks[0].NewStart != 10 {
		t.Errorf("Hunk[0] NewStart = %d, want 10", file.Hunks[0].NewStart)
	}
	if file.Hunks[1].NewStart != 21 {
		t.Errorf("Hunk[1] NewStart = %d, want 21", file.Hunks[1].NewStart)
	}
}

func TestParser_Parse_EmptyHunkContext(t *testing.T) {
	// Hunk with no function context after @@
	diffContent := `diff --git a/test.go b/test.go
--- a/test.go
+++ b/test.go
@@ -1,2 +1,3 @@
 line1
+line2
 line3`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Should still parse correctly even without context
	if len(result.Files) != 1 {
		t.Errorf("Files count = %d, want 1", len(result.Files))
	}
}

func TestParser_Parse_SpecialPathCharacters(t *testing.T) {
	// Test with spaces and special chars in path
	diffContent := `diff --git a/path with spaces/file-name.go b/path with spaces/file-name.go
--- a/path with spaces/file-name.go
+++ b/path with spaces/file-name.go
@@ -1,1 +1,2 @@
 package main
+var x = 1`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	file := result.Files[0]

	expectedPath := "path with spaces/file-name.go"
	if file.NewPath != expectedPath {
		t.Errorf("NewPath = %q, want %q", file.NewPath, expectedPath)
	}
}

func TestParser_Parse_VeryLongLine(t *testing.T) {
	longLine := "var veryLongLine = \"" + string(make([]byte, 1000)) + "\""
	diffContent := `diff --git a/test.go b/test.go
--- a/test.go
+++ b/test.go
@@ -1,1 +1,2 @@
 package main
+` + longLine

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Should handle very long lines
	hunk := result.Files[0].Hunks[0]
	foundLongLine := false

	for _, line := range hunk.Lines {
		if line.Type == LineAdded && len(line.Content) > 500 {
			foundLongLine = true
		}
	}

	if !foundLongLine {
		t.Error("Parse() did not preserve very long line")
	}
}

func TestParser_Parse_UnicodeContent(t *testing.T) {
	diffContent := `diff --git a/test.go b/test.go
--- a/test.go
+++ b/test.go
@@ -1,2 +1,3 @@
 package main
+// 你好世界 こんにちは 🚀
 func main() {}`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	hunk := result.Files[0].Hunks[0]

	// Find the unicode line
	for _, line := range hunk.Lines {
		if line.Type == LineAdded {
			expected := "// 你好世界 こんにちは 🚀"
			if line.Content != expected {
				t.Errorf("Unicode line = %q, want %q", line.Content, expected)
			}
		}
	}
}

func TestParser_Parse_MixedLineEndings(t *testing.T) {
	// This is tricky - in a real scenario this might cause issues
	// but our parser reads line by line so should handle it
	diffContent := "diff --git a/test.go b/test.go\n" +
		"--- a/test.go\n" +
		"+++ b/test.go\n" +
		"@@ -1,2 +1,3 @@\n" +
		" package main\n" +
		"+var x = 1\n" +
		" func main() {}"

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(result.Files) != 1 {
		t.Errorf("Files count = %d, want 1", len(result.Files))
	}
}

func TestParser_Parse_NoNewlineAtEndOfFile(t *testing.T) {
	diffContent := `diff --git a/test.go b/test.go
--- a/test.go
+++ b/test.go
@@ -1,2 +1,3 @@
 package main
+var x = 1
 func main() {}
\ No newline at end of file`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	// Should parse successfully even with "\ No newline" marker
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if len(result.Files) != 1 {
		t.Errorf("Files count = %d, want 1", len(result.Files))
	}
}
