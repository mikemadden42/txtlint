package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// sample triggers every rule except SuperLongSentence.
const sample = "Trailing spaces here.   \n" + // 1: TrailingSpaces
	"\n" +
	"\n" + // 3: NoConsecutiveBlankLines
	"This line has no final punctuation\r\n" + // 4: MissingPunctuation, NoMixedLineEndings
	"Last line without newline." // 5: EOFNewline

const sampleOutput = `<stdin>:1:22: error: Line has trailing whitespace [TrailingSpaces]
<stdin>:3: error: More than one consecutive blank line [NoConsecutiveBlankLines]
<stdin>:4: error: Mixed line endings: 3 LF, 1 CRLF (expected LF) [NoMixedLineEndings]
<stdin>:4:34: warning: Line might be missing sentence-ending punctuation (.!?) [MissingPunctuation]
<stdin>:5:27: error: File does not end with a newline [EOFNewline]
`

func runTxtlint(t *testing.T, stdin string, args ...string) (code int, stdout, stderr string) {
	t.Helper()
	var out, errOut bytes.Buffer
	code = run(args, strings.NewReader(stdin), &out, &errOut)
	return code, out.String(), errOut.String()
}

func writeFile(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestRunStdin(t *testing.T) {
	code, stdout, stderr := runTxtlint(t, sample)
	if code != exitIssues {
		t.Errorf("exit code = %d, want %d", code, exitIssues)
	}
	if stdout != sampleOutput {
		t.Errorf("stdout:\n%s\nwant:\n%s", stdout, sampleOutput)
	}
	if stderr != "" {
		t.Errorf("unexpected stderr: %s", stderr)
	}
}

func TestRunDashMeansStdin(t *testing.T) {
	code, stdout, _ := runTxtlint(t, sample, "-")
	if code != exitIssues || stdout != sampleOutput {
		t.Errorf("code = %d, stdout:\n%s", code, stdout)
	}
}

func TestRunCleanFile(t *testing.T) {
	path := writeFile(t, "clean.txt", "All good here.\n")
	code, stdout, stderr := runTxtlint(t, "", path)
	if code != exitOK || stdout != "" || stderr != "" {
		t.Errorf("code = %d, stdout = %q, stderr = %q", code, stdout, stderr)
	}
}

func TestRunMultipleFiles(t *testing.T) {
	clean := writeFile(t, "clean.txt", "All good here.\n")
	dirty := writeFile(t, "dirty.txt", "No newline.")
	code, stdout, _ := runTxtlint(t, "", clean, dirty)
	want := dirty + ":1:12: error: File does not end with a newline [EOFNewline]\n"
	if code != exitIssues || stdout != want {
		t.Errorf("code = %d, stdout = %q, want %q", code, stdout, want)
	}
}

func TestRunWarningsOnly(t *testing.T) {
	const input = "No punctuation here\n"
	if code, stdout, _ := runTxtlint(t, input); code != exitOK || stdout == "" {
		t.Errorf("without -strict: code = %d, stdout = %q; want 0 with a warning", code, stdout)
	}
	if code, _, _ := runTxtlint(t, input, "-strict"); code != exitIssues {
		t.Errorf("with -strict: code = %d, want %d", code, exitIssues)
	}
}

func TestRunDisable(t *testing.T) {
	code, stdout, _ := runTxtlint(t, "No newline", "-disable", "eofnewline, MissingPunctuation")
	if code != exitOK || stdout != "" {
		t.Errorf("code = %d, stdout = %q", code, stdout)
	}
}

func TestRunMaxSentenceLength(t *testing.T) {
	code, stdout, _ := runTxtlint(t, "One two three four.\n", "-max-sentence-length", "10")
	want := "<stdin>:1: warning: Sentence exceeds 10 characters (length: 19) [SuperLongSentence]\n"
	if code != exitOK || stdout != want {
		t.Errorf("code = %d, stdout = %q, want %q", code, stdout, want)
	}
}

func TestRunListRules(t *testing.T) {
	code, stdout, _ := runTxtlint(t, "", "-list-rules")
	if code != exitOK {
		t.Errorf("exit code = %d", code)
	}
	for _, name := range []string{"TrailingSpaces", "SuperLongSentence", "MissingPunctuation",
		"NoConsecutiveBlankLines", "NoMixedLineEndings", "EOFNewline"} {
		if !strings.Contains(stdout, name) {
			t.Errorf("-list-rules output is missing %s:\n%s", name, stdout)
		}
	}
}

func TestRunFailures(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.txt")
	tests := []struct {
		name       string
		args       []string
		wantStderr string
	}{
		{"unknown rule", []string{"-disable", "Nope"}, `unknown rule "Nope"`},
		{"bad length", []string{"-max-sentence-length", "0"}, "must be at least 1"},
		{"unknown flag", []string{"-bogus"}, "flag provided but not defined"},
		{"missing file", []string{missing}, "missing.txt"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, stderr := runTxtlint(t, "", tt.args...)
			if code != exitFailed {
				t.Errorf("exit code = %d, want %d", code, exitFailed)
			}
			if !strings.Contains(stderr, tt.wantStderr) {
				t.Errorf("stderr = %q, want it to contain %q", stderr, tt.wantStderr)
			}
		})
	}
}

func TestRunHelp(t *testing.T) {
	code, _, stderr := runTxtlint(t, "", "-h")
	if code != exitOK || !strings.Contains(stderr, "Usage: txtlint") {
		t.Errorf("code = %d, stderr = %q", code, stderr)
	}
}
