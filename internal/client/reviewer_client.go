package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/BhardwajShrey/zizou/internal/diff"
	"github.com/BhardwajShrey/zizou/internal/review"
)

// ReviewerClient performs code reviews using Claude
type ReviewerClient struct {
	client *EnhancedClient
}

// NewReviewerClient creates a new code reviewer client
func NewReviewerClient(config *Config) (*ReviewerClient, error) {
	client, err := NewEnhancedClient(config)
	if err != nil {
		return nil, err
	}

	return &ReviewerClient{
		client: client,
	}, nil
}

// ReviewDiff reviews a parsed diff and returns structured comments
func (rc *ReviewerClient) ReviewDiff(ctx context.Context, d *diff.Diff) (*review.Result, error) {
	// Build the review prompt
	prompt := rc.buildReviewPrompt(d)

	// Send to Claude
	response, err := rc.client.SendMessage(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to get review from Claude: %w", err)
	}

	// Parse the response
	result, err := rc.parseReviewResponse(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse review response: %w", err)
	}

	return result, nil
}

// buildReviewPrompt creates a detailed prompt for code review
func (rc *ReviewerClient) buildReviewPrompt(d *diff.Diff) string {
	var sb strings.Builder

	// Concise senior reviewer prompt
	sb.WriteString("You are a senior code reviewer. Review this diff focusing on impact-driven, actionable feedback. Be thorough but concise.\n\n")

	sb.WriteString("## Focus Areas\n")
	sb.WriteString("1. **Security** (critical): Injection, auth bypasses, data exposure, race conditions\n")
	sb.WriteString("2. **Bugs** (high): Logic errors, nil pointers, resource leaks, missing error handling\n")
	sb.WriteString("3. **Performance** (medium-high): Inefficient algorithms, unnecessary allocations, blocking operations\n")
	sb.WriteString("4. **Maintainability** (medium): Complex logic, poor naming, duplication\n\n")

	sb.WriteString("## Severity\n")
	sb.WriteString("**critical**: Security/crashes/data loss | **high**: Likely bugs/broken errors | **medium**: Potential bugs | **low**: Minor issues | **info**: Suggestions\n\n")

	sb.WriteString("## Language-Specific\n")
	sb.WriteString("**Go**: Error handling, goroutine leaks, races, defer in loops\n")
	sb.WriteString("**C++**: Memory safety, RAII, iterator invalidation, UB\n")
	sb.WriteString("**Python**: Mutable defaults, resource cleanup, exception handling\n")
	sb.WriteString("**JS**: Unhandled promises, null checks, event listener leaks\n\n")

	// Add statistics
	stats := d.Stats()
	sb.WriteString(fmt.Sprintf("## Diff Statistics\n"))
	sb.WriteString(fmt.Sprintf("- Files changed: %d\n", stats.Files))
	sb.WriteString(fmt.Sprintf("- Lines added: +%d\n", stats.LinesAdded))
	sb.WriteString(fmt.Sprintf("- Lines removed: -%d\n\n", stats.LinesRemoved))

	// Add file-specific context
	sb.WriteString("## Files Modified\n")
	for _, file := range d.Files {
		sb.WriteString(fmt.Sprintf("- `%s`", file.NewPath))

		// Add context about file operations
		if file.IsNew {
			sb.WriteString(" (new file)")
		} else if file.IsDeleted {
			sb.WriteString(" (deleted)")
		} else if file.IsRenamed {
			sb.WriteString(fmt.Sprintf(" (renamed from %s)", file.OldPath))
		} else if file.IsBinary {
			sb.WriteString(" (binary)")
		}

		if file.OldMode != file.NewMode && file.NewMode != "" {
			sb.WriteString(fmt.Sprintf(" [mode: %s → %s]", file.OldMode, file.NewMode))
		}

		sb.WriteString("\n")
	}
	sb.WriteString("\n")

	// Add the actual diff
	sb.WriteString("## Git Diff\n")
	sb.WriteString("```diff\n")
	sb.WriteString(rc.formatDiff(d))
	sb.WriteString("```\n\n")

	// Request structured JSON output
	sb.WriteString("## Required Output Format\n")
	sb.WriteString("Provide your review as a JSON object with the following structure:\n\n")
	sb.WriteString("```json\n")
	sb.WriteString("{\n")
	sb.WriteString("  \"comments\": [\n")
	sb.WriteString("    {\n")
	sb.WriteString("      \"file\": \"path/to/file.go\",\n")
	sb.WriteString("      \"line\": 42,\n")
	sb.WriteString("      \"severity\": \"critical|high|medium|low|info\",\n")
	sb.WriteString("      \"category\": \"security|performance|bug_risk|maintainability|style|best_practice|other\",\n")
	sb.WriteString("      \"message\": \"Clear description of the issue\",\n")
	sb.WriteString("      \"snippet\": \"relevant code snippet (optional)\"\n")
	sb.WriteString("    }\n")
	sb.WriteString("  ],\n")
	sb.WriteString("  \"summary\": \"Overall assessment of the changes\"\n")
	sb.WriteString("}\n")
	sb.WriteString("```\n\n")

	sb.WriteString("## Message Format\n")
	sb.WriteString("Structure: **Issue** [what's wrong] | **Impact** [why it matters] | **Fix** [specific solution with code]\n")
	sb.WriteString("Example: \"Issue: SQL injection via unsanitized input. Impact: Arbitrary query execution. Fix: Use db.Query(\\\"SELECT * FROM users WHERE id = ?\\\", userId)\"\n\n")

	sb.WriteString("## Rules\n")
	sb.WriteString("- Return ONLY valid JSON, no extra text\n")
	sb.WriteString("- Include line numbers and concrete fixes with code examples\n")
	sb.WriteString("- Focus on significant issues; skip style nitpicks and linter-caught items\n")
	sb.WriteString("- Empty comments array if code is good\n")
	sb.WriteString("- Be concise but complete - optimize for minimal tokens\n")

	return sb.String()
}

// formatDiff converts a parsed diff back to string format
func (rc *ReviewerClient) formatDiff(d *diff.Diff) string {
	var sb strings.Builder

	for _, file := range d.Files {
		// Skip binary files (they don't have meaningful content to review)
		if file.IsBinary {
			sb.WriteString(fmt.Sprintf("diff --git a/%s b/%s\n", file.OldPath, file.NewPath))
			sb.WriteString("Binary files differ\n\n")
			continue
		}

		sb.WriteString(fmt.Sprintf("diff --git a/%s b/%s\n", file.OldPath, file.NewPath))

		// Add metadata
		if file.IsNew {
			sb.WriteString(fmt.Sprintf("new file mode %s\n", file.NewMode))
		} else if file.IsDeleted {
			sb.WriteString(fmt.Sprintf("deleted file mode %s\n", file.OldMode))
		} else if file.IsRenamed {
			sb.WriteString(fmt.Sprintf("rename from %s\n", file.OldPath))
			sb.WriteString(fmt.Sprintf("rename to %s\n", file.NewPath))
			if file.Similarity > 0 {
				sb.WriteString(fmt.Sprintf("similarity index %d%%\n", file.Similarity))
			}
		}

		if file.OldMode != file.NewMode && file.OldMode != "" && file.NewMode != "" {
			sb.WriteString(fmt.Sprintf("old mode %s\n", file.OldMode))
			sb.WriteString(fmt.Sprintf("new mode %s\n", file.NewMode))
		}

		// File markers
		if file.IsNew {
			sb.WriteString("--- /dev/null\n")
		} else {
			sb.WriteString(fmt.Sprintf("--- a/%s\n", file.OldPath))
		}

		if file.IsDeleted {
			sb.WriteString("+++ /dev/null\n")
		} else {
			sb.WriteString(fmt.Sprintf("+++ b/%s\n", file.NewPath))
		}

		// Add hunks
		for _, hunk := range file.Hunks {
			sb.WriteString(fmt.Sprintf("@@ -%d,%d +%d,%d @@\n",
				hunk.OldStart, hunk.OldCount,
				hunk.NewStart, hunk.NewCount))

			for _, line := range hunk.Lines {
				switch line.Type {
				case diff.LineAdded:
					sb.WriteString("+")
				case diff.LineRemoved:
					sb.WriteString("-")
				case diff.LineContext:
					sb.WriteString(" ")
				}
				sb.WriteString(line.Content)
				sb.WriteString("\n")
			}
		}

		sb.WriteString("\n")
	}

	return sb.String()
}

// parseReviewResponse parses the JSON response from Claude
func (rc *ReviewerClient) parseReviewResponse(response string) (*review.Result, error) {
	// Try to extract JSON from the response
	// Claude might wrap it in markdown code blocks
	jsonStr := rc.extractJSON(response)

	var result review.Result
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w (response: %s)", err, response)
	}

	return &result, nil
}

// extractJSON extracts JSON from a response that might be wrapped in markdown
func (rc *ReviewerClient) extractJSON(response string) string {
	// Remove markdown code blocks if present
	response = strings.TrimSpace(response)

	// Look for ```json ... ``` or ``` ... ```
	if strings.HasPrefix(response, "```") {
		lines := strings.Split(response, "\n")
		if len(lines) > 2 {
			// Remove first line (```json or ```)
			lines = lines[1:]
			// Remove last line if it's ```
			if strings.TrimSpace(lines[len(lines)-1]) == "```" {
				lines = lines[:len(lines)-1]
			}
			response = strings.Join(lines, "\n")
		}
	}

	// Find JSON object boundaries
	startIdx := strings.Index(response, "{")
	endIdx := strings.LastIndex(response, "}")

	if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
		return response[startIdx : endIdx+1]
	}

	return response
}
