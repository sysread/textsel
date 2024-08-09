package textsel

import "strings"

// Highlights the cursor position and selected text in the widget.
func (ts *TextSel) highlightCursor() {
	text := ts.text
	startRow, startCol, endRow, endCol := ts.GetSelectionRange()

	buf := strings.Builder{}
	sel := false
	row := 0
	col := 0

	// We only highlight the cursor if the widget has focus
	showCursor := ts.HasFocus()

	// Track the current format in the text, excluding the codes we use to
	// highlight the cursor and selection.
	formatCode := newFormatCode()

	for idx := 0; idx < len(text); idx++ {
		// Check for the start of a format code using the regex. Do so until we
		// get all of the ones at the current index.
		for formatRegex.MatchString(text[idx:]) {
			match := formatRegex.FindString(text[idx:])

			// If we are selecting, we skip the format code, because it may
			// interfere with the selection format. Otherwise, we write it to
			// the buffer.
			if !sel {
				// Write the format code unchanged to the buffer
				buf.WriteString(match)
			}

			// Adjust index to skip the matched format code
			idx += len(match)

			// Parse the format code into its components and save them
			// for later use.
			formatCode = formatCode.update(match)
		}

		// Now that we've dealt with any leading format tags, we can get the
		// current character.
		char := text[idx]

		// Mark the start of the selection
		isSelStartRow := row == startRow
		isSelStartCol := col == startCol

		// If the cursor moves up or down from a column > 0, but the current
		// line is empty, pretend the cursor is on the first column.
		if !isSelStartCol && char == '\n' && col == 0 && startCol > 0 {
			isSelStartCol = true
		}

		// If this is the beginning of the selection, mark it
		if ts.isSelecting && isSelStartRow && isSelStartCol {
			sel = true
			buf.WriteString(ts.selectionColor)
		}

		// Determine if the cursor is on the current character
		isCursorRow := row == ts.cursorRow
		isCursorCol := col == ts.cursorCol

		// If the cursor moves up or down from a column > 0, but the current
		// line is empty, pretend the cursor is on the first column.
		if !isCursorCol && char == '\n' && col == 0 && ts.cursorCol > 0 {
			isCursorCol = true
		}

		// Highlight the cursor position
		if showCursor && isCursorRow && isCursorCol {
			cursorStart := ts.cursorColor
			cursorEnd := formatCode.String()

			if sel {
				cursorStart = ts.cursorInSelectionColor

				if ts.cursorRow != endRow || ts.cursorCol != endCol {
					// If the cursor is NOT at the end of the selection, the
					// selection color should be used to "reset" the format.
					cursorEnd = ts.selectionColor
				} else {
					// If we ARE in the middle of a selection, we need to use
					// a slightly modified version of the current format string,
					// to remove the in-selection styles used to highlight the
					// cursor.
					cursorEnd = formatCode.update("[::-]").String()
				}
			}

			buf.WriteString(cursorStart)

			// If the cursor is on an empty line ("\n"), add a space to make it
			// visible.
			if char == '\n' {
				buf.WriteString(" \n")
			} else {
				buf.WriteString(string(char))
			}

			buf.WriteString(cursorEnd)
		} else if char == '\n' {
			// If the cursor is not on the current character, but it's a newline,
			// we need to add a space to make it visible.
			buf.WriteString(" \n")
		} else {
			buf.WriteString(string(char))
		}

		// Mark the end of the selection
		if sel && row == endRow && col == endCol {
			sel = false
			buf.WriteString(formatCode.String())
		}

		if char == '\n' {
			row++
			col = 0
		} else {
			col++
		}
	}

	ts.TextView.SetText(buf.String())
}
