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

// Rule defines the interface that all line-by-line linting rules must implement.
type Rule interface {
	// Name returns the unique name of the linting rule.
	Name() string

	// LintLine processes a single line of text and returns any LintErrors found.
	// The 'line' parameter does NOT include the trailing newline character.
	LintLine(line string, lineNumber int) []LintError

	// Finalize is called after all lines in the file have been processed.
	// It's used for rules that need to perform checks that don't depend on raw file bytes.
	// Returns any LintErrors found globally by the rule.
	Finalize(filePath string) []LintError // Now takes filePath
}

// FileAccessRule is an optional interface for rules that need to read the raw file bytes
// (e.g., to check for mixed line endings or EOF newline).
// A rule can implement both Rule and FileAccessRule if it has both line-based and file-based checks.
type FileAccessRule interface {
	Rule // Embeds the base Rule interface

	// Finalize is overridden to specifically indicate that this rule might
	// need to re-read the file, and thus the filePath is guaranteed to be provided.
	Finalize(filePath string) []LintError // Explicitly takes filePath
}
