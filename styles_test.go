package textsel

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func TestNewFormatCode(t *testing.T) {
	fc := newFormatCode()

	if fc.fg != tview.Styles.PrimaryTextColor {
		t.Errorf("newFormatCode() failed. Expected fg = '%v', got = '%v'", tview.Styles.PrimaryTextColor, fc.fg)
	}

	if fc.bg != tview.Styles.PrimitiveBackgroundColor {
		t.Errorf("newFormatCode() failed. Expected bg = '%v', got = '%v'", tview.Styles.PrimitiveBackgroundColor, fc.bg)
	}

	if fc.style != "" {
		t.Errorf("newFormatCode() failed. Expected style = '', got = '%v'", fc.style)
	}
}

func TestFormatCodeString(t *testing.T) {
	fc := formatCode{
		fg:    tcell.ColorRed,
		bg:    tcell.ColorBlue,
		style: "b",
	}

	expected := "[red:blue:b]"
	got := fc.String()

	if got != expected {
		t.Errorf("formatCode.String() failed. Expected '%v', got '%v'", expected, got)
	}
}

func TestUpdateFormatCode1(t *testing.T) {
	original := newFormatCode()
	updated := original.update("[green::]")

	expectedFg := tcell.GetColor("green")
	if updated.fg != expectedFg {
		t.Errorf("update() failed. Expected fg = '%v', got = '%v'", expectedFg, updated.fg)
	}

	expectedBg := tview.Styles.PrimitiveBackgroundColor
	if updated.bg != expectedBg {
		t.Errorf("update() failed. Expected bg = '%v', got = '%v'", expectedBg, updated.bg)
	}

	expectedStyle := ""
	if updated.style != expectedStyle {
		t.Errorf("update() failed. Expected style = '%v', got = '%v'", expectedStyle, updated.style)
	}
}

func TestUpdateFormatCode2(t *testing.T) {
	original := newFormatCode()
	updated := original.update("[red:yellow:-]")

	expectedFg := tcell.GetColor("red")
	if updated.fg != expectedFg {
		t.Errorf("update() failed. Expected fg = '%v', got = '%v'", expectedFg, updated.fg)
	}

	expectedBg := tcell.GetColor("yellow")
	if updated.bg != expectedBg {
		t.Errorf("update() failed. Expected bg = '%v', got = '%v'", expectedBg, updated.bg)
	}

	expectedStyle := "-"
	if updated.style != expectedStyle {
		t.Errorf("update() failed. Expected style = '%v', got = '%v'", expectedStyle, updated.style)
	}
}

func TestUpdateFormatCode3(t *testing.T) {
	original := newFormatCode()
	updated := original.update("[red:yellow:b]")

	expectedFg := tcell.GetColor("red")
	if updated.fg != expectedFg {
		t.Errorf("update() failed. Expected fg = '%v', got = '%v'", expectedFg, updated.fg)
	}

	expectedBg := tcell.GetColor("yellow")
	if updated.bg != expectedBg {
		t.Errorf("update() failed. Expected bg = '%v', got = '%v'", expectedBg, updated.bg)
	}

	expectedStyle := "b"
	if updated.style != expectedStyle {
		t.Errorf("update() failed. Expected style = '%v', got = '%v'", expectedStyle, updated.style)
	}
}

func TestUpdateFormatCode4(t *testing.T) {
	original := newFormatCode()
	updated := original.update("[:-:i]")

	expectedFg := tview.Styles.PrimaryTextColor
	if updated.fg != expectedFg {
		t.Errorf("update() failed. Expected fg = '%v', got = '%v'", expectedFg, updated.fg)
	}

	expectedBg := tview.Styles.PrimitiveBackgroundColor
	if updated.bg != expectedBg {
		t.Errorf("update() failed. Expected bg = '%v', got = '%v'", expectedBg, updated.bg)
	}

	expectedStyle := "i"
	if updated.style != expectedStyle {
		t.Errorf("update() failed. Expected style = '%v', got = '%v'", expectedStyle, updated.style)
	}
}
