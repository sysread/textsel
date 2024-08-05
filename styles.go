package textsel

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type formatCode struct {
	fg    tcell.Color
	bg    tcell.Color
	style string
}

func newFormatCode() formatCode {
	return formatCode{
		fg:    tview.Styles.PrimaryTextColor,
		bg:    tview.Styles.PrimitiveBackgroundColor,
		style: "",
	}
}

func (fc formatCode) String() string {
	return fmt.Sprintf("[%s:%s:%s]", fc.fg.String(), fc.bg.String(), fc.style)
}

func (fc formatCode) update(code string) formatCode {
	// Strip the [ and ] characters from the code
	code = strings.TrimLeft(code, "[")
	code = strings.TrimRight(code, "]")

	// Split the code into its components
	parts := strings.Split(code, ":")

	if len(parts) > 0 {
		switch parts[0] {
		case "":
			// current style persists
		case "-":
			fc.fg = tview.Styles.PrimaryTextColor
		default:
			fc.fg = tcell.GetColor(parts[0])
		}

		fc.bg = tview.Styles.PrimitiveBackgroundColor
		fc.style = ""
	}

	if len(parts) > 1 {
		switch parts[1] {
		case "":
			// current style persists
		case "-":
			fc.bg = tview.Styles.PrimitiveBackgroundColor
		default:
			fc.bg = tcell.GetColor(parts[1])
		}

		fc.style = ""
	}

	if len(parts) > 2 {
		switch parts[2] {
		case "":
			// current style persists
		case "-":
			fc.style = "-"
		default:
			fc.style = parts[2]
		}
	}

	return fc
}
