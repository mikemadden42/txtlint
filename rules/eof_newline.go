package rules

import (
	"fmt"
	"unicode/utf8"

	"github.com/mikemadden42/txtlint/linter"
)

// EOFNewlineRule ensures a non-empty file ends with exactly one line ending.
type EOFNewlineRule struct{}

// NewEOFNewlineRule creates a new instance of EOFNewlineRule.
func NewEOFNewlineRule() linter.Rule {
	return &EOFNewlineRule{}
}

// Name returns the name of the rule.
func (r *EOFNewlineRule) Name() string {
	return "EOFNewline"
}

// Description returns a one-line summary of the rule.
func (r *EOFNewlineRule) Description() string {
	return "Files must end with exactly one newline."
}

// Check reports a missing final newline, or blank lines at the end of the file.
// Files that are empty or contain only blank lines are not reported.
func (r *EOFNewlineRule) Check(doc *linter.Document) []linter.Issue {
	lines := doc.Lines
	if len(lines) == 0 {
		return nil
	}

	last := lines[len(lines)-1]
	if last.Ending == "" {
		return []linter.Issue{{
			Rule:     r.Name(),
			Severity: linter.Error,
			Line:     last.Number,
			Column:   utf8.RuneCountInString(last.Text) + 1,
			Message:  "File does not end with a newline",
		}}
	}

	blank := 0
	for i := len(lines) - 1; i >= 0 && lines[i].Text == ""; i-- {
		blank++
	}
	if blank == 0 || blank == len(lines) {
		return nil
	}
	return []linter.Issue{{
		Rule:     r.Name(),
		Severity: linter.Error,
		Line:     lines[len(lines)-blank].Number,
		Message:  fmt.Sprintf("File ends with %d extra newline(s)", blank),
	}}
}
