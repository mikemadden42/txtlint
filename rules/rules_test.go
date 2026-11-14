package rules

import (
	"reflect"
	"strings"
	"testing"

	"github.com/mikemadden42/txtlint/linter"
)

// pos is the position of an expected issue.
type pos struct{ line, col int }

type ruleTest struct {
	name  string
	input string
	want  []pos
}

func runRuleTests(t *testing.T, rule linter.Rule, tests []ruleTest) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []pos
			for _, issue := range rule.Check(linter.NewDocument([]byte(tt.input))) {
				if issue.Rule != rule.Name() {
					t.Errorf("issue rule = %q, want %q", issue.Rule, rule.Name())
				}
				got = append(got, pos{issue.Line, issue.Column})
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s on %q: got %v, want %v", rule.Name(), tt.input, got, tt.want)
			}
		})
	}
}

func TestTrailingSpaces(t *testing.T) {
	runRuleTests(t, NewTrailingSpacesRule(), []ruleTest{
		{"clean", "clean.\n", nil},
		{"spaces", "trailing.   \n", []pos{{1, 10}}},
		{"tab", "tab\t\n", []pos{{1, 4}}},
		{"whitespace-only line", "ok\n  \n", []pos{{2, 1}}},
		{"CRLF is not whitespace", "x\r\n", nil},
		{"space before CRLF", "x \r\n", []pos{{1, 2}}},
		{"non-ASCII last char", "Name: Å\n", nil},
		{"column counts runes", "Café \n", []pos{{1, 5}}},
	})
}

func TestMissingPunctuation(t *testing.T) {
	runRuleTests(t, NewMissingPunctuationRule(), []ruleTest{
		{"sentence", "This is a sentence.\n", nil},
		{"missing", "This is missing\n", []pos{{1, 15}}},
		{"indented", "    indented line no stop\n", []pos{{1, 25}}},
		{"trailing space", "trailing space here \n", []pos{{1, 19}}},
		{"non-ASCII last char", "Das ist ein Café\n", []pos{{1, 16}}},
		{"single word", "Heading\n", nil},
		{"indented single word", "  Heading\n", nil},
		{"colon", "Ends with colon:\n", nil},
		{"blank", "\n", nil},
	})
}

func TestSuperLongSentence(t *testing.T) {
	runRuleTests(t, NewSuperLongSentenceRule(20), []ruleTest{
		{"short", "Short one. Another short.\n", nil},
		{"long", "This sentence is definitely longer than twenty.\n", []pos{{1, 0}}},
		{"second sentence long", "Short.\nThis one is much too long.\n", []pos{{2, 0}}},
		{"wraps across lines", "This sentence wraps\nacross two lines here.\n", []pos{{1, 0}}},
		{"blank line ends sentence", "Short words here\n\nmore short words.\n", nil},
		{"closing quote", `He said "stop." Then left.` + "\n", nil},
		{"runes not bytes", "Ünïcödé wörds ïn sé.\n", nil},
		{"no final punctuation", "a b c d e f g h i j k l m n o p", []pos{{1, 0}}},
	})

	t.Run("decimal does not split", func(t *testing.T) {
		issues := NewSuperLongSentenceRule(10).Check(linter.NewDocument([]byte("It costs 3.50 now.\n")))
		if len(issues) != 1 || !strings.Contains(issues[0].Message, "length: 18") {
			t.Errorf("got %v, want one issue with length 18", issues)
		}
	})
}

func TestNoConsecutiveBlankLines(t *testing.T) {
	runRuleTests(t, NewNoConsecutiveBlankLinesRule(), []ruleTest{
		{"single blank", "a\n\nb\n", nil},
		{"double blank", "a\n\n\nb\n", []pos{{3, 0}}},
		{"whitespace-only counts as blank", "a\n\n  \n\nb\n", []pos{{3, 0}, {4, 0}}},
		{"CRLF", "a\r\n\r\n\r\nb\r\n", []pos{{3, 0}}},
	})
}

func TestNoMixedLineEndings(t *testing.T) {
	runRuleTests(t, NewNoMixedLineEndingsRule(), []ruleTest{
		{"LF only", "a\nb\n", nil},
		{"CRLF only", "a\r\nb\r\n", nil},
		{"no endings", "a", nil},
		{"mostly CRLF", "a\r\nb\nc\r\n", []pos{{2, 0}}},
		{"mostly LF", "a\nb\r\nc\n", []pos{{2, 0}}},
		{"tie goes to first", "a\nb\r\n", []pos{{2, 0}}},
		{"CR only", "a\rb\r", []pos{{1, 0}}},
		{"CR mixed with LF", "a\nb\rc\n", []pos{{2, 0}, {2, 0}}},
		{"CR mixed with CRLF", "a\r\nb\rc\r\n", []pos{{2, 0}, {2, 0}}},
		{"CR is never expected", "a\rb\rc\r\n", []pos{{1, 0}, {1, 0}}},
	})

	t.Run("message", func(t *testing.T) {
		issues := NewNoMixedLineEndingsRule().Check(linter.NewDocument([]byte("a\nb\nc\r\n")))
		want := "Mixed line endings: 2 LF, 1 CRLF (expected LF)"
		if len(issues) != 1 || issues[0].Message != want {
			t.Errorf("got %v, want one issue %q", issues, want)
		}
	})
}

func TestEOFNewline(t *testing.T) {
	runRuleTests(t, NewEOFNewlineRule(), []ruleTest{
		{"empty", "", nil},
		{"LF", "a\n", nil},
		{"CRLF", "a\r\n", nil},
		{"missing", "a", []pos{{1, 2}}},
		{"missing on later line", "a\nbé", []pos{{2, 3}}},
		{"extra LF", "a\n\n", []pos{{2, 0}}},
		{"extra CRLF", "a\r\n\r\n", []pos{{2, 0}}},
		{"several extra", "a\n\n\n", []pos{{2, 0}}},
		{"only newlines", "\n", nil},
	})
}

func TestAllRulesHaveUniqueNames(t *testing.T) {
	seen := map[string]bool{}
	for _, rule := range All(120) {
		if seen[rule.Name()] {
			t.Errorf("duplicate rule name %q", rule.Name())
		}
		seen[rule.Name()] = true
		if rule.Description() == "" {
			t.Errorf("rule %q has no description", rule.Name())
		}
	}
}
