package diff

import (
	"testing"
)

func TestDiff_GetAddedLines(t *testing.T) {
	diffContent := `diff --git a/test.go b/test.go
--- a/test.go
+++ b/test.go
@@ -1,3 +1,5 @@
 package main
+import "fmt"
 func main() {
+    fmt.Println("hello")
 }`

	parser := NewParser()
	d, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	added := d.GetAddedLines()

	if len(added) != 2 {
		t.Fatalf("GetAddedLines() count = %d, want 2", len(added))
	}

	// Check first added line
	if added[0].Content != `import "fmt"` {
		t.Errorf("Added line[0] content = %q, want %q", added[0].Content, `import "fmt"`)
	}
	if added[0].File != "test.go" {
		t.Errorf("Added line[0] file = %q, want %q", added[0].File, "test.go")
	}
	if added[0].Type != LineAdded {
		t.Errorf("Added line[0] type = %v, want %v", added[0].Type, LineAdded)
	}

	// Check second added line
	if added[1].Content != `    fmt.Println("hello")` {
		t.Errorf("Added line[1] content = %q, want %q", added[1].Content, `    fmt.Println("hello")`)
	}
}

func TestDiff_GetRemovedLines(t *testing.T) {
	diffContent := `diff --git a/test.go b/test.go
--- a/test.go
+++ b/test.go
@@ -1,5 +1,3 @@
 package main
-import "old"
 func main() {
-    old code
 }`

	parser := NewParser()
	d, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	removed := d.GetRemovedLines()

	if len(removed) != 2 {
		t.Fatalf("GetRemovedLines() count = %d, want 2", len(removed))
	}

	// Check first removed line
	if removed[0].Content != `import "old"` {
		t.Errorf("Removed line[0] content = %q, want %q", removed[0].Content, `import "old"`)
	}
	if removed[0].Type != LineRemoved {
		t.Errorf("Removed line[0] type = %v, want %v", removed[0].Type, LineRemoved)
	}

	// Check second removed line
	if removed[1].Content != `    old code` {
		t.Errorf("Removed line[1] content = %q, want %q", removed[1].Content, `    old code`)
	}
}

func TestDiff_GetModifiedFiles(t *testing.T) {
	diffContent := `diff --git a/file1.go b/file1.go
--- a/file1.go
+++ b/file1.go
@@ -1,1 +1,2 @@
 package main
+var x = 1
diff --git a/file2.go b/file2.go
--- a/file2.go
+++ b/file2.go
@@ -1,1 +1,2 @@
 package main
+var y = 2
diff --git a/file3.go b/file3.go
--- a/file3.go
+++ b/file3.go
@@ -1,1 +1,2 @@
 package main
+var z = 3`

	parser := NewParser()
	d, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	files := d.GetModifiedFiles()

	if len(files) != 3 {
		t.Fatalf("GetModifiedFiles() count = %d, want 3", len(files))
	}

	expected := []string{"file1.go", "file2.go", "file3.go"}
	for i, want := range expected {
		if files[i] != want {
			t.Errorf("GetModifiedFiles()[%d] = %q, want %q", i, files[i], want)
		}
	}
}

func TestDiff_Stats(t *testing.T) {
	diffContent := `diff --git a/file1.go b/file1.go
--- a/file1.go
+++ b/file1.go
@@ -1,3 +1,4 @@
 package main
-var old = 1
+var new = 2
+var extra = 3
diff --git a/file2.go b/file2.go
--- a/file2.go
+++ b/file2.go
@@ -1,2 +1,1 @@
 package main
-var removed = 1`

	parser := NewParser()
	d, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	stats := d.Stats()

	if stats.Files != 2 {
		t.Errorf("Stats() Files = %d, want 2", stats.Files)
	}

	// File1: +2 lines (var new, var extra)
	// File2: no additions
	// Total: +2
	if stats.LinesAdded != 2 {
		t.Errorf("Stats() LinesAdded = %d, want 2", stats.LinesAdded)
	}

	// File1: -1 line (var old)
	// File2: -1 line (var removed)
	// Total: -2
	if stats.LinesRemoved != 2 {
		t.Errorf("Stats() LinesRemoved = %d, want 2", stats.LinesRemoved)
	}
}

func TestDiff_Stats_EmptyDiff(t *testing.T) {
	d := &Diff{Files: []FileDiff{}}
	stats := d.Stats()

	if stats.Files != 0 {
		t.Errorf("Stats() Files = %d, want 0", stats.Files)
	}
	if stats.LinesAdded != 0 {
		t.Errorf("Stats() LinesAdded = %d, want 0", stats.LinesAdded)
	}
	if stats.LinesRemoved != 0 {
		t.Errorf("Stats() LinesRemoved = %d, want 0", stats.LinesRemoved)
	}
}

func TestHunk_GetContextAroundLine(t *testing.T) {
	hunk := &Hunk{
		Lines: []Line{
			{Type: LineContext, Content: "line1"},
			{Type: LineContext, Content: "line2"},
			{Type: LineContext, Content: "line3"},
			{Type: LineAdded, Content: "line4"},   // Index 3
			{Type: LineContext, Content: "line5"},
			{Type: LineContext, Content: "line6"},
			{Type: LineContext, Content: "line7"},
		},
	}

	tests := []struct {
		name        string
		lineIndex   int
		contextSize int
		wantStart   int
		wantEnd     int
	}{
		{
			name:        "middle line with context 1",
			lineIndex:   3,
			contextSize: 1,
			wantStart:   2,
			wantEnd:     5,
		},
		{
			name:        "middle line with context 2",
			lineIndex:   3,
			contextSize: 2,
			wantStart:   1,
			wantEnd:     6,
		},
		{
			name:        "first line with context 2",
			lineIndex:   0,
			contextSize: 2,
			wantStart:   0,
			wantEnd:     3,
		},
		{
			name:        "last line with context 2",
			lineIndex:   6,
			contextSize: 2,
			wantStart:   4,
			wantEnd:     7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hunk.GetContextAroundLine(tt.lineIndex, tt.contextSize)

			if len(result) != tt.wantEnd-tt.wantStart {
				t.Errorf("GetContextAroundLine() length = %d, want %d",
					len(result), tt.wantEnd-tt.wantStart)
			}

			// Check first and last lines match expected range
			if result[0].Content != hunk.Lines[tt.wantStart].Content {
				t.Errorf("GetContextAroundLine() first line = %q, want %q",
					result[0].Content, hunk.Lines[tt.wantStart].Content)
			}

			lastIdx := len(result) - 1
			if result[lastIdx].Content != hunk.Lines[tt.wantEnd-1].Content {
				t.Errorf("GetContextAroundLine() last line = %q, want %q",
					result[lastIdx].Content, hunk.Lines[tt.wantEnd-1].Content)
			}
		})
	}
}

func TestHunk_HasChanges(t *testing.T) {
	tests := []struct {
		name string
		hunk *Hunk
		want bool
	}{
		{
			name: "has added lines",
			hunk: &Hunk{
				Lines: []Line{
					{Type: LineContext, Content: "context"},
					{Type: LineAdded, Content: "added"},
				},
			},
			want: true,
		},
		{
			name: "has removed lines",
			hunk: &Hunk{
				Lines: []Line{
					{Type: LineContext, Content: "context"},
					{Type: LineRemoved, Content: "removed"},
				},
			},
			want: true,
		},
		{
			name: "has both added and removed",
			hunk: &Hunk{
				Lines: []Line{
					{Type: LineRemoved, Content: "removed"},
					{Type: LineAdded, Content: "added"},
				},
			},
			want: true,
		},
		{
			name: "only context lines",
			hunk: &Hunk{
				Lines: []Line{
					{Type: LineContext, Content: "context1"},
					{Type: LineContext, Content: "context2"},
					{Type: LineContext, Content: "context3"},
				},
			},
			want: false,
		},
		{
			name: "empty hunk",
			hunk: &Hunk{
				Lines: []Line{},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.hunk.HasChanges()
			if got != tt.want {
				t.Errorf("HasChanges() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiff_GetAddedLines_MultipleFiles(t *testing.T) {
	diffContent := `diff --git a/file1.go b/file1.go
--- a/file1.go
+++ b/file1.go
@@ -1,1 +1,2 @@
 package main
+var x = 1
diff --git a/file2.go b/file2.go
--- a/file2.go
+++ b/file2.go
@@ -1,1 +1,2 @@
 package main
+var y = 2`

	parser := NewParser()
	d, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	added := d.GetAddedLines()

	if len(added) != 2 {
		t.Fatalf("GetAddedLines() count = %d, want 2", len(added))
	}

	// Should have additions from both files (check both are present)
	files := make(map[string]bool)
	for _, line := range added {
		files[line.File] = true
	}

	// Debug: show what files were found
	if len(files) != 2 {
		t.Logf("Found %d unique files, expected 2", len(files))
		for file := range files {
			t.Logf("  File: %q", file)
		}
	}

	if !files["file1.go"] {
		t.Error("Expected addition from file1.go not found")
	}
	if !files["file2.go"] {
		t.Error("Expected addition from file2.go not found")
	}
}

func TestDiff_GetRemovedLines_NoRemovals(t *testing.T) {
	diffContent := `diff --git a/test.go b/test.go
--- a/test.go
+++ b/test.go
@@ -1,1 +1,2 @@
 package main
+var x = 1`

	parser := NewParser()
	d, err := parser.Parse(diffContent)

	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	removed := d.GetRemovedLines()

	if len(removed) != 0 {
		t.Errorf("GetRemovedLines() count = %d, want 0", len(removed))
	}
}

func TestLineInfo_Fields(t *testing.T) {
	info := LineInfo{
		File:    "test.go",
		Line:    42,
		Content: "test content",
		Type:    LineAdded,
	}

	if info.File != "test.go" {
		t.Errorf("LineInfo.File = %q, want %q", info.File, "test.go")
	}
	if info.Line != 42 {
		t.Errorf("LineInfo.Line = %d, want %d", info.Line, 42)
	}
	if info.Content != "test content" {
		t.Errorf("LineInfo.Content = %q, want %q", info.Content, "test content")
	}
	if info.Type != LineAdded {
		t.Errorf("LineInfo.Type = %v, want %v", info.Type, LineAdded)
	}
}

func TestDiffStats_Fields(t *testing.T) {
	stats := DiffStats{
		Files:        3,
		LinesAdded:   10,
		LinesRemoved: 5,
	}

	if stats.Files != 3 {
		t.Errorf("DiffStats.Files = %d, want %d", stats.Files, 3)
	}
	if stats.LinesAdded != 10 {
		t.Errorf("DiffStats.LinesAdded = %d, want %d", stats.LinesAdded, 10)
	}
	if stats.LinesRemoved != 5 {
		t.Errorf("DiffStats.LinesRemoved = %d, want %d", stats.LinesRemoved, 5)
	}
}
