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
	g, err := gocui.NewGui(gocui.Output256)

	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("colors", -1, -1, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		var count int
		for i := 1; i <= 255; i++ {
			count++

			block := fmt.Sprintf("\x1b[48;5;%dm%3s\x1b[0m", i, "  ")
			str := fmt.Sprintf("\x1b[38;5;%dm%3d %s\x1b[0m\t", i, i, block)

			if count == 10 {
				str += "\n\n"
				count = 0
			}

			fmt.Fprint(v, str)

		}
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
