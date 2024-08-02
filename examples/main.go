package main

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/sysread/textsel"
)

const text = `
This is an example of the tview-textsel package.

Use arrow keys to move the cursor.

Press space to start/stop selecting text.
`

func main() {
	app := tview.NewApplication()

	textSel := textsel.NewTextSel().SetText(text)
	textSel.SetSelectFunc(func(selectedText string) {
		app.Stop()
		fmt.Println("Selected text:\n\n", selectedText)
	})

	if err := app.SetRoot(textSel, true).Run(); err != nil {
		panic(err)
	}
}
