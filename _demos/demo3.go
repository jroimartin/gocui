// Copyright 2014 The gocui Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/jroimartin/gocui"
)

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if _, err := g.SetView("side", -1, -1, int(0.2*float32(maxX)), maxY-5); err != nil &&
		err != gocui.ErrorUnkView {
		return err
	}
	if _, err := g.SetView("main", int(0.2*float32(maxX)), -1, maxX, maxY-5); err != nil &&
		err != gocui.ErrorUnkView {
		return err
	}
	if _, err := g.SetView("cmdline", -1, maxY-5, maxX, maxY); err != nil &&
		err != gocui.ErrorUnkView {
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

	g.SetLayout(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, 0, quit); err != nil {
		log.Panicln(err)
	}

	err = g.MainLoop()
	if err != nil && err != gocui.ErrorQuit {
		log.Panicln(err)
	}
}
