package textsel

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
