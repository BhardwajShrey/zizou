package review

import (
	"fmt"
	"strings"

	"github.com/BhardwajShrey/zizou/internal/diff"
)

// BuildPrompt creates a prompt for Claude to review the diff
func BuildPrompt(d *diff.Diff) string {
	var sb strings.Builder

	sb.WriteString("You are an expert code reviewer. Review the following git diff and provide structured feedback.\n\n")
	sb.WriteString("Focus on:\n")
	sb.WriteString("- Security vulnerabilities\n")
	sb.WriteString("- Potential bugs or logic errors\n")
	sb.WriteString("- Performance issues\n")
	sb.WriteString("- Code maintainability and readability\n")
	sb.WriteString("- Best practices violations\n\n")

	sb.WriteString("Git diff:\n")
	sb.WriteString("```\n")
	sb.WriteString(formatDiff(d))
	sb.WriteString("```\n\n")

	sb.WriteString("Provide your review in the following JSON format:\n")
	sb.WriteString("{\n")
	sb.WriteString("  \"comments\": [\n")
	sb.WriteString("    {\n")
	sb.WriteString("      \"file\": \"path/to/file.go\",\n")
	sb.WriteString("      \"line\": 42,\n")
	sb.WriteString("      \"severity\": \"high\",  // critical, high, medium, low, info\n")
	sb.WriteString("      \"category\": \"security\",  // security, performance, bug_risk, maintainability, style, best_practice, other\n")
	sb.WriteString("      \"message\": \"Description of the issue\",\n")
	sb.WriteString("      \"snippet\": \"relevant code snippet (optional)\"\n")
	sb.WriteString("    }\n")
	sb.WriteString("  ],\n")
	sb.WriteString("  \"summary\": \"Overall review summary\"\n")
	sb.WriteString("}\n\n")

	sb.WriteString("IMPORTANT: Return ONLY valid JSON. Do not include any explanatory text before or after the JSON.\n")

	return sb.String()
}

// formatDiff converts a parsed Diff back into a string representation
func formatDiff(d *diff.Diff) string {
	var sb strings.Builder

	for _, file := range d.Files {
		sb.WriteString(fmt.Sprintf("diff --git a/%s b/%s\n", file.OldPath, file.NewPath))
		sb.WriteString(fmt.Sprintf("--- a/%s\n", file.OldPath))
		sb.WriteString(fmt.Sprintf("+++ b/%s\n", file.NewPath))

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
	}

	return sb.String()
}
