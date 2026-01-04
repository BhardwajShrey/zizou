package diff

import (
	"testing"
)

func TestParser_Parse_BasicDiff(t *testing.T) {
	diffContent := `diff --git a/app/main.go b/app/main.go
index 30dcdee..2cb2e68 100644
--- a/app/main.go
+++ b/app/main.go
@@ -23,7 +23,8 @@ func main() {

     image := os.Args[2]

-    EnterNewJail(os.Args[3], image)
+    jailPath := EnterNewJail(os.Args[3], image)
+    defer os.Remove(jailPath)

     cmd := exec.Command(command, args...)`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v, want nil", err)
	}

	if result == nil {
		t.Fatal("Parse() returned nil result")
	}

	// Check file count
	if len(result.Files) != 1 {
		t.Errorf("Parse() files count = %d, want 1", len(result.Files))
	}

	// Check file path
	file := result.Files[0]
	if file.NewPath != "app/main.go" {
		t.Errorf("Parse() file path = %q, want %q", file.NewPath, "app/main.go")
	}

	// Check hunk count
	if len(file.Hunks) != 1 {
		t.Errorf("Parse() hunks count = %d, want 1", len(file.Hunks))
	}

	hunk := file.Hunks[0]

	// Check hunk line ranges
	if hunk.OldStart != 23 {
		t.Errorf("Parse() hunk.OldStart = %d, want 23", hunk.OldStart)
	}
	if hunk.NewStart != 23 {
		t.Errorf("Parse() hunk.NewStart = %d, want 23", hunk.NewStart)
	}
	if hunk.OldCount != 7 {
		t.Errorf("Parse() hunk.OldCount = %d, want 7", hunk.OldCount)
	}
	if hunk.NewCount != 8 {
		t.Errorf("Parse() hunk.NewCount = %d, want 8", hunk.NewCount)
	}

	// Check line types
	addedCount := 0
	removedCount := 0
	contextCount := 0

	for _, line := range hunk.Lines {
		switch line.Type {
		case LineAdded:
			addedCount++
		case LineRemoved:
			removedCount++
		case LineContext:
			contextCount++
		}
	}

	if addedCount != 2 {
		t.Errorf("Parse() added lines = %d, want 2", addedCount)
	}
	if removedCount != 1 {
		t.Errorf("Parse() removed lines = %d, want 1", removedCount)
	}
	// Context lines: empty, "image := ...", empty, cmd line = varies based on trailing content
	// The actual count depends on whether trailing context is included
	if contextCount < 2 {
		t.Errorf("Parse() context lines = %d, want at least 2", contextCount)
	}
}

func TestParser_Parse_MultipleHunks(t *testing.T) {
	diffContent := `diff --git a/app/main.go b/app/main.go
--- a/app/main.go
+++ b/app/main.go
@@ -23,7 +23,8 @@ func main() {
     image := os.Args[2]
-    EnterNewJail(os.Args[3], image)
+    jailPath := EnterNewJail(os.Args[3], image)
+    defer os.Remove(jailPath)
     cmd := exec.Command(command, args...)
@@ -40,3 +41,4 @@ func EnterNewJail(filepath string, image string) {
         throwerror.ThrowError(err, "Error")
     }
+    return tempDirPath
 }`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v, want nil", err)
	}

	file := result.Files[0]
	if len(file.Hunks) != 2 {
		t.Errorf("Parse() hunks count = %d, want 2", len(file.Hunks))
	}

	// Check first hunk
	hunk1 := file.Hunks[0]
	if hunk1.OldStart != 23 || hunk1.NewStart != 23 {
		t.Errorf("Parse() hunk1 starts = %d,%d, want 23,23", hunk1.OldStart, hunk1.NewStart)
	}

	// Check second hunk
	hunk2 := file.Hunks[1]
	if hunk2.OldStart != 40 || hunk2.NewStart != 41 {
		t.Errorf("Parse() hunk2 starts = %d,%d, want 40,41", hunk2.OldStart, hunk2.NewStart)
	}
}

func TestParser_Parse_MultipleFiles(t *testing.T) {
	diffContent := `diff --git a/file1.go b/file1.go
--- a/file1.go
+++ b/file1.go
@@ -1,3 +1,4 @@
 package main
+import "fmt"
 func main() {}
diff --git a/file2.go b/file2.go
--- a/file2.go
+++ b/file2.go
@@ -1,2 +1,3 @@
 package main
+var x = 1`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v, want nil", err)
	}

	if len(result.Files) != 2 {
		t.Fatalf("Parse() files count = %d, want 2", len(result.Files))
	}

	// Check first file
	if result.Files[0].NewPath != "file1.go" {
		t.Errorf("Parse() file[0] path = %q, want %q", result.Files[0].NewPath, "file1.go")
	}

	// Check second file
	if result.Files[1].NewPath != "file2.go" {
		t.Errorf("Parse() file[1] path = %q, want %q", result.Files[1].NewPath, "file2.go")
	}
}

func TestParser_Parse_EmptyDiff(t *testing.T) {
	parser := NewParser()
	_, err := parser.Parse("")

	if err == nil {
		t.Error("Parse() with empty string should return error")
	}
}

func TestParser_Parse_LineTypes(t *testing.T) {
	diffContent := `diff --git a/test.go b/test.go
--- a/test.go
+++ b/test.go
@@ -1,4 +1,5 @@
 package main
-var old = 1
+var new = 2
+var added = 3
 func test() {}`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	hunk := result.Files[0].Hunks[0]

	// Verify we can identify each line type
	foundContext := false
	foundRemoved := false
	foundAdded := false

	for _, line := range hunk.Lines {
		switch line.Type {
		case LineContext:
			foundContext = true
		case LineRemoved:
			foundRemoved = true
			if line.Content != "var old = 1" {
				t.Errorf("Removed line content = %q, want %q", line.Content, "var old = 1")
			}
		case LineAdded:
			foundAdded = true
		}
	}

	if !foundContext {
		t.Error("Parse() did not find any context lines")
	}
	if !foundRemoved {
		t.Error("Parse() did not find any removed lines")
	}
	if !foundAdded {
		t.Error("Parse() did not find any added lines")
	}
}

func TestParser_Parse_HunkHeaderFormats(t *testing.T) {
	tests := []struct {
		name        string
		hunkHeader  string
		wantOldStart int
		wantOldCount int
		wantNewStart int
		wantNewCount int
	}{
		{
			name:        "standard format",
			hunkHeader:  "@@ -23,7 +23,8 @@",
			wantOldStart: 23,
			wantOldCount: 7,
			wantNewStart: 23,
			wantNewCount: 8,
		},
		{
			name:        "single line old",
			hunkHeader:  "@@ -23 +23,2 @@",
			wantOldStart: 23,
			wantOldCount: 1,
			wantNewStart: 23,
			wantNewCount: 2,
		},
		{
			name:        "single line new",
			hunkHeader:  "@@ -23,2 +23 @@",
			wantOldStart: 23,
			wantOldCount: 2,
			wantNewStart: 23,
			wantNewCount: 1,
		},
		{
			name:        "single line both",
			hunkHeader:  "@@ -23 +23 @@",
			wantOldStart: 23,
			wantOldCount: 1,
			wantNewStart: 23,
			wantNewCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diffContent := `diff --git a/test.go b/test.go
--- a/test.go
+++ b/test.go
` + tt.hunkHeader + `
 context line`

			parser := NewParser()
			result, err := parser.Parse(diffContent)

			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			hunk := result.Files[0].Hunks[0]

			if hunk.OldStart != tt.wantOldStart {
				t.Errorf("OldStart = %d, want %d", hunk.OldStart, tt.wantOldStart)
			}
			if hunk.OldCount != tt.wantOldCount {
				t.Errorf("OldCount = %d, want %d", hunk.OldCount, tt.wantOldCount)
			}
			if hunk.NewStart != tt.wantNewStart {
				t.Errorf("NewStart = %d, want %d", hunk.NewStart, tt.wantNewStart)
			}
			if hunk.NewCount != tt.wantNewCount {
				t.Errorf("NewCount = %d, want %d", hunk.NewCount, tt.wantNewCount)
			}
		})
	}
}

func TestParser_Parse_SpecialCharactersInContent(t *testing.T) {
	diffContent := `diff --git a/test.go b/test.go
--- a/test.go
+++ b/test.go
@@ -1,2 +1,2 @@
-fmt.Println("Hello---World")
+fmt.Println("Hello+++World")`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	hunk := result.Files[0].Hunks[0]

	// Should have 1 removed and 1 added line
	removedCount := 0
	addedCount := 0

	for _, line := range hunk.Lines {
		if line.Type == LineRemoved {
			removedCount++
			if line.Content != `fmt.Println("Hello---World")` {
				t.Errorf("Removed line = %q, want %q", line.Content, `fmt.Println("Hello---World")`)
			}
		}
		if line.Type == LineAdded {
			addedCount++
			if line.Content != `fmt.Println("Hello+++World")` {
				t.Errorf("Added line = %q, want %q", line.Content, `fmt.Println("Hello+++World")`)
			}
		}
	}

	if removedCount != 1 || addedCount != 1 {
		t.Errorf("Line counts: removed=%d, added=%d, want 1, 1", removedCount, addedCount)
	}
}

func TestParser_Parse_EmptyLines(t *testing.T) {
	diffContent := `diff --git a/test.go b/test.go
--- a/test.go
+++ b/test.go
@@ -1,4 +1,5 @@
 line1

-line3
+line3modified
+
 line4`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	hunk := result.Files[0].Hunks[0]

	// Should handle empty lines correctly
	emptyContextLines := 0
	emptyAddedLines := 0

	for _, line := range hunk.Lines {
		if line.Content == "" {
			if line.Type == LineContext {
				emptyContextLines++
			} else if line.Type == LineAdded {
				emptyAddedLines++
			}
		}
	}

	// Should have at least one empty line (context or added)
	totalEmptyLines := emptyContextLines + emptyAddedLines
	if totalEmptyLines == 0 {
		t.Error("Parse() did not preserve any empty lines")
	}

	// Specifically, we added an empty line
	if emptyAddedLines == 0 {
		t.Error("Parse() did not detect the added empty line")
	}
}

func TestParser_Parse_RenamedFile(t *testing.T) {
	diffContent := `diff --git a/oldname.go b/newname.go
similarity index 95%
rename from oldname.go
rename to newname.go
--- a/oldname.go
+++ b/newname.go
@@ -1,2 +1,3 @@
 package main
+// renamed file`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	file := result.Files[0]

	if file.OldPath != "oldname.go" {
		t.Errorf("OldPath = %q, want %q", file.OldPath, "oldname.go")
	}

	if file.NewPath != "newname.go" {
		t.Errorf("NewPath = %q, want %q", file.NewPath, "newname.go")
	}
}

func TestParser_Parse_LineNumbers(t *testing.T) {
	diffContent := `diff --git a/test.go b/test.go
--- a/test.go
+++ b/test.go
@@ -10,5 +10,6 @@
 line10
 line11
-line12
+line12modified
+line13new
 line14`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	hunk := result.Files[0].Hunks[0]

	// Check that line numbers are tracked correctly
	for _, line := range hunk.Lines {
		if line.Type == LineAdded && line.Content == "line12modified" {
			if line.Number != 12 {
				t.Errorf("Line number for 'line12modified' = %d, want 12", line.Number)
			}
		}
		if line.Type == LineAdded && line.Content == "line13new" {
			if line.Number != 13 {
				t.Errorf("Line number for 'line13new' = %d, want 13", line.Number)
			}
		}
	}
}

func TestParser_NewParser(t *testing.T) {
	parser := NewParser()

	if parser == nil {
		t.Error("NewParser() returned nil")
	}
}

func TestParser_Parse_ComplexRealWorldDiff(t *testing.T) {
	// Test with the actual example from the user
	diffContent := `diff --git a/app/main.go b/app/main.go
index 30dcdee..2cb2e68 100644
--- a/app/main.go
+++ b/app/main.go
@@ -23,7 +23,8 @@ func main() {

     image := os.Args[2]

-    EnterNewJail(os.Args[3], image)
+    jailPath := EnterNewJail(os.Args[3], image)
+    defer os.Remove(jailPath)

     cmd := exec.Command(command, args...)
     cmd.Stderr = os.Stderr
@@ -37,14 +38,12 @@ func main() {
     os.Exit(cmd.ProcessState.ExitCode())
 }

-func EnterNewJail(filepath string, image string) {
+func EnterNewJail(filepath string, image string) string {
     tempDirPath, err := os.MkdirTemp("", "temp_folder_*")
     if err != nil {
         throwerror.ThrowError(err, "Unable to create temp directory")
     }

-    defer os.Remove(tempDirPath)
-
     err = os.Chmod(tempDirPath, 0777)
     if err != nil {
         throwerror.ThrowError(err, "Error modifying rwx on tempDirPath")
@@ -58,4 +57,6 @@ func EnterNewJail(filepath string, image string) {
     if err != nil {
         throwerror.ThrowError(err, fmt.Sprintf("Error in executing chroot on %s", tempDirPath))
     }
+
+    return tempDirPath
 }`

	parser := NewParser()
	result, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Should have 1 file
	if len(result.Files) != 1 {
		t.Fatalf("Files count = %d, want 1", len(result.Files))
	}

	file := result.Files[0]

	// Should have 3 hunks
	if len(file.Hunks) != 3 {
		t.Fatalf("Hunks count = %d, want 3", len(file.Hunks))
	}

	// Verify each hunk has the correct ranges
	expectedHunks := []struct {
		oldStart, oldCount, newStart, newCount int
	}{
		{23, 7, 23, 8},
		{37, 14, 38, 12},
		{58, 4, 57, 6},
	}

	for i, expected := range expectedHunks {
		hunk := file.Hunks[i]
		if hunk.OldStart != expected.oldStart {
			t.Errorf("Hunk %d: OldStart = %d, want %d", i, hunk.OldStart, expected.oldStart)
		}
		if hunk.OldCount != expected.oldCount {
			t.Errorf("Hunk %d: OldCount = %d, want %d", i, hunk.OldCount, expected.oldCount)
		}
		if hunk.NewStart != expected.newStart {
			t.Errorf("Hunk %d: NewStart = %d, want %d", i, hunk.NewStart, expected.newStart)
		}
		if hunk.NewCount != expected.newCount {
			t.Errorf("Hunk %d: NewCount = %d, want %d", i, hunk.NewCount, expected.newCount)
		}
	}
}
