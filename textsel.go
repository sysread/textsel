// Package textsel provides a `tview.TextView` widget that supports selecting text with the keyboard.
package textsel

import (
	"fmt"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TextSel is a `tview.TextView` widget that supports selecting text with the keyboard.
type TextSel struct {
	*tview.TextView

	text string

	// Cursor position
	cursorRow, cursorCol int

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
//	textSel := textsel.NewTextSel().
//	    SetText("Hello, World!")
func NewTextSel() *TextSel {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true)

	ts := &TextSel{
		TextView:               textView,
		cursorRow:              0,
		cursorCol:              0,
		defaultColor:           fmt.Sprintf("[%s:%s:-]", colorToHex(tview.Styles.PrimaryTextColor), colorToHex(tview.Styles.PrimitiveBackgroundColor)),
		cursorColor:            fmt.Sprintf("[%s:%s:-]", colorToHex(tview.Styles.PrimitiveBackgroundColor), colorToHex(tview.Styles.PrimaryTextColor)),
		selectionColor:         fmt.Sprintf("[%s:%s:-]", colorToHex(tview.Styles.PrimitiveBackgroundColor), colorToHex(tview.Styles.SecondaryTextColor)),
		cursorInSelectionColor: fmt.Sprintf("[%s:%s:bu]", colorToHex(tview.Styles.PrimitiveBackgroundColor), colorToHex(tview.Styles.SecondaryTextColor)),
	}

	// Handle key events for moving the cursor and selecting text
	ts.SetInputCapture(ts.handleKeyEvents)

	return ts
}

func (ts *TextSel) debug(format string, args ...interface{}) {
	file, _ := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer file.Close()
	line := fmt.Sprintf(format, args...)
	file.WriteString(line + "\n")
}

func (ts *TextSel) debugColors() *TextSel {
	ts.debug("          defaultColor: %s", ts.defaultColor)
	ts.debug("           cursorColor: %s", ts.cursorColor)
	ts.debug("        selectionColor: %s", ts.selectionColor)
	ts.debug("cursorInSelectionColor: %s", ts.cursorInSelectionColor)
	return ts
}

func (ts *TextSel) debugCursor() *TextSel {
	ts.debug("Cursor: (%d, %d)", ts.cursorRow, ts.cursorCol)
	return ts
}

func (ts *TextSel) debugSelection() *TextSel {
	if ts.isSelecting {
		ts.debug("Selection range: (%d, %d) - (%d, %d)", ts.selectionStartRow, ts.selectionStartCol, ts.selectionEndRow, ts.selectionEndCol)
	} else {
		ts.debug("No selection")
	}
	return ts
}

// SetText sets the text content of the TextSel widget, retaining the cursor
// and selection positions. Note that if the text is significantly different
// from the previous text, the cursor and selection positions may not be where
// you expect them to be. If the text is shortened, making the cursor or
// selection positions invalid, they will be adjusted to the end of the text.
//
// Example:
//
//	textSel.SetText("New text content")
func (ts *TextSel) SetText(text string) *TextSel {
	ts.text = text
	ts.TextView.SetText(text)
	ts.highlightCursor()

	lines := strings.Split(text, "\n")

	if ts.cursorRow >= len(lines) {
		ts.cursorRow = len(lines) - 1
	}

	if ts.selectionStartRow >= len(lines) {
		ts.selectionStartRow = len(lines) - 1
	}

	if ts.selectionEndRow >= len(lines) {
		ts.selectionEndRow = len(lines) - 1
	}

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
		char := text[idx]

		if char == '\n' {
			row++
			col = 0
		}

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
	}

	return buf.String()
}

// Helper function to convert tcell.Color to a hex string.
func colorToHex(color tcell.Color) string {
	return fmt.Sprintf("#%06X", color.Hex())
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

// Moves the cursor up by one row.
func (ts *TextSel) moveUp() {
	if ts.cursorRow > 0 {
		ts.cursorRow--
		if ts.cursorCol > len(ts.getCurrentLine()) {
			ts.cursorCol = len(ts.getCurrentLine())
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
	lines := strings.Split(ts.text, "\n")
	if ts.cursorRow < len(lines)-1 {
		ts.cursorRow++
		if ts.cursorCol > len(ts.getCurrentLine()) {
			ts.cursorCol = len(ts.getCurrentLine()) - 1
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
	} else if ts.cursorRow < len(strings.Split(ts.text, "\n"))-1 {
		ts.cursorRow++
		ts.cursorCol = 0
	}

	if ts.isSelecting {
		ts.selectionEndRow = ts.cursorRow
		ts.selectionEndCol = ts.cursorCol
	}

	ts.highlightCursor()
}

func (ts *TextSel) moveToStartOfLine() {
	ts.cursorCol = 0
	if ts.isSelecting {
		ts.selectionEndCol = ts.cursorCol
	}

	ts.highlightCursor()
}

func (ts *TextSel) moveToEndOfLine() {
	ts.cursorCol = len(ts.getCurrentLine()) - 1
	if ts.isSelecting {
		ts.selectionEndCol = ts.cursorCol
	}

	ts.highlightCursor()
}

func (ts *TextSel) startSelection() {
	ts.isSelecting = true
	ts.selectionStartRow = ts.cursorRow
	ts.selectionStartCol = ts.cursorCol
	ts.selectionEndRow = ts.cursorRow
	ts.selectionEndCol = ts.cursorCol
}

func (ts *TextSel) finishSelection() {
	// Call the selectFunc callback with the selected text
	if ts.selectFunc != nil {
		ts.selectFunc(ts.GetSelectedText())
	}

	ts.isSelecting = false
	ts.selectionEndRow = ts.cursorRow
	ts.selectionEndCol = ts.cursorCol
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

// Determines if the cursor location is within the selection.
func (ts *TextSel) cursorLocationIsWithinSelection() bool {
	if !ts.isSelecting {
		return false
	}

	startRow, startCol, endRow, endCol := ts.getSelectionRange()

	start := [2]int{startRow, startCol}
	end := [2]int{endRow, endCol}
	cursor := [2]int{ts.cursorRow, ts.cursorCol}

	lines := strings.Split(ts.text, "\n")

	return isCursorWithinRange(cursor, start, end, lines)
}

// Converts a cursor position to an absolute position in the text.
func toAbsolutePosition(point [2]int, lines []string) int {
	row, col := point[0], point[1]
	absolutePos := 0
	for i := 0; i < row; i++ {
		absolutePos += len(lines[i]) + 1 // Add 1 for the newline character
	}
	absolutePos += col
	return absolutePos
}

// Determines if the cursor is within the specified range in the text.
func isCursorWithinRange(cursor, start, end [2]int, lines []string) bool {
	cursorAbs := toAbsolutePosition(cursor, lines)
	startAbs := toAbsolutePosition(start, lines)
	endAbs := toAbsolutePosition(end, lines)

	return startAbs <= cursorAbs && cursorAbs <= endAbs
}

// Highlights the cursor position and selected text in the widget.
func (ts *TextSel) highlightCursor() {
	text := ts.text
	startRow, startCol, endRow, endCol := ts.getSelectionRange()

	buf := strings.Builder{}
	sel := false
	row := 0
	col := 0

	for idx := 0; idx < len(text); idx++ {
		char := text[idx]

		// Mark the start of the selection
		if ts.isSelecting && row == startRow && col == startCol {
			sel = true
			buf.WriteString(ts.selectionColor)
		}

		// Highlight the cursor position
		if row == ts.cursorRow && col == ts.cursorCol {
			if sel {
				buf.WriteString(ts.cursorInSelectionColor)
			} else {
				buf.WriteString(ts.cursorColor)
			}

			// If the cursor is on an empty line ("\n"), add a space to make it
			// visible.
			if char == '\n' {
				buf.WriteString(" \n")
			} else {
				buf.WriteString(string(char))
			}

			if sel {
				buf.WriteString(ts.selectionColor)
			} else {
				buf.WriteString(ts.defaultColor)
			}
		} else {
			buf.WriteString(string(char))
		}

		// Mark the end of the selection
		if row == endRow && col == endCol {
			sel = false
			buf.WriteString(ts.defaultColor)
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

// Retrieves the current line the cursor is on.
func (ts *TextSel) getCurrentLine() string {
	lines := strings.Split(ts.text, "\n")

	if ts.cursorRow >= len(lines) {
		return ""
	}

	return lines[ts.cursorRow]
}
