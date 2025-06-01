package rules

import (
	"bytes"
	"io/ioutil"

	"github.com/mikemadden42/txtlint/linter"
)

// NoMixedLineEndingsRule implements linter.FileAccessRule to check for consistent line endings.
type NoMixedLineEndingsRule struct {
	// No state needed during line-by-line scan, check happens in Finalize
}

// NewNoMixedLineEndingsRule creates a new instance of NoMixedLineEndingsRule.
func NewNoMixedLineEndingsRule() linter.FileAccessRule { // Returns FileAccessRule
	return &NoMixedLineEndingsRule{}
}

// Name returns the name of the rule.
func (r *NoMixedLineEndingsRule) Name() string {
	return "NoMixedLineEndings"
}

// LintLine does nothing for this rule, as line endings are checked in Finalize.
func (r *NoMixedLineEndingsRule) LintLine(line string, lineNumber int) []linter.LintError {
	return nil
}

// Finalize checks the entire file for mixed line endings.
func (r *NoMixedLineEndingsRule) Finalize(filePath string) []linter.LintError {
	if filePath == "" {
		return nil // Should not happen if main.go passes it correctly for FileAccessRule
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return []linter.LintError{{
			RuleName: r.Name(),
			Line:     0, // File-level error
			Message:  "Failed to read file for line ending check: " + err.Error(),
		}}
	}

	hasLF := bytes.ContainsRune(content, '\n')
	hasCRLF := bytes.Contains(content, []byte{'\r', '\n'})
	hasCR := bytes.ContainsRune(content, '\r') && !hasCRLF // Check for old Mac style CR

	var errors []linter.LintError

	// If both LF and CRLF are present, it's mixed
	if hasLF && hasCRLF {
		errors = append(errors, linter.LintError{
			RuleName: r.Name(),
			Line:     0, // File-level error
			Message:  "Mixed line endings (LF and CRLF) detected in file",
		})
	}
	// If CR is present (and not part of CRLF), it's also a mixed or non-standard ending
	if hasCR {
		errors = append(errors, linter.LintError{
			RuleName: r.Name(),
			Line:     0, // File-level error
			Message:  "Carriage Return (CR) line endings detected (old Mac format)",
		})
	}

	return errors
}
