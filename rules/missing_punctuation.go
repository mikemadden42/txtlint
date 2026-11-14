package rules

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mikemadden42/txtlint/linter"
)

// MissingPunctuationRule reports lines that look like sentences but do not
// end with sentence-ending punctuation. It is a heuristic, so its issues are
// warnings.
type MissingPunctuationRule struct{}

// NewMissingPunctuationRule creates a new instance of MissingPunctuationRule.
func NewMissingPunctuationRule() linter.Rule {
	return &MissingPunctuationRule{}
}

// Name returns the name of the rule.
func (r *MissingPunctuationRule) Name() string {
	return "MissingPunctuation"
}

// Description returns a one-line summary of the rule.
func (r *MissingPunctuationRule) Description() string {
	return "Sentence-like lines should end with . ! or ? (warning)."
}

// Check reports lines of more than one word that end in a letter or digit.
func (r *MissingPunctuationRule) Check(doc *linter.Document) []linter.Issue {
	var issues []linter.Issue
	for _, line := range doc.Lines {
		text := strings.TrimRightFunc(line.Text, unicode.IsSpace)
		last, _ := utf8.DecodeLastRuneInString(text)
		if !unicode.IsLetter(last) && !unicode.IsDigit(last) {
			continue
		}
		// Single words are usually headings or list items, not sentences.
		if !strings.ContainsFunc(strings.TrimSpace(text), unicode.IsSpace) {
			continue
		}
		issues = append(issues, linter.Issue{
			Rule:     r.Name(),
			Severity: linter.Warning,
			Line:     line.Number,
			Column:   utf8.RuneCountInString(text),
			Message:  "Line might be missing sentence-ending punctuation (.!?)",
		})
	}
	return issues
}
