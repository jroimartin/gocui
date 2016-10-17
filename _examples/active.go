// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

const LEN = 4

var (
	viewArr = [LEN]string{"v1", "v2", "v3", "v4"}
	active  = 0
)

func nextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (active + 1) % LEN
	name := viewArr[nextIndex]
	out, _ := g.View("v2")
	fmt.Fprintln(out, "Going from view "+v.Name()+" to "+name)

	if err := g.SetCurrentViewOnTop(name); err != nil {
		fmt.Fprintln(out, err)
		return err
	}

	if nextIndex == 0 || nextIndex == 3 {
		g.Cursor = true
	} else {
		g.Cursor = false
	}

	active = nextIndex
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("v1", 0, 0, maxX/2-1, maxY/2-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "v1 (editable)"
		v.Editable = true
		v.Wrap = true
	}

	if v, err := g.SetView("v2", maxX/2-1, 0, maxX-1, maxY/2-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "v2"
		v.Editable = false
		v.Wrap = true
		v.Autoscroll = true
		v.ActiveColor = gocui.ColorRed
	}
	if v, err := g.SetView("v3", 0, maxY/2-1, maxX/2-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "v3"
		v.Write([]byte("Press TAB to change current view"))
		v.Wrap = true
		v.Editable = false
		v.Autoscroll = true
		v.ActiveColor = gocui.ColorRed
	}
	if v, err := g.SetView("v4", maxX/2, maxY/2, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "v4 (editable)"
		v.Editable = true

		if err = g.SetCurrentViewOnTop("v1"); err != nil {
			return err
		}
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func main() {
	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.ActiveColor = gocui.ColorGreen
	g.SetLayout(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
