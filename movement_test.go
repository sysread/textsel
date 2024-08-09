package textsel

import (
	"testing"
)

func TestGetAndSetCursorPosition(t *testing.T) {
	ts := NewTextSel().SetText("Hello\nWorld")

	ts.SetCursorPosition(1, 2)

	row, col := ts.GetCursorPosition()
	if row != 1 || col != 2 {
		t.Errorf("GetCursorPosition failed. Expected cursorRow = 1, cursorCol = 2, got = %d, %d", row, col)
	}

	ts.ResetCursor()
	ts.StartSelection()
	ts.SetCursorPosition(1, 2)

	selectedText := ts.GetSelectedText()
	if selectedText != "Hello\nWor" {
		t.Errorf("SetCursorPosition failed to correctly update selection. Expected 'Hello\nWo', got = %s", selectedText)
	}
}

func TestResetCursor(t *testing.T) {
	ts := NewTextSel().
		SetText("Hello\nWorld").
		SetCursorPosition(1, 2).
		ResetCursor()

	row, col := ts.GetCursorPosition()
	if row != 0 || col != 0 {
		t.Errorf("ResetCursor failed. Expected cursorRow = 0, cursorCol = 0, got = %d, %d", row, col)
	}
}

func TestMoveRight(t *testing.T) {
	ts := NewTextSel().SetText("a\nb\n")

	ts.MoveRight()
	row, col := ts.GetCursorPosition()
	if row != 0 || col != 1 {
		t.Errorf("MoveRight failed. Expected cursorRow = 0, cursorCol = 1, got = %d, %d", row, col)
	}

	ts.MoveRight()
	row, col = ts.GetCursorPosition()
	if row != 1 || col != 0 {
		t.Errorf("MoveRight failed to wrap to the next row. Expected cursorRow = 1, cursorCol = 0, got = %d, %d", row, col)
	}

	ts.MoveRight()
	row, col = ts.GetCursorPosition()
	if row != 1 || col != 1 {
		t.Errorf("MoveRight failed. Expected cursorRow = 1, cursorCol = 1, got = %d, %d", row, col)
	}

	ts.MoveRight()
	row, col = ts.GetCursorPosition()
	if row != 1 || col != 1 {
		t.Errorf("MoveRight failed to stop at EOF. Expected cursorRow = 1, cursorCol = 1, got = %d, %d", row, col)
	}
}

func TestMoveLeft(t *testing.T) {
	ts := NewTextSel().SetText("a\nb\n")

	ts.cursorRow = 1
	ts.cursorCol = 1

	ts.MoveLeft()
	row, col := ts.GetCursorPosition()
	if row != 1 || col != 0 {
		t.Errorf("MoveLeft failed. Expected cursorRow = 1, cursorCol = 0, got = %d, %d", row, col)
	}

	ts.MoveLeft()
	row, col = ts.GetCursorPosition()
	if row != 0 || col != 1 {
		t.Errorf("MoveLeft failed to wrap to the previous row. Expected cursorRow = 0, cursorCol = 1, got = %d, %d", row, col)
	}

	ts.MoveLeft()
	row, col = ts.GetCursorPosition()
	if row != 0 || col != 0 {
		t.Errorf("MoveLeft failed. Expected cursorRow = 0, cursorCol = 0, got = %d, %d", row, col)
	}

	ts.MoveLeft()
	row, col = ts.GetCursorPosition()
	if row != 0 || col != 0 {
		t.Errorf("MoveLeft failed to stop at BOF. Expected cursorRow = 0, cursorCol = 0, got = %d, %d", row, col)
	}
}

func TestMoveDown(t *testing.T) {
	ts := NewTextSel().SetText("a\nb\n")

	ts.MoveDown()
	row, col := ts.GetCursorPosition()
	if row != 1 || col != 0 {
		t.Errorf("MoveDown failed. Expected cursorRow = 1, cursorCol = 0, got = %d, %d", row, col)
	}

	ts.MoveDown()
	row, col = ts.GetCursorPosition()
	if row != 1 || col != 0 {
		t.Errorf("MoveDown failed to stop at EOF. Expected cursorRow = 1, cursorCol = 0, got = %d, %d", row, col)
	}
}

func TestMoveUp(t *testing.T) {
	ts := NewTextSel().SetText("a\nb\n")

	ts.cursorRow = 1
	ts.cursorCol = 1

	ts.MoveUp()
	row, col := ts.GetCursorPosition()
	if row != 0 || col != 1 {
		t.Errorf("MoveUp failed. Expected cursorRow = 0, cursorCol = 1, got = %d, %d", row, col)
	}

	ts.MoveUp()
	row, col = ts.GetCursorPosition()
	if row != 0 || col != 1 {
		t.Errorf("MoveUp failed to stop at BOF. Expected cursorRow = 0, cursorCol = 1, got = %d, %d", row, col)
	}
}

func TestMoveToEndOfLine(t *testing.T) {
	ts := NewTextSel().SetText("Hello\nWorld\n")

	ts.MoveToEndOfLine()
	row, col := ts.GetCursorPosition()
	if row != 0 || col != 5 {
		t.Errorf("MoveToEndOfLine failed. Expected cursorRow = 0, cursorCol = 5, got = %d, %d", row, col)
	}

	ts.MoveToEndOfLine()
	row, col = ts.GetCursorPosition()
	if row != 0 || col != 5 {
		t.Errorf("MoveToEndOfLine failed to stop at EOL. Expected cursorRow = 1, cursorCol = 5, got = %d, %d", row, col)
	}

	ts.MoveDown()
	ts.MoveToEndOfLine()
	row, col = ts.GetCursorPosition()
	if row != 1 || col != 5 {
		t.Errorf("MoveToEndOfLine failed. Expected cursorRow = 1, cursorCol = 5, got = %d, %d", row, col)
	}
}

func TestMoveToStartOfLine(t *testing.T) {
	ts := NewTextSel().SetText("Hello\nWorld\n")

	ts.MoveToEndOfLine().MoveToStartOfLine()
	row, col := ts.GetCursorPosition()
	if row != 0 || col != 0 {
		t.Errorf("MoveToStartOfLine failed. Expected cursorRow = 0, cursorCol = 0, got = %d, %d", row, col)
	}

	ts.MoveDown().MoveToEndOfLine().MoveToStartOfLine()
	row, col = ts.GetCursorPosition()
	if row != 1 || col != 0 {
		t.Errorf("MoveToStartOfLine failed. Expected cursorRow = 1, cursorCol = 0, got = %d, %d", row, col)
	}

	ts.MoveToStartOfLine()
	row, col = ts.GetCursorPosition()
	if row != 1 || col != 0 {
		t.Errorf("MoveToStartOfLine failed to stop at BOL. Expected cursorRow = 1, cursorCol = 0, got = %d, %d", row, col)
	}
}
