#!/bin/bash

# Define the filename
FILE="test.txt"

# --- Content for the test file ---
# This content is carefully crafted to trigger line-based rules:
# 1. Trailing spaces on the first line.
# 2. Super long sentence on the second line.
# 3. Two consecutive blank lines (triggers NoConsecutiveBlankLines on line 6, the second blank line).
# 4. Trailing space AND missing punctuation on line 7.
# 5. Lines 8 and 9 will be used to demonstrate mixed line endings.

# We'll use printf to write each line, allowing us to control newlines precisely.
# By default, we use '\n' (LF) for most lines as it's common on macOS.

printf "This line has trailing spaces.     \n" >"$FILE"                                                                                                                                                                                                                                                    # Line 1: TrailingSpaces (LF)
printf "This sentence is extremely long and will definitely exceed the 120-character limit that we've set for the SuperLongSentence rule, proving that it can detect overly verbose prose that might be difficult for readers to process quickly and efficiently within a single visual scan.\n" >>"$FILE" # Line 2: SuperLongSentence (LF)
printf "\n" >>"$FILE"                                                                                                                                                                                                                                                                                      # Line 3: Blank line (LF)
printf "\n" >>"$FILE"                                                                                                                                                                                                                                                                                      # Line 4: Blank line (LF)
printf "This line is followed by two blank lines, triggering the NoConsecutiveBlankLines rule.\n" >>"$FILE"                                                                                                                                                                                                # Line 5: Followed by blank line (LF)
printf "\n" >>"$FILE"                                                                                                                                                                                                                                                                                      # Line 6: Second consecutive blank line (LF) -> triggers NoConsecutiveBlankLines here

# Note: The problem description implies "This sentence might be missing punctuation at the end " (with a trailing space) was line 7
printf "This sentence might be missing punctuation at the end \n" >>"$FILE" # Line 7: TrailingSpaces, MissingPunctuation (LF for now)

printf "This is another line that will have a CRLF ending.\n" >>"$FILE" # Line 8: To be converted to CRLF (LF for now)
printf "This is a line with an LF ending." >>"$FILE"                    # Line 9: No trailing newline for EOFNewline rule (LF)

# --- Manipulations for File-Level Rules ---

# 1. Trigger "EOFNewline" (Missing):
#    The last printf for Line 9 above does NOT add a newline.
#    This ensures the file ends without a trailing newline character.

# 2. Trigger "NoMixedLineEndings":
#    We will convert Line 8's ending from LF to CRLF (\r\n).
#    macOS's `sed` (BSD sed) requires a literal carriage return.
#    We use `printf '\r'` to insert the actual CR character.
#    Line 8 has 7 lines of content before it, so it's the 8th line overall.
LINE_TO_CONVERT_TO_CRLF=8
# -i '' is for in-place editing on macOS's sed (empty string for backup extension)
# s/$/\r/ appends a carriage return before the existing newline of the line
sed -i '' "${LINE_TO_CONVERT_TO_CRLF}s/$/$(printf '\r')/" "$FILE"

echo "Successfully generated '$FILE' to trigger all rules."
echo ""
echo "Expected linting issues for '$FILE':"
echo "  - [TrailingSpaces] on Line 1"
echo "  - [SuperLongSentence] on Line 2"
echo "  - [NoConsecutiveBlankLines] on Line 6 (the second blank line)"
echo "  - [TrailingSpaces] on Line 7"
echo "  - [MissingPunctuation] on Line 7"
echo "  - [NoMixedLineEndings] (file-level, Line 0)"
echo "  - [EOFNewline] on Line 9 (last line, as it lacks a final newline)"

echo ""
echo "You can now run: go run main.go $FILE"
