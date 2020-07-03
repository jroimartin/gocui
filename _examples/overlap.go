// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/awesome-gocui/gocui"
)

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if _, err := g.SetView("v1", -1, -1, 10, 10, 0); err != nil &&
		!gocui.IsUnknownView(err) {
		return err
	}
	if _, err := g.SetView("v2", maxX-10, -1, maxX, 10, 0); err != nil &&
		!gocui.IsUnknownView(err) {
		return err
	}
	if _, err := g.SetView("v3", maxX/2-5, -1, maxX/2+5, 10, 0); err != nil &&
		!gocui.IsUnknownView(err) {
		return err
	}
	if _, err := g.SetView("v4", -1, maxY/2-5, 10, maxY/2+5, 0); err != nil &&
		!gocui.IsUnknownView(err) {
		return err
	}
	if _, err := g.SetView("v5", maxX-10, maxY/2-5, maxX, maxY/2+5, 0); err != nil &&
		!gocui.IsUnknownView(err) {
		return err
	}
	if _, err := g.SetView("v6", -1, maxY-10, 10, maxY, 0); err != nil &&
		!gocui.IsUnknownView(err) {
		return err
	}
	if _, err := g.SetView("v7", maxX-10, maxY-10, maxX, maxY, 0); err != nil &&
		!gocui.IsUnknownView(err) {
		return err
	}
	if _, err := g.SetView("v8", maxX/2-5, maxY-10, maxX/2+5, maxY, 0); err != nil &&
		!gocui.IsUnknownView(err) {
		return err
	}
	if _, err := g.SetView("v9", maxX/2-5, maxY/2-5, maxX/2+5, maxY/2+5, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		if _, err := g.SetCurrentView("v9"); err != nil {
			return err
		}
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && !gocui.IsQuit(err) {
		log.Panicln(err)
	}
}
