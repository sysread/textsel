package textsel

// Resets the cursor position to the beginning of the text.
func (ts *TextSel) ResetCursor() *TextSel {
	ts.SetCursorPosition(0, 0)
	ts.ResetSelection()
	return ts
}

func (ts *TextSel) SetCursorPosition(row int, col int) *TextSel {
	ts.cursorRow = row
	ts.cursorCol = col

	if ts.isSelecting {
		ts.selectionEndRow = ts.cursorRow
		ts.selectionEndCol = ts.cursorCol
	}

	ts.highlightCursor()

	return ts
}

// Returns the current cursor position as row and column.
func (ts *TextSel) GetCursorPosition() (int, int) {
	return ts.cursorRow, ts.cursorCol
}

// Moves the cursor up by one row.
func (ts *TextSel) MoveUp() *TextSel {
	if ts.cursorRow > 0 {
		ts.cursorRow--

		currentLine := ts.getCurrentLine()

		if ts.cursorCol > len(currentLine) {
			ts.cursorCol = len(currentLine) - 1
			if ts.cursorCol < 0 {
				ts.cursorCol = 0
			}
		}

		if ts.isSelecting {
			ts.selectionEndRow = ts.cursorRow
			ts.selectionEndCol = ts.cursorCol
		}
	}

	ts.highlightCursor()

	return ts
}

// Moves the cursor down by one row.
func (ts *TextSel) MoveDown() *TextSel {
	if ts.cursorRow < ts.lastRow() {
		ts.cursorRow++

		currentLine := ts.getCurrentLine()

		if ts.cursorCol > len(currentLine) {
			ts.cursorCol = len(currentLine) - 1

			if ts.cursorCol < 0 {
				ts.cursorCol = 0
			}
		}

		if ts.isSelecting {
			ts.selectionEndRow = ts.cursorRow
			ts.selectionEndCol = ts.cursorCol
		}
	}

	ts.highlightCursor()

	return ts
}

// Moves the cursor left by one column, wrapping to the previous row if necessary.
func (ts *TextSel) MoveLeft() *TextSel {
	if ts.cursorCol > 0 {
		ts.cursorCol--
	} else if ts.cursorRow > 0 {
		ts.cursorRow--
		ts.cursorCol = len(ts.getCurrentLine()) - 1 // Adjust to the last valid column in the previous row
	}

	if ts.isSelecting {
		ts.selectionEndRow = ts.cursorRow
		ts.selectionEndCol = ts.cursorCol
	}

	ts.highlightCursor()

	return ts
}

// Moves the cursor right by one column, wrapping to the next row if necessary.
func (ts *TextSel) MoveRight() *TextSel {
	if ts.cursorCol < len(ts.getCurrentLine())-1 {
		ts.cursorCol++
	} else if ts.cursorRow < ts.lastRow() {
		ts.cursorRow++
		ts.cursorCol = 0
	}

	if ts.isSelecting {
		ts.selectionEndRow = ts.cursorRow
		ts.selectionEndCol = ts.cursorCol
	}

	ts.highlightCursor()

	return ts
}

// Moves the cursor to the start of the current line.
func (ts *TextSel) MoveToStartOfLine() *TextSel {
	ts.cursorCol = 0

	if ts.isSelecting {
		ts.selectionEndCol = ts.cursorCol
	}

	ts.highlightCursor()

	return ts
}

// Moves the cursor to the end of the current line.
func (ts *TextSel) MoveToEndOfLine() *TextSel {
	ts.cursorCol = len(ts.getCurrentLine()) - 1

	if ts.isSelecting {
		ts.selectionEndCol = ts.cursorCol
	}

	ts.highlightCursor()

	return ts
}

// Moves the cursor to the first line of the text. The current column is
// preserved if it is within the bounds of the first line. Otherwise, the
// cursor is placed at the end of the first line.
func (ts *TextSel) MoveToFirstLine() *TextSel {
	ts.cursorRow = 0

	if ts.isSelecting {
		ts.selectionEndRow = ts.cursorRow
	}

	if ts.cursorCol > len(ts.getCurrentLine()) {
		ts.MoveToEndOfLine()
	}

	ts.highlightCursor()

	return ts
}

// Moves the cursor to the last line of the text. The current column is
// preserved if it is within the bounds of the last line. Otherwise, the
// cursor is placed at the end of the last line.
func (ts *TextSel) MoveToLastLine() *TextSel {
	ts.cursorRow = ts.lastRow()

	if ts.isSelecting {
		ts.selectionEndRow = ts.cursorRow
	}

	if ts.cursorCol > len(ts.getCurrentLine()) {
		ts.MoveToEndOfLine()
	}

	ts.highlightCursor()

	return ts
}
