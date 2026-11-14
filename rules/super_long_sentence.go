package rules

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/mikemadden42/txtlint/linter"
)

// SuperLongSentenceRule reports sentences longer than a configured number of
// characters. Sentences may wrap across lines; a blank line ends a paragraph
// and therefore a sentence.
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

// Description returns a one-line summary of the rule.
func (r *SuperLongSentenceRule) Description() string {
	return fmt.Sprintf("Sentences should be at most %d characters (warning).", r.maxSentenceLength)
}

// Check reports each overly long sentence on the line where it starts.
// Length is measured in characters, with runs of whitespace (including line
// breaks) counted as a single space.
func (r *SuperLongSentenceRule) Check(doc *linter.Document) []linter.Issue {
	var issues []linter.Issue
	length, start := 0, 0
	flush := func() {
		if length > r.maxSentenceLength {
			issues = append(issues, linter.Issue{
				Rule:     r.Name(),
				Severity: linter.Warning,
				Line:     start,
				Message:  fmt.Sprintf("Sentence exceeds %d characters (length: %d)", r.maxSentenceLength, length),
			})
		}
		length = 0
	}

	for _, line := range doc.Lines {
		words := strings.Fields(line.Text)
		if len(words) == 0 {
			flush()
			continue
		}
		for _, word := range words {
			if length == 0 {
				start = line.Number
			} else {
				length++ // the space joining this word to the previous one
			}
			length += utf8.RuneCountInString(word)
			if endsSentence(word) {
				flush()
			}
		}
	}
	flush()
	return issues
}

// endsSentence reports whether word ends with sentence-ending punctuation,
// optionally followed by closing quotes or brackets. Punctuation inside a
// word, as in "3.50", does not end a sentence.
func endsSentence(word string) bool {
	word = strings.TrimRight(word, `"')]”’`)
	return strings.HasSuffix(word, ".") || strings.HasSuffix(word, "!") || strings.HasSuffix(word, "?")
}
