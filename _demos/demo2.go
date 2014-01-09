package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if _, err := g.SetView("side", -1, -1, 30, maxY-5); err != nil {
		return err
	}
	if _, err := g.SetView("main", 30, -1, maxX, maxY-5); err != nil {
		return err
	}
	if _, err := g.SetView("cmdline", -1, maxY-5, maxX, maxY); err != nil {
		return err
	}
	return nil
}

func start(g *gocui.Gui) error {
	if err := keybindings(g); err != nil {
		return err
	}
	if err := g.SetCurrentView("main"); err != nil {
		return err
	}
	if v := g.GetView("main"); v != nil {
		fmt.Fprintln(v, "This is a test")
	}
	if v := g.GetView("side"); v != nil {
		fmt.Fprintln(v, "Item 1")
		fmt.Fprintln(v, "Item 2")
		fmt.Fprintln(v, "Item 3")
		fmt.Fprintln(v, "Item 4")
	}
	if v := g.GetView("cmdline"); v != nil {
		fmt.Fprintln(v, "Buffer test")
	}
	g.ShowCursor = true
	return nil
}

func focusMain(g *gocui.Gui, v *gocui.View) error {
	return g.SetCurrentView("main")
}

func focusSide(g *gocui.Gui, v *gocui.View) error {
	return g.SetCurrentView("side")

}

func focusCmdLine(g *gocui.Gui, v *gocui.View) error {
	return g.SetCurrentView("cmdline")

}

func showHideCursor(g *gocui.Gui, v *gocui.View) error {
	g.ShowCursor = !g.ShowCursor
	return nil

}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		v.SetCursor(v.CX, v.CY+1)
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		v.SetCursor(v.CX, v.CY-1)
	}
	return nil
}

func cursorLeft(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		v.SetCursor(v.CX-1, v.CY)
	}
	return nil
}

func cursorRight(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		v.SetCursor(v.CX+1, v.CY)
	}
	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlM, 0, focusMain); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlS, 0, focusSide); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlL, 0, focusCmdLine); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'c', gocui.ModAlt, showHideCursor); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'j', 0, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'k', 0, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'h', 0, cursorLeft); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'l', 0, cursorRight); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyCtrlC, 0, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", 'q', 0, quit); err != nil {
		return err
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrorQuit
}

func main() {
	var err error

	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Layout = layout
	g.Start = start

	err = g.MainLoop()
	if err != nil && err != gocui.ErrorQuit {
		log.Panicln(err)
	}
}
