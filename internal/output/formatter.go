package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/BhardwajShrey/zizou/internal/review"
)

// Format represents the output format
type Format string

const (
	FormatText     Format = "text"
	FormatJSON     Format = "json"
	FormatMarkdown Format = "markdown"
)

// Formatter formats review results for output
type Formatter struct {
	format Format
}

// NewFormatter creates a new formatter
func NewFormatter(format string) *Formatter {
	return &Formatter{
		format: Format(format),
	}
}

// Format formats a review result according to the configured format
func (f *Formatter) Format(result *review.Result) (string, error) {
	switch f.format {
	case FormatJSON:
		return f.formatJSON(result)
	case FormatMarkdown:
		return f.formatMarkdown(result)
	case FormatText:
		return f.formatText(result)
	default:
		return "", fmt.Errorf("unsupported format: %s", f.format)
	}
}

func (f *Formatter) formatJSON(result *review.Result) (string, error) {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (f *Formatter) formatText(result *review.Result) (string, error) {
	var sb strings.Builder

	sb.WriteString("Code Review Results\n")
	sb.WriteString(strings.Repeat("=", 60))
	sb.WriteString("\n\n")

	if len(result.Comments) == 0 {
		sb.WriteString("No issues found.\n")
	} else {
		for i, comment := range result.Comments {
			sb.WriteString(fmt.Sprintf("[%d] %s:%d\n", i+1, comment.File, comment.Line))
			sb.WriteString(fmt.Sprintf("    Severity: %s | Category: %s\n", comment.Severity, comment.Category))
			sb.WriteString(fmt.Sprintf("    %s\n", comment.Message))
			if comment.Snippet != "" {
				sb.WriteString(fmt.Sprintf("    Code: %s\n", comment.Snippet))
			}
			sb.WriteString("\n")
		}
	}

	if result.Summary != "" {
		sb.WriteString(strings.Repeat("-", 60))
		sb.WriteString("\n")
		sb.WriteString("Summary:\n")
		sb.WriteString(result.Summary)
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

func (f *Formatter) formatMarkdown(result *review.Result) (string, error) {
	var sb strings.Builder

	sb.WriteString("# Code Review Results\n\n")

	if len(result.Comments) == 0 {
		sb.WriteString("✅ No issues found.\n")
	} else {
		// Group by severity
		bySeverity := make(map[review.Severity][]review.Comment)
		for _, comment := range result.Comments {
			bySeverity[comment.Severity] = append(bySeverity[comment.Severity], comment)
		}

		// Output in severity order
		severities := []review.Severity{
			review.SeverityCritical,
			review.SeverityHigh,
			review.SeverityMedium,
			review.SeverityLow,
			review.SeverityInfo,
		}

		for _, severity := range severities {
			comments := bySeverity[severity]
			if len(comments) == 0 {
				continue
			}

			icon := getSeverityIcon(severity)
			sb.WriteString(fmt.Sprintf("## %s %s (%d)\n\n", icon, strings.Title(string(severity)), len(comments)))

			for _, comment := range comments {
				sb.WriteString(fmt.Sprintf("### %s:%d\n", comment.File, comment.Line))
				sb.WriteString(fmt.Sprintf("**Category:** %s\n\n", comment.Category))
				sb.WriteString(fmt.Sprintf("%s\n\n", comment.Message))
				if comment.Snippet != "" {
					sb.WriteString("```\n")
					sb.WriteString(comment.Snippet)
					sb.WriteString("\n```\n\n")
				}
			}
		}
	}

	if result.Summary != "" {
		sb.WriteString("---\n\n")
		sb.WriteString("## Summary\n\n")
		sb.WriteString(result.Summary)
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

func getSeverityIcon(severity review.Severity) string {
	switch severity {
	case review.SeverityCritical:
		return "🔴"
	case review.SeverityHigh:
		return "🟠"
	case review.SeverityMedium:
		return "🟡"
	case review.SeverityLow:
		return "🔵"
	case review.SeverityInfo:
		return "ℹ️"
	default:
		return "•"
	}
}
