package textsel

import (
	"testing"
)

func TestCursorMovements(t *testing.T) {
	textSel := NewTextSel()
	textSel.SetText("Line 1\nLine 2\nLine 3")

	// Initial cursor position
	textSel.resetCursor()
	textSel.moveRight()
	if textSel.cursorCol != 1 {
		t.Errorf("Cursor movement failed. Expected cursorCol = 1, got = %d", textSel.cursorCol)
	}

	textSel.moveDown()
	if textSel.cursorRow != 1 {
		t.Errorf("Cursor movement failed. Expected cursorRow = 1, got = %d", textSel.cursorRow)
	}

	textSel.moveToEndOfLine()
	if textSel.cursorCol != len("Line 2\n")-1 {
		t.Errorf("Cursor move to end of line failed. Expected cursorCol = %d, got = %d", len("Line 2\n")-1, textSel.cursorCol)
	}
}
