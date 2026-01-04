package diff

// Diff represents a parsed git diff
type Diff struct {
	Files []FileDiff
}

// FileDiff represents changes to a single file
type FileDiff struct {
	OldPath string
	NewPath string
	Hunks   []Hunk
}

// Hunk represents a section of changes within a file
type Hunk struct {
	OldStart int
	OldCount int
	NewStart int
	NewCount int
	Lines    []Line
}

// Line represents a single line in a diff
type Line struct {
	Type    LineType
	Content string
	Number  int // Line number in the new file
}

// LineType represents the type of diff line
type LineType string

const (
	LineAdded   LineType = "added"
	LineRemoved LineType = "removed"
	LineContext LineType = "context"
)
