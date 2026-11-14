package rules

import (
	"fmt"
	"strings"

	"github.com/mikemadden42/txtlint/linter"
)

// lineEndings lists the recognized line endings in reporting order.
var lineEndings = []struct{ seq, name string }{
	{"\n", "LF"},
	{"\r\n", "CRLF"},
	{"\r", "CR"},
}

// NoMixedLineEndingsRule reports files that mix LF, CRLF and CR line endings,
// and files that use classic Mac (CR) line endings at all.
type NoMixedLineEndingsRule struct{}

// NewNoMixedLineEndingsRule creates a new instance of NoMixedLineEndingsRule.
func NewNoMixedLineEndingsRule() linter.Rule {
	return &NoMixedLineEndingsRule{}
}

// Name returns the name of the rule.
func (r *NoMixedLineEndingsRule) Name() string {
	return "NoMixedLineEndings"
}

// Description returns a one-line summary of the rule.
func (r *NoMixedLineEndingsRule) Description() string {
	return "All lines must use the same line ending, either LF or CRLF."
}

// Check reports the first CR line ending, and, if endings are mixed, the first
// line that differs from the most common ending.
func (r *NoMixedLineEndingsRule) Check(doc *linter.Document) []linter.Issue {
	counts := map[string]int{}
	first := map[string]int{} // first line number using each ending
	for _, line := range doc.Lines {
		if line.Ending == "" {
			continue
		}
		if counts[line.Ending] == 0 {
			first[line.Ending] = line.Number
		}
		counts[line.Ending]++
	}

	var issues []linter.Issue
	if counts["\r"] > 0 {
		issues = append(issues, linter.Issue{
			Rule:     r.Name(),
			Severity: linter.Error,
			Line:     first["\r"],
			Message:  "Carriage return (CR) line ending (classic Mac format)",
		})
	}
	if len(counts) < 2 {
		return issues
	}

	// The more common of LF and CRLF is the expected ending (CR is already
	// reported above, so it is never expected); ties go to whichever appears
	// first. Report the first line that uses any other ending.
	dominant := ""
	for _, e := range lineEndings {
		if e.seq == "\r" {
			continue
		}
		c := counts[e.seq]
		if c > counts[dominant] || (c > 0 && c == counts[dominant] && first[e.seq] < first[dominant]) {
			dominant = e.seq
		}
	}
	var summary []string
	deviation, expected := 0, ""
	for _, e := range lineEndings {
		if counts[e.seq] == 0 {
			continue
		}
		summary = append(summary, fmt.Sprintf("%d %s", counts[e.seq], e.name))
		if e.seq == dominant {
			expected = e.name
		} else if deviation == 0 || first[e.seq] < deviation {
			deviation = first[e.seq]
		}
	}
	issues = append(issues, linter.Issue{
		Rule:     r.Name(),
		Severity: linter.Error,
		Line:     deviation,
		Message:  fmt.Sprintf("Mixed line endings: %s (expected %s)", strings.Join(summary, ", "), expected),
	})
	return issues
}
