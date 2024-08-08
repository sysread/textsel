// Package textsel provides a `tview.TextView` widget that supports selecting text with the keyboard.
package textsel

import (
	"fmt"
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
		defaultColor:           fmt.Sprintf("[%s:%s:-]", tview.Styles.PrimaryTextColor, tview.Styles.PrimitiveBackgroundColor),
		cursorColor:            fmt.Sprintf("[%s:%s:-]", tview.Styles.PrimitiveBackgroundColor, tview.Styles.PrimaryTextColor),
		selectionColor:         fmt.Sprintf("[%s:%s:-]", tview.Styles.PrimitiveBackgroundColor, tview.Styles.SecondaryTextColor),
		cursorInSelectionColor: fmt.Sprintf("[%s:%s:bu]", tview.Styles.PrimitiveBackgroundColor, tview.Styles.SecondaryTextColor),
	}

	// Handle key events for moving the cursor and selecting text
	ts.SetInputCapture(ts.handleKeyEvents)

	// Ensure that we redraw when we are focused or blurred to update whether
	// the cursor is visible or not.
	ts.SetFocusFunc(func() {
		go ts.highlightCursor()
	})

	ts.SetBlurFunc(func() {
		go ts.highlightCursor()
	})

	ts.resetCursor()
	ts.resetSelection()

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
