package textsel

import (
	"testing"
)

func TestGetCurrentLine(t *testing.T) {
	ts := NewTextSel()
	ts.SetText("Hello, World!\nThis is a test.\nThis is only a test.\n")

	if ts.getCurrentLine() != "Hello, World!\n" {
		t.Error("getCurrentLine() failed on first line")
	}

	ts.moveDown()

	if ts.getCurrentLine() != "This is a test.\n" {
		t.Error("getCurrentLine() failed on second line after moving cursor down")
	}

	ts.moveRight().moveRight().moveRight().moveRight()

	if ts.getCurrentLine() != "This is a test.\n" {
		t.Error("getCurrentLine() failed on second line after moving cursor right")
	}
}

func TestLastRow(t *testing.T) {
	ts := NewTextSel()
	ts.SetText("Hello, World!\nThis is a test.\nThis is only a test.\n")

	if ts.lastRow() != 2 {
		t.Error("lastRow() failed")
	}
}
