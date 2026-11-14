// Package linter defines the types shared by txtlint rules and runs rules
// against a document.
package linter

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
)

// Severity indicates how serious an issue is.
type Severity int

const (
	// Error marks an issue that is definitely a problem.
	Error Severity = iota
	// Warning marks a heuristic issue that may be a false positive.
	Warning
)

// String returns the lowercase name of the severity.
func (s Severity) String() string {
	if s == Warning {
		return "warning"
	}
	return "error"
}

// Issue is a single problem found by a rule.
type Issue struct {
	Rule     string
	Severity Severity
	Line     int // 1-based; 0 means the issue applies to the whole file
	Column   int // 1-based, in runes; 0 means no specific column
	Message  string
}

// Format renders the issue as "file:line:col: severity: message [Rule]",
// omitting any position parts that are not set.
func (i Issue) Format(file string) string {
	var pos []string
	if file != "" {
		pos = append(pos, file)
	}
	if i.Line > 0 {
		pos = append(pos, fmt.Sprint(i.Line))
		if i.Column > 0 {
			pos = append(pos, fmt.Sprint(i.Column))
		}
	}
	prefix := ""
	if len(pos) > 0 {
		prefix = strings.Join(pos, ":") + ": "
	}
	return fmt.Sprintf("%s%s: %s [%s]", prefix, i.Severity, i.Message, i.Rule)
}

// String renders the issue without a file name.
func (i Issue) String() string {
	return i.Format("")
}

// Line is a single line of a document.
type Line struct {
	Number int    // 1-based
	Text   string // the line without its line ending
	Ending string // "\n", "\r\n", "\r", or "" for a final line with no terminator
}

// Document is a text file split into lines.
type Document struct {
	Content []byte
	Lines   []Line
}

// NewDocument splits content into lines, recognizing LF, CRLF and CR endings.
// A trailing line ending does not produce an extra empty line.
func NewDocument(content []byte) *Document {
	doc := &Document{Content: content}
	s := string(content)
	for n := 1; len(s) > 0; n++ {
		i := strings.IndexAny(s, "\r\n")
		if i < 0 {
			doc.Lines = append(doc.Lines, Line{Number: n, Text: s})
			break
		}
		end := i + 1
		if s[i] == '\r' && end < len(s) && s[end] == '\n' {
			end++
		}
		doc.Lines = append(doc.Lines, Line{Number: n, Text: s[:i], Ending: s[i:end]})
		s = s[end:]
	}
	return doc
}

// Rule is a single check run against a whole document.
type Rule interface {
	// Name returns the unique name of the rule.
	Name() string
	// Description returns a one-line summary of what the rule checks.
	Description() string
	// Check returns the issues the rule finds in doc.
	Check(doc *Document) []Issue
}

// Lint runs rules against content and returns their issues ordered by
// position. File-level issues (line 0) come first.
func Lint(content []byte, rules []Rule) []Issue {
	doc := NewDocument(content)
	var issues []Issue
	for _, rule := range rules {
		issues = append(issues, rule.Check(doc)...)
	}
	slices.SortStableFunc(issues, func(a, b Issue) int {
		return cmp.Or(cmp.Compare(a.Line, b.Line), cmp.Compare(a.Column, b.Column))
	})
	return issues
}
