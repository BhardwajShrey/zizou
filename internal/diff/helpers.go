package diff

// GetAddedLines returns all added lines from the diff
func (d *Diff) GetAddedLines() []LineInfo {
	var lines []LineInfo
	for _, file := range d.Files {
		for _, hunk := range file.Hunks {
			for _, line := range hunk.Lines {
				if line.Type == LineAdded {
					lines = append(lines, LineInfo{
						File:    file.NewPath,
						Line:    line.Number,
						Content: line.Content,
						Type:    line.Type,
					})
				}
			}
		}
	}
	return lines
}

// GetRemovedLines returns all removed lines from the diff
func (d *Diff) GetRemovedLines() []LineInfo {
	var lines []LineInfo
	for _, file := range d.Files {
		for _, hunk := range file.Hunks {
			for _, line := range hunk.Lines {
				if line.Type == LineRemoved {
					lines = append(lines, LineInfo{
						File:    file.OldPath,
						Line:    line.Number,
						Content: line.Content,
						Type:    line.Type,
					})
				}
			}
		}
	}
	return lines
}

// GetModifiedFiles returns a list of all modified file paths
func (d *Diff) GetModifiedFiles() []string {
	files := make([]string, 0, len(d.Files))
	for _, file := range d.Files {
		files = append(files, file.NewPath)
	}
	return files
}

// Stats returns statistics about the diff
func (d *Diff) Stats() DiffStats {
	stats := DiffStats{
		Files: len(d.Files),
	}

	for _, file := range d.Files {
		for _, hunk := range file.Hunks {
			for _, line := range hunk.Lines {
				switch line.Type {
				case LineAdded:
					stats.LinesAdded++
				case LineRemoved:
					stats.LinesRemoved++
				}
			}
		}
	}

	return stats
}

// GetContextAroundLine returns lines of context around a specific line
func (h *Hunk) GetContextAroundLine(lineIndex, contextSize int) []Line {
	start := max(0, lineIndex-contextSize)
	end := min(len(h.Lines), lineIndex+contextSize+1)
	return h.Lines[start:end]
}

// HasChanges returns true if the hunk contains any added or removed lines
func (h *Hunk) HasChanges() bool {
	for _, line := range h.Lines {
		if line.Type == LineAdded || line.Type == LineRemoved {
			return true
		}
	}
	return false
}

// LineInfo contains information about a specific line in the diff
type LineInfo struct {
	File    string
	Line    int
	Content string
	Type    LineType
}

// DiffStats contains statistics about a diff
type DiffStats struct {
	Files        int
	LinesAdded   int
	LinesRemoved int
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
