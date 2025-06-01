package rules

import (
	"io/ioutil"
	"strings"

	"github.com/mikemadden42/txtlint/linter"
)

// EOFNewlineRule implements linter.FileAccessRule to ensure the file ends with a single newline.
type EOFNewlineRule struct {
	// No state needed
}

// NewEOFNewlineRule creates a new instance of EOFNewlineRule.
func NewEOFNewlineRule() linter.FileAccessRule { // Returns FileAccessRule
	return &EOFNewlineRule{}
}

// Name returns the name of the rule.
func (r *EOFNewlineRule) Name() string {
	return "EOFNewline"
}

// LintLine does nothing for this rule, as it's a file-level check.
func (r *EOFNewlineRule) LintLine(line string, lineNumber int) []linter.LintError {
	return nil
}

// Finalize checks if the file ends with a single newline character.
func (r *EOFNewlineRule) Finalize(filePath string) []linter.LintError {
	if filePath == "" {
		return nil // Should not happen
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return []linter.LintError{{
			RuleName: r.Name(),
			Line:     0, // File-level error
			Message:  "Failed to read file for EOF newline check: " + err.Error(),
		}}
	}

	var errors []linter.LintError

	if len(content) == 0 {
		return nil // Empty file, no newline needed
	}

	// Check if the last character is a newline
	if content[len(content)-1] != '\n' {
		errors = append(errors, linter.LintError{
			RuleName: r.Name(),
			Line:     len(strings.Split(string(content), "\n")), // Report on the last logical line
			Message:  "File does not end with a newline character",
		})
	} else if len(content) >= 2 && content[len(content)-2] == '\n' && content[len(content)-1] == '\n' {
		// Check for multiple consecutive newlines at EOF (e.g., "\n\n")
		// This handles both LF and CRLF files where the very last char is '\n'
		// but the char before it is also '\n' (or '\r' followed by '\n' but then another '\n')
		// This specifically targets the scenario of *more than one* trailing newline.
		// For example: "text\n\n" should be flagged as "more than one".
		// "text\n" should be okay.
		// "text\r\n" should be okay.
		// "text\r\n\r\n" should be flagged.
		if len(content) >= 2 && content[len(content)-1] == '\n' {
			// Check for two consecutive newlines at the very end
			if content[len(content)-2] == '\n' { // ...\n\n
				errors = append(errors, linter.LintError{
					RuleName: r.Name(),
					Line:     len(strings.Split(string(content), "\n")),
					Message:  "File ends with more than one newline character",
				})
			} else if len(content) >= 3 && content[len(content)-1] == '\n' && content[len(content)-2] == '\r' && content[len(content)-3] == '\n' {
				// Special case for CRLF: ...\n\r\n
				// If it was ...LFCRLF then the second to last char would be '\r'.
				// This specifically targets `\n\n` or `\r\n\n` at the very end.
				// For CRLF, we want to ensure it ends *only* with \r\n.
				// If it's \r\n\r\n, we need to flag.
				// This check is tricky. A simpler way:
				// If you trim all final newlines, should only one remain?
				cleaned := strings.TrimRight(string(content), "\r\n")
				if len(content)-len(cleaned) > 1 { // More than one newline sequence at end
					errors = append(errors, linter.LintError{
						RuleName: r.Name(),
						Line:     len(strings.Split(string(content), "\n")),
						Message:  "File ends with more than one newline character",
					})
				}
			}
		}
	}

	return errors
}
