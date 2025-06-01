package main

import (
	"bufio"
	"fmt"
	"io" // Import for io.EOF
	"os"

	"github.com/mikemadden42/txtlint/linter"
	"github.com/mikemadden42/txtlint/rules"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: txtlint <file.txt>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filePath, err)
		os.Exit(1)
	}
	defer file.Close() // Ensure the file is closed when main exits

	// Initialize your rules.
	// In a real application, you might read these from a config file
	// or command-line flags to enable/disable specific rules.
	lintRules := []linter.Rule{
		rules.NewTrailingSpacesRule(),
		rules.NewSuperLongSentenceRule(120), // Example: Configure max length
		rules.NewMissingPunctuationRule(),
		rules.NewNoConsecutiveBlankLinesRule(), // New rule!
		rules.NewNoMixedLineEndingsRule(),      // New rule!
		rules.NewEOFNewlineRule(),              // New rule!
		// Add instances of other rules here as you create them
	}

	var allErrors []linter.LintError
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text() // scanner.Text() strips the newline characters

		// Store raw line (including original newline if present) if needed by Finalize() rules
		// To get the raw line *including* the newline, we need to manually read or adapt.
		// For simplicity, let's assume scanner.Bytes() works, or handle it as a specific rule's concern.
		// For NoMixedLineEndingsRule, we need to read the raw bytes.
		// Let's modify the loop to provide raw line data.
		// For simplicity in this main.go, we'll let NoMixedLineEndingsRule handle its own file reading.
		// The `Finalize` method will likely need to re-read the file in its own way if it needs raw byte content.

		for _, rule := range lintRules {
			allErrors = append(allErrors, rule.LintLine(line, lineNumber)...)
		}
	}

	if err := scanner.Err(); err != nil {
		// bufio.Scanner only returns io.EOF if the last token is empty and
		// the reader returns io.EOF. For our purposes, other errors are more critical.
		if err != io.EOF {
			fmt.Printf("Error reading file %s: %v\n", filePath, err)
			os.Exit(1)
		}
	}

	// After processing all lines, call Finalize for rules that need it.
	// Pass the file path to Finalize if rules need to re-read the file raw.
	for _, rule := range lintRules {
		if ruleWithFileAccess, ok := rule.(linter.FileAccessRule); ok {
			allErrors = append(allErrors, ruleWithFileAccess.Finalize(filePath)...)
		} else {
			allErrors = append(allErrors, rule.Finalize("")...) // Pass empty string if file path isn't needed
		}
	}

	// Report results
	if len(allErrors) > 0 {
		fmt.Printf("\nLinting results for %s:\n", filePath)
		for _, err := range allErrors {
			fmt.Println(err)
		}
		os.Exit(1) // Exit with non-zero status to indicate errors for scripting
	} else {
		fmt.Printf("No linting issues found in %s.\n", filePath)
		os.Exit(0) // Exit with zero status for success
	}
}
