package diff

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	fileHeaderRegex  = regexp.MustCompile(`^diff --git a/(.*) b/(.*)$`)
	oldFileRegex     = regexp.MustCompile(`^--- a/(.*)$`)
	newFileRegex     = regexp.MustCompile(`^\+\+\+ b/(.*)$`)
	hunkHeaderRegex  = regexp.MustCompile(`^@@ -(\d+),?(\d*) \+(\d+),?(\d*) @@`)
	binaryFileRegex  = regexp.MustCompile(`^Binary files (.*) and (.*) differ$`)
	newFileModeRegex = regexp.MustCompile(`^new file mode (\d+)$`)
	deletedFileModeRegex = regexp.MustCompile(`^deleted file mode (\d+)$`)
	oldModeRegex     = regexp.MustCompile(`^old mode (\d+)$`)
	newModeRegex     = regexp.MustCompile(`^new mode (\d+)$`)
	renameFromRegex  = regexp.MustCompile(`^rename from (.*)$`)
	renameToRegex    = regexp.MustCompile(`^rename to (.*)$`)
	copyFromRegex    = regexp.MustCompile(`^copy from (.*)$`)
	copyToRegex      = regexp.MustCompile(`^copy to (.*)$`)
	similarityRegex  = regexp.MustCompile(`^similarity index (\d+)%$`)
	indexRegex       = regexp.MustCompile(`^index ([0-9a-f]+)\.\.([0-9a-f]+)`)
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
		if matches := oldFileRegex.FindStringSubmatch(line); matches != nil {
			if currentFile != nil && matches[1] == "/dev/null" {
				currentFile.IsNew = true
			}
			continue
		}

		// Check for new file marker
		if matches := newFileRegex.FindStringSubmatch(line); matches != nil {
			if currentFile != nil && matches[1] == "/dev/null" {
				currentFile.IsDeleted = true
			}
			continue
		}

		// Check for binary file
		if currentFile != nil && binaryFileRegex.MatchString(line) {
			currentFile.IsBinary = true
			continue
		}

		// Check for new file mode
		if currentFile != nil {
			if matches := newFileModeRegex.FindStringSubmatch(line); matches != nil {
				currentFile.IsNew = true
				currentFile.NewMode = matches[1]
				continue
			}
		}

		// Check for deleted file mode
		if currentFile != nil {
			if matches := deletedFileModeRegex.FindStringSubmatch(line); matches != nil {
				currentFile.IsDeleted = true
				currentFile.OldMode = matches[1]
				continue
			}
		}

		// Check for old mode
		if currentFile != nil {
			if matches := oldModeRegex.FindStringSubmatch(line); matches != nil {
				currentFile.OldMode = matches[1]
				continue
			}
		}

		// Check for new mode
		if currentFile != nil {
			if matches := newModeRegex.FindStringSubmatch(line); matches != nil {
				currentFile.NewMode = matches[1]
				continue
			}
		}

		// Check for rename from
		if currentFile != nil {
			if matches := renameFromRegex.FindStringSubmatch(line); matches != nil {
				currentFile.IsRenamed = true
				currentFile.OldPath = matches[1]
				continue
			}
		}

		// Check for rename to
		if currentFile != nil {
			if matches := renameToRegex.FindStringSubmatch(line); matches != nil {
				currentFile.IsRenamed = true
				currentFile.NewPath = matches[1]
				continue
			}
		}

		// Check for copy from
		if currentFile != nil {
			if matches := copyFromRegex.FindStringSubmatch(line); matches != nil {
				currentFile.IsCopied = true
				currentFile.OldPath = matches[1]
				continue
			}
		}

		// Check for copy to
		if currentFile != nil {
			if matches := copyToRegex.FindStringSubmatch(line); matches != nil {
				currentFile.IsCopied = true
				currentFile.NewPath = matches[1]
				continue
			}
		}

		// Check for similarity index
		if currentFile != nil {
			if matches := similarityRegex.FindStringSubmatch(line); matches != nil {
				similarity, _ := strconv.Atoi(matches[1])
				currentFile.Similarity = similarity
				continue
			}
		}

		// Skip index lines
		if indexRegex.MatchString(line) {
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
