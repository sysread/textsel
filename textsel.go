// Package textsel provides a `tview.TextView` widget that supports selecting text with the keyboard.
package textsel

import (
	"fmt"
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
	if !ts.isSelecting {
		return ""
	}

	startRow, startCol, endRow, endCol := ts.getSelectionRange()

	lines := strings.Split(ts.text, "\n")
	selectedLines := []string{}

	for row := startRow; row <= endRow; row++ {
		line := lines[row]

		if row == startRow && row == endRow {
			selectedLines = append(selectedLines, line[startCol:endCol+1])
		} else if row == startRow {
			selectedLines = append(selectedLines, line[startCol:])
		} else if row == endRow {
			selectedLines = append(selectedLines, line[:endCol+1])
		} else {
			selectedLines = append(selectedLines, line)
		}
	}

	return strings.Join(selectedLines, "\n")
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
		// Call the selectFunc callback with the selected text
		if ts.selectFunc != nil {
			ts.selectFunc(ts.GetSelectedText())
		}

		// Clear the selection
		ts.isSelecting = false
		ts.selectionEndRow = ts.cursorRow
		ts.selectionEndCol = ts.cursorCol
	case tcell.KeyRune:
		switch event.Rune() {
		case ' ':
			// Start selection
			ts.isSelecting = true
			ts.selectionStartRow = ts.cursorRow
			ts.selectionStartCol = ts.cursorCol
			ts.selectionEndRow = ts.cursorRow
			ts.selectionEndCol = ts.cursorCol
		case 'k':
			ts.moveUp()
		case 'j':
			ts.moveDown()
		case 'h':
			ts.moveLeft()
		case 'l':
			ts.moveRight()
		case '^':
			ts.cursorCol = 0
			if ts.isSelecting {
				ts.selectionEndCol = ts.cursorCol
			}
		case '$':
			ts.cursorCol = len(ts.getCurrentLine()) - 1
			if ts.isSelecting {
				ts.selectionEndCol = ts.cursorCol
			}
		}
	}

	ts.highlightCursor()
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
	lines := strings.Split(ts.text, "\n")
	originalLines := strings.Split(ts.text, "\n")

	if ts.cursorRow >= len(lines) {
		return
	}

	// Reset the entire text to its original form
	for i := range lines {
		lines[i] = originalLines[i]
	}

	// Apply selection highlighting
	if ts.isSelecting {
		startRow, startCol, endRow, endCol := ts.getSelectionRange()

		for row := startRow; row <= endRow; row++ {
			if row >= len(lines) {
				continue
			}

			line := lines[row]

			// Handle empty lines
			if len(line) == 0 {
				lines[row] = fmt.Sprintf("%s %s", ts.selectionColor, ts.defaultColor)
				continue
			}

			if row == startRow && row == endRow {
				lines[row] = fmt.Sprintf("%s%s%s%s%s", line[:startCol], ts.selectionColor, line[startCol:endCol+1], ts.defaultColor, line[endCol+1:])
			} else if row == startRow {
				lines[row] = fmt.Sprintf("%s%s%s%s", line[:startCol], ts.selectionColor, line[startCol:], ts.defaultColor)
			} else if row == endRow {
				lines[row] = fmt.Sprintf("%s%s%s%s", ts.selectionColor, line[:endCol+1], ts.defaultColor, line[endCol+1:])
			} else {
				lines[row] = fmt.Sprintf("%s%s%s", ts.selectionColor, line, ts.defaultColor)
			}
		}
	}

	// Highlight the cursor position
	line := lines[ts.cursorRow]
	cursorCol := ts.cursorCol

	// Adjust cursorCol if selecting by adding the length of the selection color
	if ts.isSelecting {
		cursorCol += len(ts.selectionColor)
	}

	// Adjust cursorCol if it is past the end of the line
	if cursorCol > len(line) {
		cursorCol = len(line)
	}

	// Alter the cursor color if the cursor is within the selection
	cursorColor := ts.cursorColor
	resetColor := ts.defaultColor
	if ts.cursorLocationIsWithinSelection() {
		cursorColor = ts.cursorInSelectionColor
		resetColor = ts.selectionColor
	}

	if len(line) == 0 {
		// If the line is empty, set the cursor at the start of the line
		lines[ts.cursorRow] = fmt.Sprintf("%s %s", cursorColor, resetColor)
	} else if cursorCol >= len(line) {
		// Ensure the cursor stays within the line boundaries
		ts.cursorCol = len(line) - 1
		line = lines[ts.cursorRow]
		highlightedText := fmt.Sprintf("%s%s%c%s%s", line[:cursorCol], cursorColor, line[cursorCol], resetColor, "")
		lines[ts.cursorRow] = highlightedText
	} else {
		// Highlight the character at the cursor position
		highlightedText := fmt.Sprintf("%s%s%c%s%s", line[:cursorCol], cursorColor, line[cursorCol], resetColor, line[cursorCol+1:])
		lines[ts.cursorRow] = highlightedText
	}

	ts.TextView.SetText(strings.Join(lines, "\n"))
}

// Retrieves the current line the cursor is on.
func (ts *TextSel) getCurrentLine() string {
	lines := strings.Split(ts.text, "\n")

	if ts.cursorRow >= len(lines) {
		return ""
	}

	return lines[ts.cursorRow]
}
