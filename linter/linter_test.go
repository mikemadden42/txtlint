package linter

import (
	"reflect"
	"testing"
)

func TestNewDocument(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []Line
	}{
		{"empty", "", nil},
		{"no final newline", "a", []Line{{1, "a", ""}}},
		{"LF", "a\nb\n", []Line{{1, "a", "\n"}, {2, "b", "\n"}}},
		{"CRLF", "a\r\nb\r\n", []Line{{1, "a", "\r\n"}, {2, "b", "\r\n"}}},
		{"CR", "a\rb", []Line{{1, "a", "\r"}, {2, "b", ""}}},
		{"mixed", "a\r\nb\nc\r", []Line{{1, "a", "\r\n"}, {2, "b", "\n"}, {3, "c", "\r"}}},
		{"LF CR is two endings", "a\n\rb", []Line{{1, "a", "\n"}, {2, "", "\r"}, {3, "b", ""}}},
		{"blank lines", "\n\n", []Line{{1, "", "\n"}, {2, "", "\n"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDocument([]byte(tt.input)).Lines
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDocument(%q).Lines = %#v, want %#v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIssueFormat(t *testing.T) {
	tests := []struct {
		issue Issue
		file  string
		want  string
	}{
		{Issue{Rule: "R", Line: 3, Column: 7, Message: "m"}, "f.txt", "f.txt:3:7: error: m [R]"},
		{Issue{Rule: "R", Severity: Warning, Line: 3, Message: "m"}, "f.txt", "f.txt:3: warning: m [R]"},
		{Issue{Rule: "R", Message: "m"}, "f.txt", "f.txt: error: m [R]"},
		{Issue{Rule: "R", Line: 2, Column: 1, Message: "m"}, "", "2:1: error: m [R]"},
		{Issue{Rule: "R", Message: "m"}, "", "error: m [R]"},
	}
	for _, tt := range tests {
		if got := tt.issue.Format(tt.file); got != tt.want {
			t.Errorf("Format(%q) = %q, want %q", tt.file, got, tt.want)
		}
	}
}

type fixedRule []Issue

func (r fixedRule) Name() string            { return "Fixed" }
func (r fixedRule) Description() string     { return "" }
func (r fixedRule) Check(*Document) []Issue { return r }

func TestLintSortsByPosition(t *testing.T) {
	rules := []Rule{
		fixedRule{{Line: 2, Column: 5, Message: "a"}, {Line: 0, Message: "b"}},
		fixedRule{{Line: 2, Column: 1, Message: "c"}, {Line: 1, Message: "d"}, {Line: 2, Column: 1, Message: "e"}},
	}
	var got []string
	for _, issue := range Lint(nil, rules) {
		got = append(got, issue.Message)
	}
	want := []string{"b", "d", "c", "e", "a"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Lint order = %v, want %v", got, want)
	}
}
