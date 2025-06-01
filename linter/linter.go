package linter

import (
	"fmt"
)

// LintError represents a single linting issue found.
type LintError struct {
	RuleName string
	Line     int
	Column   int // Optional, but useful for pointing to exact location
	Message  string
}

// String returns a formatted string representation of the LintError.
func (e LintError) String() string {
	if e.Column > 0 {
		return fmt.Sprintf("[%s] Line %d, Col %d: %s", e.RuleName, e.Line, e.Column, e.Message)
	}
	return fmt.Sprintf("[%s] Line %d: %s", e.RuleName, e.Line, e.Message)
}

// Rule defines the interface that all linting rules must implement.
type Rule interface {
	// Name returns the unique name of the linting rule.
	Name() string

	// LintLine processes a single line of text and returns any LintErrors found.
	// The 'line' parameter does NOT include the trailing newline character.
	LintLine(line string, lineNumber int) []LintError

	// Finalize is called after all lines in the file have been processed.
	// It's used for rules that need to perform checks on the entire file content,
	// or require aggregated information. Returns any LintErrors found globally.
	Finalize() []LintError
}
