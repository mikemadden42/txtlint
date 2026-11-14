// Command txtlint checks plain-text files for common formatting and style
// problems.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/mikemadden42/txtlint/linter"
	"github.com/mikemadden42/txtlint/rules"
)

// Exit codes.
const (
	exitOK     = 0 // no errors found
	exitIssues = 1 // lint errors found (or warnings, with -strict)
	exitFailed = 2 // bad usage, or a file could not be read
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

// run executes txtlint with the given arguments and returns the exit code.
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("txtlint", flag.ContinueOnError)
	fs.SetOutput(stderr)
	maxLength := fs.Int("max-sentence-length", 120, "maximum sentence length, in characters, for SuperLongSentence")
	disable := fs.String("disable", "", "comma-separated list of rules to disable")
	strict := fs.Bool("strict", false, "exit with status 1 on warnings as well as errors")
	listRules := fs.Bool("list-rules", false, "list the available rules and exit")
	fs.Usage = func() {
		fmt.Fprint(stderr, "Usage: txtlint [flags] [file ...]\n\n"+
			"Lints each file, or standard input if no file (or \"-\") is given.\n\nFlags:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return exitOK
		}
		return exitFailed
	}
	if *maxLength < 1 {
		fmt.Fprintln(stderr, "txtlint: -max-sentence-length must be at least 1")
		return exitFailed
	}

	all := rules.All(*maxLength)
	if *listRules {
		for _, rule := range all {
			fmt.Fprintf(stdout, "%-24s %s\n", rule.Name(), rule.Description())
		}
		return exitOK
	}
	enabled, err := enabledRules(all, *disable)
	if err != nil {
		fmt.Fprintf(stderr, "txtlint: %v\n", err)
		return exitFailed
	}

	paths := fs.Args()
	if len(paths) == 0 {
		paths = []string{"-"}
	}

	code := exitOK
	for _, path := range paths {
		content, err := readInput(path, stdin)
		if err != nil {
			fmt.Fprintf(stderr, "txtlint: %v\n", err)
			code = exitFailed
			continue
		}
		name := path
		if path == "-" {
			name = "<stdin>"
		}
		for _, issue := range linter.Lint(content, enabled) {
			fmt.Fprintln(stdout, issue.Format(name))
			if code == exitOK && (issue.Severity == linter.Error || *strict) {
				code = exitIssues
			}
		}
	}
	return code
}

// enabledRules returns all rules except those named in disable, a
// comma-separated, case-insensitive list.
func enabledRules(all []linter.Rule, disable string) ([]linter.Rule, error) {
	disabled := map[string]bool{}
	for _, name := range strings.Split(disable, ",") {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		known := slices.ContainsFunc(all, func(r linter.Rule) bool {
			return strings.EqualFold(r.Name(), name)
		})
		if !known {
			return nil, fmt.Errorf("unknown rule %q (see -list-rules)", name)
		}
		disabled[strings.ToLower(name)] = true
	}

	var enabled []linter.Rule
	for _, rule := range all {
		if !disabled[strings.ToLower(rule.Name())] {
			enabled = append(enabled, rule)
		}
	}
	return enabled, nil
}

// readInput reads the named file, or stdin when path is "-".
func readInput(path string, stdin io.Reader) ([]byte, error) {
	if path == "-" {
		return io.ReadAll(stdin)
	}
	return os.ReadFile(path) //nolint:gosec // G304: reading user-specified files is the tool's purpose.
}
