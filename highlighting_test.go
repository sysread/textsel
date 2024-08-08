package textsel

import (
	"fmt"
	"testing"

	"github.com/rivo/tview"
)

func visualizeString(s string) string {
	var result string
	for _, r := range s {
		if r == '\n' {
			result += "\\n\n"
		} else if r == '\t' {
			result += "\\t"
		} else if r == ' ' {
			result += "·"
		} else {
			result += fmt.Sprintf("%c", r)
		}
	}
	return result
}

func TestHighlightCursor(t *testing.T) {
	ts := NewTextSel()

	app := tview.NewApplication()
	app.SetRoot(ts, true)
	app.SetFocus(ts)

	if !ts.HasFocus() {
		t.Errorf("TextSel does not have focus.")
	}

	// Test Case 1: Cursor at the beginning
	ts.SetText("Hello World").resetCursor().highlightCursor()
	expectedOutput1 := "[black:white:-]H[white:black:]ello World"
	actualOutput1 := ts.TextView.GetText(false)

	if actualOutput1 != expectedOutput1 {
		t.Errorf("Cursor highlight at start of string failed.\n\nExpected: '%s'\n\n  Actual: '%s'\n\n", visualizeString(expectedOutput1), visualizeString(actualOutput1))
	}

	// Test Case 2: Cursor at the end
	ts.SetText("Hello\nWorld").
		resetCursor().
		moveDown().
		moveRight().
		moveRight().
		moveRight().
		moveRight().
		highlightCursor()

	// note extra space at end of line; it is used to make the cursor visible
	expectedOutput2 := "Hello \nWorl[black:white:-]d[white:black:]"
	actualOutput2 := ts.TextView.GetText(false)

	if actualOutput2 != expectedOutput2 {
		t.Errorf("Cursor highlight at end of string failed.\n\nExpected: '%s'\n\n  Actual: '%s'\n\n", visualizeString(expectedOutput2), visualizeString(actualOutput2))
	}

	// Test Case 3: Selection highlighting
	ts.SetText("Hello\nWorld").
		resetCursor().
		startSelection().
		moveRight().
		moveRight().
		moveRight().
		moveRight().
		moveRight().
		highlightCursor()

	// note extra space at end of line; it is used to make the cursor visible
	expectedOutput3 := "[black:yellow:-]Hello[black:yellow:bu] \n[white:black:-][white:black:]World"
	actualOutput3 := ts.TextView.GetText(false)

	if actualOutput3 != expectedOutput3 {
		t.Errorf("Selection highlight failed.\n\nExpected: '%s'\n\n  Actual: '%s'\n\n", visualizeString(expectedOutput3), visualizeString(actualOutput3))
	}
}
