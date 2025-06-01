package main

import (
	"bufio"
	"fmt"
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
	defer file.Close()

	// Initialize your rules.
	// In a real application, you might read these from a config file
	// or command-line flags to enable/disable specific rules.
	lintRules := []linter.Rule{
		rules.NewTrailingSpacesRule(),       // Use New functions for rules with potential config
		rules.NewSuperLongSentenceRule(120), // Example: Configure max length
		rules.NewMissingPunctuationRule(),
		// Add instances of other rules here as you create them
	}

	var allErrors []linter.LintError
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text() // scanner.Text() strips the newline characters

		for _, rule := range lintRules {
			allErrors = append(allErrors, rule.LintLine(line, lineNumber)...)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
		os.Exit(1)
	}

	// After processing all lines, call Finalize for rules that need it.
	for _, rule := range lintRules {
		allErrors = append(allErrors, rule.Finalize()...)
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
