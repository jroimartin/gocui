# GOCUI - Go Console User Interface

Minimalist Go package aimed at creating Console User Interfaces.

## Installation

`go get github.com/jroimartin/gocui`

## Documentation

`godoc github.com/jroimartin/gocui`

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
		if err != gocui.ErrorUnkView {
			return err
		}
		fmt.Fprintln(v, "This is an example")
	}
	return nil
}
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.Quit
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
	if err != nil && err != gocui.Quit {
		log.Panicln(err)
	}
}
```
