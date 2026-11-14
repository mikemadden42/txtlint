package rules

import (
	"strings"

	"github.com/mikemadden42/txtlint/linter"
)

// NoConsecutiveBlankLinesRule reports runs of more than one blank line.
type NoConsecutiveBlankLinesRule struct{}

// NewNoConsecutiveBlankLinesRule creates a new instance of NoConsecutiveBlankLinesRule.
func NewNoConsecutiveBlankLinesRule() linter.Rule {
	return &NoConsecutiveBlankLinesRule{}
}

// Name returns the name of the rule.
func (r *NoConsecutiveBlankLinesRule) Name() string {
	return "NoConsecutiveBlankLines"
}

// Description returns a one-line summary of the rule.
func (r *NoConsecutiveBlankLinesRule) Description() string {
	return "There must not be more than one blank line in a row."
}

// Check reports every blank line that directly follows another blank line.
func (r *NoConsecutiveBlankLinesRule) Check(doc *linter.Document) []linter.Issue {
	var issues []linter.Issue
	previousLineWasBlank := false
	for _, line := range doc.Lines {
		blank := strings.TrimSpace(line.Text) == ""
		if blank && previousLineWasBlank {
			issues = append(issues, linter.Issue{
				Rule:     r.Name(),
				Severity: linter.Error,
				Line:     line.Number,
				Message:  "More than one consecutive blank line",
			})
		}
		previousLineWasBlank = blank
	}
	return issues
}
