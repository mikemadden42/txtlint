# txtlint

A linter for plain-text files. It checks for formatting problems (trailing
whitespace, line endings, blank lines, final newline) and a couple of style
heuristics (long sentences, missing punctuation).

## Install

```sh
go install github.com/mikemadden42/txtlint@latest
```

## Usage

```sh
txtlint [flags] [file ...]
```

With no file arguments (or `-`), txtlint reads standard input. Flags must come
before file names.

| Flag | Default | Description |
| --- | --- | --- |
| `-disable` | | Comma-separated, case-insensitive list of rules to skip |
| `-max-sentence-length` | `120` | Maximum sentence length, in characters |
| `-strict` | `false` | Exit with status 1 on warnings as well as errors |
| `-list-rules` | | List the available rules and exit |

Issues are printed one per line in a format editors and `grep` understand:

```text
notes.txt:4:34: warning: Line might be missing sentence-ending punctuation (.!?) [MissingPunctuation]
notes.txt:5:27: error: File does not end with a newline [EOFNewline]
```

### Exit status

| Status | Meaning |
| --- | --- |
| 0 | No errors (warnings may have been printed) |
| 1 | At least one error, or any warning with `-strict` |
| 2 | Bad usage, or a file could not be read |

## Rules

| Rule | Severity | Checks |
| --- | --- | --- |
| `TrailingSpaces` | error | Lines must not end with whitespace |
| `NoConsecutiveBlankLines` | error | No more than one blank line in a row |
| `NoMixedLineEndings` | error | One line ending style (LF or CRLF) throughout; no CR endings |
| `EOFNewline` | error | Files end with exactly one newline |
| `SuperLongSentence` | warning | Sentences, which may wrap across lines, stay under the length limit |
| `MissingPunctuation` | warning | Lines of more than one word end with `.`, `!` or `?` |

## Development

```sh
make test    # go test ./...
make check   # golangci-lint run
make build   # goreleaser snapshot build
```
