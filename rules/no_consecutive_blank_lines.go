package rules

import (
	"strings"

	"github.com/mikemadden42/txtlint/linter"
)

// NoConsecutiveBlankLinesRule implements the linter.Rule interface
// to check for more than one consecutive blank line.
type NoConsecutiveBlankLinesRule struct {
	previousLineWasBlank bool
}

// NewNoConsecutiveBlankLinesRule creates a new instance of NoConsecutiveBlankLinesRule.
func NewNoConsecutiveBlankLinesRule() linter.Rule {
	return &NoConsecutiveBlankLinesRule{}
}

// Name returns the name of the rule.
func (r *NoConsecutiveBlankLinesRule) Name() string {
	return "NoConsecutiveBlankLines"
}

// LintLine checks for consecutive blank lines.
func (r *NoConsecutiveBlankLinesRule) LintLine(line string, lineNumber int) []linter.LintError {
	var errors []linter.LintError
	currentLineIsBlank := (strings.TrimSpace(line) == "")

	if currentLineIsBlank && r.previousLineWasBlank {
		errors = append(errors, linter.LintError{
			RuleName: r.Name(),
			Line:     lineNumber,
			Message:  "More than one consecutive blank line detected",
		})
	}
	r.previousLineWasBlank = currentLineIsBlank
	return errors
}

// Finalize does nothing for this rule as its state is reset per file.
func (r *NoConsecutiveBlankLinesRule) Finalize(filePath string) []linter.LintError {
	return nil
}
