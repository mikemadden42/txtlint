package rules

import (
	"fmt"
	"strings"

	"github.com/mikemadden42/txtlint/linter"
)

// SuperLongSentenceRule implements the linter.Rule interface to check for long sentences.
type SuperLongSentenceRule struct {
	maxSentenceLength int
}

// NewSuperLongSentenceRule creates a new instance of SuperLongSentenceRule with a configurable maximum length.
func NewSuperLongSentenceRule(maxLength int) linter.Rule {
	return &SuperLongSentenceRule{
		maxSentenceLength: maxLength,
	}
}

// Name returns the name of the rule.
func (r *SuperLongSentenceRule) Name() string {
	return "SuperLongSentence"
}

// LintLine checks for super long sentences on the given line.
// This is a simplified approach; true sentence segmentation is complex.
func (r *SuperLongSentenceRule) LintLine(line string, lineNumber int) []linter.LintError {
	var errors []linter.LintError

	// A very basic way to split by common sentence terminators.
	// This will not be perfect for all cases (e.g., abbreviations, quoted text).
	sentences := strings.FieldsFunc(line, func(r rune) bool {
		return r == '.' || r == '!' || r == '?'
	})

	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if len(sentence) > r.maxSentenceLength {
			// For simplicity, we'll just report the line number.
			// A column number for a whole sentence is less meaningful here.
			errors = append(errors, linter.LintError{
				RuleName: r.Name(),
				Line:     lineNumber,
				Message:  fmt.Sprintf("Sentence exceeds %d characters (length: %d)", r.maxSentenceLength, len(sentence)),
			})
		}
	}
	return errors
}

// Finalize does nothing for this rule as it's purely line-by-line.
func (r *SuperLongSentenceRule) Finalize(filePath string) []linter.LintError {
	return nil
}
