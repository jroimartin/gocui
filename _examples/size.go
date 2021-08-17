// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)
	g.SetResizeFunc(onresize)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	_, err := g.SetView("size", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	return nil
}

func onresize(g *gocui.Gui, x, y int) error {
	v, err := g.View("size")
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		return nil
	}
	v.Clear()
	fmt.Fprintf(v, "%d, %d", x, y)
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
