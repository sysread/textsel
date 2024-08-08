package textsel

import (
	"sync"
	"testing"
)

func TestTextSelection(t *testing.T) {
	var sem sync.WaitGroup

	textSel := NewTextSel()
	textSel.SetText("Hello, World!")

	selectedText := ""
	sem.Add(1)
	textSel.SetSelectFunc(func(text string) {
		selectedText = text
		sem.Done()
	})

	if textSel.cursorRow != 0 || textSel.cursorCol != 0 {
		t.Errorf("Initial cursor position failed. Expected cursorRow = 0, cursorCol = 0, got cursorRow = %d, cursorCol = %d", textSel.cursorRow, textSel.cursorCol)
	}

	textSel.startSelection() // H
	textSel.moveRight()   // He
	textSel.moveRight()   // Hel
	textSel.finishSelection()

	sem.Wait()

	expected := "Hel"
	if selectedText != expected {
		t.Errorf("Text selection failed. Expected %v, got %v", expected, selectedText)
	}

	if textSel.selectionStartRow != 0 {
		t.Errorf("selectionStartRow was not reset (actual: %v)", textSel.selectionStartRow)
	}

	if textSel.selectionStartCol != 0 {
		t.Errorf("selectionStartCol was not reset (actual: %v)", textSel.selectionStartCol)
	}

	if textSel.selectionEndRow != 0 {
		t.Errorf("selectionEndRow was not reset (actual: %v)", textSel.selectionEndRow)
	}

	if textSel.selectionEndCol != 0 {
		t.Errorf("selectionEndCol was not reset (actual: %v)", textSel.selectionEndCol)
	}
}

func TestSetSelectFunc(t *testing.T) {
	textSel := NewTextSel()
	textSel.SetText("Hello, World!")

	var selectedText string
	textSel.SetSelectFunc(func(text string) {
		selectedText = text
	})

	textSel.startSelection()
	textSel.finishSelection()

	expected := "H"
	if selectedText != expected {
		t.Errorf("SelectFunc failed. Expected %v, got %v", expected, selectedText)
	}
}
