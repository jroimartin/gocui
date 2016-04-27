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
	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetLayout(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("hello", maxX/2-6, maxY/2-5, maxX/2+6, maxY/2+5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "\x1b[0;30;47mHello world")
		fmt.Fprintln(v, "\x1b[0;31mHello world")
		fmt.Fprintln(v, "\x1b[0;32mHello world")
		fmt.Fprintln(v, "\x1b[0;33mHello world")
		fmt.Fprintln(v, "\x1b[0;34mHello world")
		fmt.Fprintln(v, "\x1b[0;35mHello world")
		fmt.Fprintln(v, "\x1b[0;36mHello world")
		fmt.Fprintln(v, "\x1b[0;37mHello world")
		fmt.Fprintln(v, "\x1b[0;30;41mHello world")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
