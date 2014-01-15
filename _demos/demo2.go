// Copyright 2014 The gocui Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jroimartin/gocui"
)

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
		if err := v.SetCursor(v.CX, v.CY+1); err != nil {
			v.SetOrigin(v.OX, v.OY+1)
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		if err := v.SetCursor(v.CX, v.CY-1); err != nil && v.OY > 0 {
			v.SetOrigin(v.OX, v.OY-1)
		}
	}
	return nil
}

func cursorLeft(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		if err := v.SetCursor(v.CX-1, v.CY); err != nil && v.OX > 0 {
			v.SetOrigin(v.OX-1, v.OY)
		}
	}
	return nil
}

func cursorRight(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		if err := v.SetCursor(v.CX+1, v.CY); err != nil {
			v.SetOrigin(v.OX+1, v.OY)
		}
	}
	return nil
}

func clear(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		v.Clear()
	}
	return nil
}

func writeTest(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		fmt.Fprintln(v, "This is a test")
	}
	return nil
}
func setLayout1(g *gocui.Gui, v *gocui.View) error {
	g.SetLayout(layout)
	return nil
}
func setLayout2(g *gocui.Gui, v *gocui.View) error {
	g.SetLayout(layout2)
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
	if err := g.SetKeybinding("", gocui.KeyCtrlC, 0, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", 'q', 0, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'c', 0, clear); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 't', 0, writeTest); err != nil {
		return err
	}
	if err := g.SetKeybinding("", '1', 0, setLayout1); err != nil {
		return err
	}
	if err := g.SetKeybinding("", '2', 0, setLayout2); err != nil {
		return err
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrorQuit
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("side", -1, -1, 30, maxY-5); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		v.Highlight = true
		fmt.Fprintln(v, "Item 1")
		fmt.Fprintln(v, "Item 2")
		fmt.Fprintln(v, "Item 3")
		fmt.Fprintln(v, "Item 4")
	}
	if v, err := g.SetView("main", 30, -1, maxX, maxY-5); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		b, err := ioutil.ReadFile("Mark.Twain-Tom.Sawyer.txt")
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(v, "%s", b)
		if err := g.SetCurrentView("main"); err != nil {
			return err
		}
	}
	if v, err := g.SetView("cmdline", -1, maxY-5, maxX, maxY); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		fmt.Fprintln(v, "Command line test")
	}
	return nil
}

func layout2(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("center", maxX/2-10, maxY/2-10, maxX/2+10, maxY/2+10); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		fmt.Fprintln(v, "Center view test")
		if err := g.SetCurrentView("center"); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	var err error

	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetLayout(layout)
	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}
	g.SelBgColor = gocui.ColorGreen
	g.SelFgColor = gocui.ColorBlack
	g.ShowCursor = true

	err = g.MainLoop()
	if err != nil && err != gocui.ErrorQuit {
		log.Panicln(err)
	}
}
