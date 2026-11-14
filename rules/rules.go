// Package rules contains the built-in txtlint rules.
package rules

import "github.com/mikemadden42/txtlint/linter"

// All returns every built-in rule, configured with the given maximum
// sentence length.
func All(maxSentenceLength int) []linter.Rule {
	return []linter.Rule{
		NewTrailingSpacesRule(),
		NewSuperLongSentenceRule(maxSentenceLength),
		NewMissingPunctuationRule(),
		NewNoConsecutiveBlankLinesRule(),
		NewNoMixedLineEndingsRule(),
		NewEOFNewlineRule(),
	}
}
