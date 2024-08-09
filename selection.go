package textsel

import (
	"strings"
)

// SetSelectFunc sets the callback function that will be called when text is
// selected.
//
// Example:
//
//	textSel.SetSelectFunc(func(selectedText string) {
//		fmt.Println("Selected text:\n\n", selectedText)
//	})
func (ts *TextSel) SetSelectFunc(f func(string)) *TextSel {
	ts.selectFunc = f
	return ts
}

// GetSelectedText returns the currently selected text. If no text is selected,
// an empty string is returned.
//
// Example:
//
//	selectedText := textSel.GetSelectedText()
//	fmt.Println("Selected text:", selectedText)
func (ts *TextSel) GetSelectedText() string {
	text := ts.text
	startRow, startCol, endRow, endCol := ts.getSelectionRange()

	buf := strings.Builder{}
	sel := false
	row := 0
	col := 0

	for idx := 0; idx < len(text); idx++ {
		// Skip any format codes
		for formatRegex.MatchString(text[idx:]) {
			match := formatRegex.FindString(text[idx:])
			idx += len(match)
		}

		char := text[idx]

		if row == startRow && col == startCol {
			sel = true
		}

		if sel {
			buf.WriteString(string(char))
		}

		if row == endRow && col == endCol {
			sel = false
			break
		}

		if char == '\n' {
			row++
			col = 0
		} else {
			col++
		}
	}

	return buf.String()
}

// Resets the selection state.
func (ts *TextSel) resetSelection() *TextSel {
	ts.isSelecting = false

	ts.selectionStartRow = 0
	ts.selectionStartCol = 0
	ts.selectionEndRow = 0
	ts.selectionEndCol = 0

	ts.highlightCursor()

	return ts
}

func (ts *TextSel) getSelectionRange() (int, int, int, int) {
	startRow, startCol := ts.selectionStartRow, ts.selectionStartCol
	endRow, endCol := ts.selectionEndRow, ts.selectionEndCol

	if startRow > endRow || (startRow == endRow && startCol > endCol) {
		startRow, startCol, endRow, endCol = endRow, endCol, startRow, startCol
	}

	return startRow, startCol, endRow, endCol
}

// Starts the selection process.
func (ts *TextSel) startSelection() *TextSel {
	ts.isSelecting = true
	ts.selectionStartRow = ts.cursorRow
	ts.selectionStartCol = ts.cursorCol
	ts.selectionEndRow = ts.cursorRow
	ts.selectionEndCol = ts.cursorCol

	return ts
}

// Finishes the selection process and calls the selectFunc callback.
func (ts *TextSel) finishSelection() *TextSel {
	if ts.selectFunc != nil {
		ts.selectFunc(ts.GetSelectedText())
	}

	ts.resetSelection()

	return ts
}
