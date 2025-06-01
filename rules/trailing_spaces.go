package rules

import (
	"unicode"

	"github.com/mikemadden42/txtlint/linter"
)

// TrailingSpacesRule implements the linter.Rule interface to check for trailing whitespace.
type TrailingSpacesRule struct{}

// NewTrailingSpacesRule creates a new instance of TrailingSpacesRule.
func NewTrailingSpacesRule() linter.Rule {
	return &TrailingSpacesRule{}
}

// Name returns the name of the rule.
func (r *TrailingSpacesRule) Name() string {
	return "TrailingSpaces"
}

// LintLine checks for trailing whitespace on the given line.
func (r *TrailingSpacesRule) LintLine(line string, lineNumber int) []linter.LintError {
	var errors []linter.LintError
	if len(line) > 0 && unicode.IsSpace(rune(line[len(line)-1])) {
		// Find the column where trailing spaces begin
		firstTrailingSpaceCol := len(line)
		for i := len(line) - 1; i >= 0; i-- {
			if !unicode.IsSpace(rune(line[i])) {
				break
			}
			firstTrailingSpaceCol = i + 1 // +1 for 1-based column number
		}
		errors = append(errors, linter.LintError{
			RuleName: r.Name(),
			Line:     lineNumber,
			Column:   firstTrailingSpaceCol,
			Message:  "Line has trailing whitespace",
		})
	}
	return errors
}

// Finalize does nothing for this rule as it's line-by-line.
func (r *TrailingSpacesRule) Finalize() []linter.LintError {
	return nil
}
