package diff

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	fileHeaderRegex = regexp.MustCompile(`^diff --git a/(.*) b/(.*)$`)
	oldFileRegex    = regexp.MustCompile(`^--- a/(.*)$`)
	newFileRegex    = regexp.MustCompile(`^\+\+\+ b/(.*)$`)
	hunkHeaderRegex = regexp.MustCompile(`^@@ -(\d+),?(\d*) \+(\d+),?(\d*) @@`)
)

// Parser parses git diffs
type Parser struct{}

// NewParser creates a new diff parser
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses a git diff string into a Diff structure
func (p *Parser) Parse(content string) (*Diff, error) {
	if content == "" {
		return nil, fmt.Errorf("empty diff content")
	}

	scanner := bufio.NewScanner(strings.NewReader(content))
	diff := &Diff{
		Files: []FileDiff{},
	}

	var currentFile *FileDiff
	var currentHunk *Hunk
	lineNumber := 0

	for scanner.Scan() {
		line := scanner.Text()

		// Check for file header
		if matches := fileHeaderRegex.FindStringSubmatch(line); matches != nil {
			// Add current hunk to current file before switching files
			if currentHunk != nil && currentFile != nil {
				currentFile.Hunks = append(currentFile.Hunks, *currentHunk)
			}
			// Add current file to diff before switching to new file
			if currentFile != nil {
				diff.Files = append(diff.Files, *currentFile)
			}
			currentFile = &FileDiff{
				OldPath: matches[1],
				NewPath: matches[2],
				Hunks:   []Hunk{},
			}
			currentHunk = nil
			continue
		}

		// Check for old file marker
		if oldFileRegex.MatchString(line) {
			continue
		}

		// Check for new file marker
		if newFileRegex.MatchString(line) {
			continue
		}

		// Check for hunk header
		if matches := hunkHeaderRegex.FindStringSubmatch(line); matches != nil {
			if currentHunk != nil && currentFile != nil {
				currentFile.Hunks = append(currentFile.Hunks, *currentHunk)
			}

			oldStart, _ := strconv.Atoi(matches[1])
			oldCount := 0
			if matches[2] != "" {
				oldCount, _ = strconv.Atoi(matches[2])
			} else {
				oldCount = 1
			}

			newStart, _ := strconv.Atoi(matches[3])
			newCount := 0
			if matches[4] != "" {
				newCount, _ = strconv.Atoi(matches[4])
			} else {
				newCount = 1
			}

			currentHunk = &Hunk{
				OldStart: oldStart,
				OldCount: oldCount,
				NewStart: newStart,
				NewCount: newCount,
				Lines:    []Line{},
			}
			lineNumber = newStart
			continue
		}

		// Parse diff lines
		if currentHunk != nil {
			if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
				currentHunk.Lines = append(currentHunk.Lines, Line{
					Type:    LineAdded,
					Content: strings.TrimPrefix(line, "+"),
					Number:  lineNumber,
				})
				lineNumber++
			} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
				currentHunk.Lines = append(currentHunk.Lines, Line{
					Type:    LineRemoved,
					Content: strings.TrimPrefix(line, "-"),
					Number:  lineNumber - 1,
				})
			} else if strings.HasPrefix(line, " ") {
				currentHunk.Lines = append(currentHunk.Lines, Line{
					Type:    LineContext,
					Content: strings.TrimPrefix(line, " "),
					Number:  lineNumber,
				})
				lineNumber++
			}
		}
	}

	// Add last hunk and file
	if currentHunk != nil && currentFile != nil {
		currentFile.Hunks = append(currentFile.Hunks, *currentHunk)
	}
	if currentFile != nil {
		diff.Files = append(diff.Files, *currentFile)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning diff: %w", err)
	}

	return diff, nil
}
