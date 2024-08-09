package textsel

import (
	"strings"
)

// Retrieves the current line the cursor is on.
func (ts *TextSel) getCurrentLine() string {
	text := ts.text
	buf := strings.Builder{}
	row := 0

	for idx := 0; idx < len(text); idx++ {
		// Skip any format codes
		for formatRegex.MatchString(text[idx:]) {
			match := formatRegex.FindString(text[idx:])
			idx += len(match)
		}

		char := text[idx]

		if row == ts.cursorRow {
			buf.WriteString(string(char))
		}

		if char == '\n' {
			row++

			if row > ts.cursorRow {
				break
			}
		}
	}

	return buf.String()
}

// Returns the row index (zero-based) the last line in the text.
func (ts *TextSel) lastRow() int {
	text := ts.text
	lastIndex := len(text) - 1
	count := 0

	for idx := 0; idx < len(text); idx++ {
		if text[idx] == '\n' && idx != lastIndex {
			count++
		}
	}

	return count
}
