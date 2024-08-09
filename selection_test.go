package textsel

import (
	"sync"
	"testing"
)

func TestResetSelection(t *testing.T) {
	ts := NewTextSel().SetText("Hello, World!")

	ts.StartSelection()
	ts.MoveRight()
	ts.MoveRight()

	got := ts.GetSelectedText()
	if got != "Hel" {
		t.Errorf("StartSelection failed. Expected 'Hel', got: '%s'", got)
	}

	ts.ResetSelection()
	got = ts.GetSelectedText()
	if got != "" {
		t.Errorf("ResetSelection failed. Expected '', got: '%s'", got)
	}
}

func TestGetSelectionRange(t *testing.T) {
	ts := NewTextSel().
		SetText("Hello\nWorld")

	ts.
		ResetCursor().
		MoveRight().
		MoveRight().
		StartSelection().
		MoveDown().
		MoveRight()

	startRow, startCol, endRow, endCol := ts.GetSelectionRange()
	if startRow != 0 || startCol != 2 || endRow != 1 || endCol != 3 {
		t.Errorf("GetSelectionRange failed. Expected (0, 2, 1, 3), got (%d, %d, %d, %d)", startRow, startCol, endRow, endCol)
	}

	// Backwards selection
	ts.
		ResetCursor().
		MoveDown().
		MoveToEndOfLine().
		StartSelection().
		MoveUp().
		MoveToStartOfLine()

	startRow, startCol, endRow, endCol = ts.GetSelectionRange()

	if startRow != 0 || startCol != 0 || endRow != 1 || endCol != 4 {
		t.Errorf("GetSelectionRange failed. Expected (0, 0, 1, 4), got (%d, %d, %d, %d)", startRow, startCol, endRow, endCol)
	}
}

func TestTextSelection(t *testing.T) {
	var sem sync.WaitGroup

	ts := NewTextSel()
	ts.SetText("Hello, World!")

	selectedText := ""
	sem.Add(1)
	ts.SetSelectFunc(func(text string) {
		selectedText = text
		sem.Done()
	})

	if ts.cursorRow != 0 || ts.cursorCol != 0 {
		t.Errorf("Initial cursor position failed. Expected cursorRow = 0, cursorCol = 0, got cursorRow = %d, cursorCol = %d", ts.cursorRow, ts.cursorCol)
	}

	ts.StartSelection() // H
	ts.MoveRight()      // He
	ts.MoveRight()      // Hel
	ts.FinishSelection()

	sem.Wait()

	expected := "Hel"
	if selectedText != expected {
		t.Errorf("Text selection failed. Expected %v, got %v", expected, selectedText)
	}

	if ts.selectionStartRow != 0 {
		t.Errorf("selectionStartRow was not reset (actual: %v)", ts.selectionStartRow)
	}

	if ts.selectionStartCol != 0 {
		t.Errorf("selectionStartCol was not reset (actual: %v)", ts.selectionStartCol)
	}

	if ts.selectionEndRow != 0 {
		t.Errorf("selectionEndRow was not reset (actual: %v)", ts.selectionEndRow)
	}

	if ts.selectionEndCol != 0 {
		t.Errorf("selectionEndCol was not reset (actual: %v)", ts.selectionEndCol)
	}
}

func TestSetSelectFunc(t *testing.T) {
	ts := NewTextSel()
	ts.SetText("Hello, World!")

	var selectedText string
	ts.SetSelectFunc(func(text string) {
		selectedText = text
	})

	ts.StartSelection()
	ts.FinishSelection()

	expected := "H"
	if selectedText != expected {
		t.Errorf("SelectFunc failed. Expected %v, got %v", expected, selectedText)
	}
}
