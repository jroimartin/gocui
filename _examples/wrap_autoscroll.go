// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/awesome-gocui/gocui"
)

var tabCount = 0

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
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, print); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("main", maxX/2-7, maxY/2-1, maxX/2+7, maxY/2+4, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		if _, err := g.SetCurrentView("main"); err != nil {
			return err
		}
		v.Clear()
		v.Autoscroll = true
		v.Wrap = true

		fmt.Fprintln(v, "Hello world!")
	}

	return nil
}

func print(g *gocui.Gui, v *gocui.View) error {
	tabCount++
	if tabCount%10 == 0 {
		fmt.Fprintln(v, tabCount, "Hello!")
	} else {
		fmt.Fprintln(v, tabCount, "Hello Cruel World !")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
