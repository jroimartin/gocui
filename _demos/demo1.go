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
	if _, err := g.SetView("v1", -1, -1, 10, 10); err != nil &&
		err != gocui.ErrorUnkView {
		return err
	}
	if _, err := g.SetView("v2", maxX-10, -1, maxX, 10); err != nil &&
		err != gocui.ErrorUnkView {
		return err
	}
	if _, err := g.SetView("v3", maxX/2-5, -1, maxX/2+5, 10); err != nil &&
		err != gocui.ErrorUnkView {
		return err
	}
	if _, err := g.SetView("v4", -1, maxY/2-5, 10, maxY/2+5); err != nil &&
		err != gocui.ErrorUnkView {
		return err
	}
	if _, err := g.SetView("v5", maxX-10, maxY/2-5, maxX, maxY/2+5); err != nil &&
		err != gocui.ErrorUnkView {
		return err
	}
	if _, err := g.SetView("v6", -1, maxY-10, 10, maxY); err != nil &&
		err != gocui.ErrorUnkView {
		return err
	}
	if _, err := g.SetView("v7", maxX-10, maxY-10, maxX, maxY); err != nil &&
		err != gocui.ErrorUnkView {
		return err
	}
	if _, err := g.SetView("v8", maxX/2-5, maxY-10, maxX/2+5, maxY); err != nil &&
		err != gocui.ErrorUnkView {
		return err
	}
	if _, err := g.SetView("v9", maxX/2-5, maxY/2-5, maxX/2+5, maxY/2+5); err != nil &&
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
