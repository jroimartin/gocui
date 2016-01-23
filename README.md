# GOCUI - Go Console User Interface

Minimalist Go package aimed at creating Console User Interfaces.

## Features

* Minimalist API.
* Views (the "windows" in the GUI) implement the interface io.Writer.
* Support for overlapping views.
* The GUI can be modified at runtime.
* Global and view-level keybindings.
* Edit mode.

## Example

```go
func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("center", maxX/2-10, maxY/2, maxX/2+10, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "This is an example")
	}
	return nil
}
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
func main() {
	var err error
	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.SetLayout(layout)
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	err = g.MainLoop()
	if err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
```

## Screenshots

_examples/demo.go:

![_examples/demo.go](https://cloud.githubusercontent.com/assets/1223476/5992750/720b84f0-aa36-11e4-88ec-296fa3247b52.png)

_examples/delete.go:

![_examples/delete.go](https://cloud.githubusercontent.com/assets/1223476/5992751/76ad5cc2-aa36-11e4-8204-6a90269db827.png)

## Installation

`go get github.com/jroimartin/gocui`

## Documentation

`godoc github.com/jroimartin/gocui`
