package textsel

import (
	"testing"
)

func TestNewTextSel(t *testing.T) {
	ts := NewTextSel()

	if ts == nil {
		t.Fatal("Failed to create TextSel instance")
	}
}

func TestGetAndSetText(t *testing.T) {
	ts := NewTextSel().SetText("Hello, World!")

	got := ts.GetText(false)
	if got != "Hello, World!" {
		t.Errorf("GetText() returned the wrong text:\nExpected: '%v'\nGot: '%v'", "Hello, World!", got)
	}
}
