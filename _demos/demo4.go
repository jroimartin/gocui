// Copyright 2014 The gocui Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("main", 1, 1, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		v.Wrap = true
		v.WrapPrefix = "> "
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

	go func() {
		var line bytes.Buffer
		for i := 0; i < 10; i++ {
			line.WriteString("This is a long line -- ")
		}
		fmt.Fprint(g.View("main"), line.String())

		fmt.Fprintln(g.View("main"), "")

		fmt.Fprint(g.View("main"), "Short")
	}()

	err = g.MainLoop()
	if err != nil && err != gocui.ErrorQuit {
		log.Panicln(err)
	}

}
