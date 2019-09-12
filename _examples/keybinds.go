package main

import (
	"log"

	"github.com/awesome-gocui/gocui"
)

// layout generates the view
func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		v.Write([]byte("Hello"))

		if _, err := g.SetCurrentView("hello"); err != nil {
			return err
		}
	}

	return nil
}

// quit stops the gui
func quit(_ *gocui.Gui, _ *gocui.View) error {
	return gocui.ErrQuit
}

func main() {
	// Create a gui
	g, err := gocui.NewGui(gocui.OutputNormal, false)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	// Add a manager function
	g.SetManagerFunc(layout)

	// This will set up the recovery for MustParse
	defer func() {
		if r := recover(); r != nil {
			log.Panicln("Error caught: ", r)
		}
	}()

	// The MustParse can panic, but only returns 2 values instead of 3
	keyForced, modForced := gocui.MustParse("q")
	if err := g.SetKeybinding("", keyForced, modForced, quit); err != nil {
		log.Panicln(err)
	}

	// The normal parse returns an key, a modifier and an error
	keyNormal, modNormal, err := gocui.Parse("Ctrl+c")
	if err != nil {
		log.Panicln(err)
	}
	if err = g.SetKeybinding("", keyNormal, modNormal, quit); err != nil {
		log.Panicln(err)
	}

	// Now just start a mainloop for the demo
	if err = g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
