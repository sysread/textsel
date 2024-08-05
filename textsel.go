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
		startRow, startCol, endRow, endCol := ts.getSelectionRange()
		ts.debug("Selection range: (%d, %d) - (%d, %d)", startRow, startCol, endRow, endCol)
	} else {
		ts.debug("No selection")
	}
	return ts
}

func (ts *TextSel) resetCursor() *TextSel {
	ts.cursorRow = 0
	ts.cursorCol = 0

	ts.highlightCursor()

	return ts
}

func (ts *TextSel) resetSelection() *TextSel {
	ts.isSelecting = false

	ts.selectionStartRow = 0
	ts.selectionStartCol = 0
	ts.selectionEndRow = 0
	ts.selectionEndCol = 0

	ts.highlightCursor()

	return ts
}

// SetText sets the text content of the TextSel widget, resetting the cursor
// position and selection state.
//
// Example:
//
//	textSel.SetText("New text content")
func (ts *TextSel) SetText(text string) *TextSel {
	ts.text = text
	ts.TextView.SetText(text)

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
	lastRow := strings.Count(ts.text, "\n")

	if ts.cursorRow < lastRow {
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
	lastRow := strings.Count(ts.text, "\n")

	if ts.cursorCol < len(ts.getCurrentLine())-1 {
		ts.cursorCol++
	} else if ts.cursorRow < lastRow {
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

	for idx := 0; idx < len(text); idx++ {
		char := text[idx]

		// Mark the start of the selection
		if ts.isSelecting && row == startRow && col == startCol {
			sel = true
			buf.WriteString(ts.selectionColor)
		}

		// Highlight the cursor position
		if row == ts.cursorRow && col == ts.cursorCol {
			cursorStart := ts.cursorColor
			cursorEnd := ts.defaultColor
			if sel {
				cursorStart = ts.cursorInSelectionColor
				cursorEnd = ts.selectionColor
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
			buf.WriteString(" \n")
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
