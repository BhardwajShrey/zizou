package review

// Severity represents the severity level of a review comment
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
	SeverityInfo     Severity = "info"
)

// Category represents the category of a review comment
type Category string

const (
	CategorySecurity     Category = "security"
	CategoryPerformance  Category = "performance"
	CategoryBugRisk      Category = "bug_risk"
	CategoryMaintenance  Category = "maintainability"
	CategoryStyle        Category = "style"
	CategoryBestPractice Category = "best_practice"
	CategoryOther        Category = "other"
)

// Comment represents a code review comment
type Comment struct {
	File     string   `json:"file"`
	Line     int      `json:"line"`
	Severity Severity `json:"severity"`
	Category Category `json:"category"`
	Message  string   `json:"message"`
	Snippet  string   `json:"snippet,omitempty"`
}

// Result represents the complete review result
type Result struct {
	Comments []Comment `json:"comments"`
	Summary  string    `json:"summary"`
}
