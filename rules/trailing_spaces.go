package rules

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mikemadden42/txtlint/linter"
)

// TrailingSpacesRule reports lines that end with whitespace.
type TrailingSpacesRule struct{}

// NewTrailingSpacesRule creates a new instance of TrailingSpacesRule.
func NewTrailingSpacesRule() linter.Rule {
	return &TrailingSpacesRule{}
}

// Name returns the name of the rule.
func (r *TrailingSpacesRule) Name() string {
	return "TrailingSpaces"
}

// Description returns a one-line summary of the rule.
func (r *TrailingSpacesRule) Description() string {
	return "Lines must not end with whitespace."
}

// Check reports each line with trailing whitespace, pointing at the first
// trailing whitespace character.
func (r *TrailingSpacesRule) Check(doc *linter.Document) []linter.Issue {
	var issues []linter.Issue
	for _, line := range doc.Lines {
		trimmed := strings.TrimRightFunc(line.Text, unicode.IsSpace)
		if len(trimmed) == len(line.Text) {
			continue
		}
		issues = append(issues, linter.Issue{
			Rule:     r.Name(),
			Severity: linter.Error,
			Line:     line.Number,
			Column:   utf8.RuneCountInString(trimmed) + 1,
			Message:  "Line has trailing whitespace",
		})
	}
	return issues
}
