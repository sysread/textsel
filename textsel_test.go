package textsel

import (
	"testing"
)

func TestNewTextSel(t *testing.T) {
	textSel := NewTextSel()
	if textSel == nil {
		t.Fatal("Failed to create TextSel instance")
	}
}

func TestSetText(t *testing.T) {
	textSel := NewTextSel()

	testText := "Hello, World!"
	textSel.SetText(testText)

	got := textSel.GetText(false)
	if got != testText {
		t.Errorf("SetText() failed. Got %v, expected %v", got, testText)
	}
}
