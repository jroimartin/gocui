// Copyright 2015 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

func main() {
	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Fatalln(err)
	}
	defer g.Close()

	g.SetLayout(layout)
	if err := initKeybindings(g); err != nil {
		log.Fatalln(err)
	}
	g.Cursor = true

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("legend", maxX-23, 0, maxX-1, 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Keybindings"
		fmt.Fprintln(v, "^a: Set mask")
		fmt.Fprintln(v, "^c: Exit")
	}

	if v, err := g.SetView("input", 0, 0, maxX-24, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if err := g.SetCurrentView("input"); err != nil {
			return err
		}
		v.Editable = true
		v.Wrap = true
	}

	return nil
}

func initKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("input", gocui.KeyCtrlA, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			v.Mask ^= '*'
			return nil
		}); err != nil {
		return err
	}
	return nil
}
