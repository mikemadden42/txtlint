package rules

import (
	"strings"
	"unicode"

	"github.com/mikemadden42/txtlint/linter"
)

// MissingPunctuationRule implements the linter.Rule interface to check for missing end punctuation.
type MissingPunctuationRule struct{}

// NewMissingPunctuationRule creates a new instance of MissingPunctuationRule.
func NewMissingPunctuationRule() linter.Rule {
	return &MissingPunctuationRule{}
}

// Name returns the name of the rule.
func (r *MissingPunctuationRule) Name() string {
	return "MissingPunctuation"
}

// LintLine checks for lines that might be missing a sentence-ending punctuation.
// This is a heuristic and not foolproof.
func (r *MissingPunctuationRule) LintLine(line string, lineNumber int) []linter.LintError {
	var errors []linter.LintError

	trimmedLine := strings.TrimSpace(line)

	if len(trimmedLine) == 0 {
		return errors // Ignore empty lines
	}

	lastChar := rune(trimmedLine[len(trimmedLine)-1])

	// Define common sentence-ending punctuation
	isPunctuation := func(r rune) bool {
		return r == '.' || r == '!' || r == '?'
	}

	// Check if the line ends with a letter or digit, and not with punctuation.
	// Also ignore lines that seem to be code, paths, or short fragments.
	if unicode.IsLetter(lastChar) || unicode.IsDigit(lastChar) {
		if !isPunctuation(lastChar) {
			// If the line contains at least one space, it's more likely a sentence.
			// This helps avoid flagging single words or very short fragments.
			if strings.ContainsRune(trimmedLine, ' ') {
				errors = append(errors, linter.LintError{
					RuleName: r.Name(),
					Line:     lineNumber,
					Column:   len(trimmedLine), // Point to the last character
					Message:  "Line might be missing sentence-ending punctuation (.!?)",
				})
			}
		}
	}
	return errors
}

// Finalize does nothing for this rule as it's line-by-line.
func (r *MissingPunctuationRule) Finalize() []linter.LintError {
	return nil
}
