// Package textsel provides a `tview.TextView` widget that supports selecting text with the keyboard.
package textsel

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var formatRegex = regexp.MustCompile(`^\[[-:a-zA-Z0-9]+\]`)
var formatRegexGlobal = regexp.MustCompile(`\[[-:a-zA-Z0-9]+\]`)

// TextSel is a `tview.TextView` widget that supports selecting text with the keyboard.
type TextSel struct {
	*tview.TextView

	text string

	// Cursor position
	cursorRow int
	cursorCol int

	// Selection state
	isSelecting       bool
	selectionStartRow int
	selectionStartCol int
	selectionEndRow   int
	selectionEndCol   int

	// Color codes for highlighting text
	defaultColor           string
	cursorColor            string
	selectionColor         string
	cursorInSelectionColor string

	// Callback for handling selected text
	selectFunc func(string)
}

// NewTextSel creates and returns a new TextSel instance.
//
// Example:
//
//	textSel := textsel.NewTextSel().SetText("Hello, World!")
func NewTextSel() *TextSel {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)

	ts := &TextSel{
		TextView:               textView,
		defaultColor:           fmt.Sprintf("[%s:%s:-]", colorToHex(tview.Styles.PrimaryTextColor), colorToHex(tview.Styles.PrimitiveBackgroundColor)),
		cursorColor:            fmt.Sprintf("[%s:%s:-]", colorToHex(tview.Styles.PrimitiveBackgroundColor), colorToHex(tview.Styles.PrimaryTextColor)),
		selectionColor:         fmt.Sprintf("[%s:%s:-]", colorToHex(tview.Styles.PrimitiveBackgroundColor), colorToHex(tview.Styles.SecondaryTextColor)),
		cursorInSelectionColor: fmt.Sprintf("[%s:%s:bu]", colorToHex(tview.Styles.PrimitiveBackgroundColor), colorToHex(tview.Styles.SecondaryTextColor)),
	}

	// Handle key events for moving the cursor and selecting text
	ts.SetInputCapture(ts.handleKeyEvents)

	ts.resetCursor()
	ts.resetSelection()

	return ts
}

// Helper function to convert tcell.Color to a hex string.
func colorToHex(color tcell.Color) string {
	return fmt.Sprintf("#%06X", color.Hex())
}

// Debugging function to write to `debug.log`.
func (ts *TextSel) debug(format string, args ...interface{}) {
	file, _ := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer file.Close()
	line := fmt.Sprintf(format, args...)
	file.WriteString(line + "\n")
}

// Debugging function to log the color codes being used for the cursor and
// selection highlighting.
func (ts *TextSel) debugColors() *TextSel {
	ts.debug("          defaultColor: %s", ts.defaultColor)
	ts.debug("           cursorColor: %s", ts.cursorColor)
	ts.debug("        selectionColor: %s", ts.selectionColor)
	ts.debug("cursorInSelectionColor: %s", ts.cursorInSelectionColor)
	return ts
}

// Debugging function to log the cursor position.
func (ts *TextSel) debugCursor() *TextSel {
	ts.debug("Cursor: (%d, %d)", ts.cursorRow, ts.cursorCol)
	return ts
}

// Debugging function to log the selection range.
func (ts *TextSel) debugSelection() *TextSel {
	if ts.isSelecting {
		startRow, startCol, endRow, endCol := ts.getSelectionRange()
		ts.debug("Selection range: (%d, %d) - (%d, %d)", startRow, startCol, endRow, endCol)
	} else {
		ts.debug("No selection")
	}
	return ts
}

// Resets the cursor position to the beginning of the text.
func (ts *TextSel) resetCursor() *TextSel {
	ts.cursorRow = 0
	ts.cursorCol = 0

	ts.highlightCursor()

	return ts
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

// GetText returns the text content of the TextSel widget. If `stripFormatting`
// is true, any format codes in the text will be removed.
//
// Example:
//
//	text := textSel.GetText(false)
func (ts *TextSel) GetText(stripFormatting bool) string {
	text := ts.text

	if stripFormatting {
		// Remove any format codes from the text
		text = formatRegexGlobal.ReplaceAllString(text, "")
	}

	return text
}

// SetText sets the text content of the TextSel widget, resetting the cursor
// position and selection state.
//
// Example:
//
//	textSel.SetText("New text content")
func (ts *TextSel) SetText(text string) *TextSel {
	//ts.text = text
	ts.TextView.SetText(text)
	ts.text = ts.TextView.GetText(false)

	ts.resetCursor()
	ts.resetSelection()

	return ts
}

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

// Helper function to get the selection range in the text.
func (ts *TextSel) getSelectionRange() (int, int, int, int) {
	startRow, startCol := ts.selectionStartRow, ts.selectionStartCol
	endRow, endCol := ts.selectionEndRow, ts.selectionEndCol

	if startRow > endRow || (startRow == endRow && startCol > endCol) {
		startRow, startCol, endRow, endCol = endRow, endCol, startRow, startCol
	}

	return startRow, startCol, endRow, endCol
}

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

// Moves the cursor up by one row.
func (ts *TextSel) moveUp() {
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
}

// Moves the cursor down by one row.
func (ts *TextSel) moveDown() {
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
}

// Moves the cursor left by one column.
func (ts *TextSel) moveLeft() {
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
}

// Moves the cursor right by one column.
func (ts *TextSel) moveRight() {
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
}

// Moves the cursor to the start of the current line.
func (ts *TextSel) moveToStartOfLine() {
	ts.cursorCol = 0

	if ts.isSelecting {
		ts.selectionEndCol = ts.cursorCol
	}

	ts.highlightCursor()
}

// Moves the cursor to the end of the current line.
func (ts *TextSel) moveToEndOfLine() {
	ts.cursorCol = len(ts.getCurrentLine()) - 1

	if ts.isSelecting {
		ts.selectionEndCol = ts.cursorCol
	}

	ts.highlightCursor()
}

// Starts the selection process.
func (ts *TextSel) startSelection() {
	ts.isSelecting = true
	ts.selectionStartRow = ts.cursorRow
	ts.selectionStartCol = ts.cursorCol
	ts.selectionEndRow = ts.cursorRow
	ts.selectionEndCol = ts.cursorCol
}

// Finishes the selection process and calls the selectFunc callback.
func (ts *TextSel) finishSelection() {
	if ts.selectFunc != nil {
		ts.selectFunc(ts.GetSelectedText())
	}

	ts.resetSelection()
}

// Handles key events for moving the cursor and selecting text.
func (ts *TextSel) handleKeyEvents(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyUp:
		ts.moveUp()
	case tcell.KeyDown:
		ts.moveDown()
	case tcell.KeyLeft:
		ts.moveLeft()
	case tcell.KeyRight:
		ts.moveRight()
	case tcell.KeyEnter:
		ts.finishSelection()
	case tcell.KeyRune:
		switch event.Rune() {
		case ' ':
			ts.startSelection()
		case 'k':
			ts.moveUp()
		case 'j':
			ts.moveDown()
		case 'h':
			ts.moveLeft()
		case 'l':
			ts.moveRight()
		case '^':
			ts.moveToStartOfLine()
		case '$':
			ts.moveToEndOfLine()
		}
	}

	return event
}

// Highlights the cursor position and selected text in the widget.
func (ts *TextSel) highlightCursor() {
	text := ts.text
	startRow, startCol, endRow, endCol := ts.getSelectionRange()

	buf := strings.Builder{}
	sel := false
	row := 0
	col := 0

	// We only highlight the cursor if the widget has focus
	showCursor := ts.TextView.HasFocus()

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
