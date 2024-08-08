package textsel

import (
	"fmt"
	"os"
)

// Debug function to write to `debug.log`.
func (ts *TextSel) debug(format string, args ...interface{}) {
	file, _ := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer file.Close()
	line := fmt.Sprintf(format, args...)
	file.WriteString(line + "\n")
}

// Debug function to log the color codes being used for the cursor and
// selection highlighting.
func (ts *TextSel) debugColors() *TextSel {
	ts.debug("          defaultColor: %s", ts.defaultColor)
	ts.debug("           cursorColor: %s", ts.cursorColor)
	ts.debug("        selectionColor: %s", ts.selectionColor)
	ts.debug("cursorInSelectionColor: %s", ts.cursorInSelectionColor)
	return ts
}

// Debug function to log the cursor position.
func (ts *TextSel) debugCursor() *TextSel {
	ts.debug("Cursor: (%d, %d)", ts.cursorRow, ts.cursorCol)
	return ts
}

// Debug function to log the selection range.
func (ts *TextSel) debugSelection() *TextSel {
	if ts.isSelecting {
		startRow, startCol, endRow, endCol := ts.getSelectionRange()
		ts.debug("Selection range: (%d, %d) - (%d, %d)", startRow, startCol, endRow, endCol)
	} else {
		ts.debug("No selection")
	}
	return ts
}
