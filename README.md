# GOCUI - Go Console User Interface

[![GoDoc](https://godoc.org/github.com/jroimartin/gocui?status.svg)](https://godoc.org/github.com/jroimartin/gocui)

Minimalist Go package aimed at creating Console User Interfaces.

## Features

* Minimalist API.
* Views (the "windows" in the GUI) implement the interface io.ReadWriter.
* Support for overlapping views.
* The GUI can be modified at runtime (concurrent-safe).
* Global and view-level keybindings.
* Mouse support.
* Customizable edition mode.

## Installation

Execute:

```
$ go get github.com/jroimartin/gocui
```

## Documentation

Execute:

```
$ godoc github.com/jroimartin/gocui
```

Or visit [godoc.org](https://godoc.org/github.com/jroimartin/gocui) to read it
online.

## Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

func main() {
	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetLayout(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Hello world!")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
```

## Screenshots

_examples/demo.go:

![_examples/demo.go](https://cloud.githubusercontent.com/assets/1223476/5992750/720b84f0-aa36-11e4-88ec-296fa3247b52.png)

_examples/dynamic.go:

![_examples/delete.go](https://cloud.githubusercontent.com/assets/1223476/5992751/76ad5cc2-aa36-11e4-8204-6a90269db827.png)
