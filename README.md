# TextSel

[![Go Reference](https://pkg.go.dev/badge/github.com/sysread/textsel.svg)](https://pkg.go.dev/github.com/sysread/textsel#section-readme)

`textsel` is a Go package that extends the `tview.TextView` widget to include
cursor movement and text selection capabilities. It allows users to navigate
and select text within the `TextView`, making it easier to copy and highlight
text.


## Features

- Cursor movement (up, down, left, right).
- Text selection with visual highlighting.
- Customizable colors for cursor and selection.

## Installation

To install `textsel`, you need to have Go installed and set up on your machine.
Use the following command to install the package:

```bash
go get github.com/sysread/textsel
```

## Usage

Here's an example of how to use `textsel` in your project:

```go
package main

import (
    "github.com/sysread/textsel"
    "github.com/rivo/tview"
)

func main() {
    app := tview.NewApplication()

    textSel := textsel.NewTextSel().SetText("Hello, World!")

    if err := app.SetRoot(textSel, true).Run(); err != nil {
        panic(err)
    }
}
```

## Customization

You can customize the colors used for cursor and selection highlighting by
modifying the `defaultColor`, `cursorColor`, `selectionColor`, and
`cursorInSelectionColor` fields in the `TextSel` struct.

Example:

```go
textSel := textsel.NewTextSel()
textSel.defaultColor = "[#FFFFFF:#000000:-]"
textSel.cursorColor = "[#000000:#FFFFFF:-]"
textSel.selectionColor = "[#000000:#FF0000:-]"
textSel.cursorInSelectionColor = "[#000000:#FF0000:bu]"
```

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request with your improvements.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
