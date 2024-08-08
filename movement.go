package textsel

// Resets the cursor position to the beginning of the text.
func (ts *TextSel) resetCursor() *TextSel {
	ts.cursorRow = 0
	ts.cursorCol = 0

	ts.highlightCursor()

	return ts
}

// Moves the cursor up by one row.
func (ts *TextSel) moveUp() *TextSel {
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
func (ts *TextSel) moveDown() *TextSel {
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

// Moves the cursor left by one column.
func (ts *TextSel) moveLeft() *TextSel {
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

// Moves the cursor right by one column.
func (ts *TextSel) moveRight() *TextSel {
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
func (ts *TextSel) moveToStartOfLine() *TextSel {
	ts.cursorCol = 0

	if ts.isSelecting {
		ts.selectionEndCol = ts.cursorCol
	}

	ts.highlightCursor()

	return ts
}

// Moves the cursor to the end of the current line.
func (ts *TextSel) moveToEndOfLine() *TextSel {
	ts.cursorCol = len(ts.getCurrentLine()) - 1

	if ts.isSelecting {
		ts.selectionEndCol = ts.cursorCol
	}

	ts.highlightCursor()

	return ts
}
