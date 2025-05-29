// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jroimartin/gocui"
)

var (
	viewArr = []string{"v1", "v2", "v3", "v4"}
	active  = 0
)

func nextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (active + 1) % len(viewArr)
	active = nextIndex
	if _, err := g.SetCurrentView(viewArr[active]); err != nil {
		return err
	}
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("v1", 0, 0, maxX/2-1, maxY/2-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "v1"

		v.FrameColor = gocui.ColorCyan
		fmt.Fprintf(v, "I am \033[36;1mcyan\033[0m, and \033[32;1mgreen\033[0m when focused\n")
		fmt.Fprintln(v, "Press TAB to change current view")
		if _, err := g.SetCurrentView("v1"); err != nil {
			return err
		}
	}

	if v, err := g.SetView("v2", maxX/2, 0, maxX-1, maxY/2-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "v2"

		v.FrameColor = gocui.ColorRed
		fmt.Fprint(v, "I am \033[31;1mred\033[0m")
	}
	if v, err := g.SetView("v3", 0, maxY/2, maxX/2-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "v3"

		v.FrameBgColor = gocui.ColorBlue

		fmt.Fprint(v, "I am default but my border is \033[34;1mblue\033[0m")
	}
	if v, err := g.SetView("v4", maxX/2, maxY/2, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "v4"
		rainbow(g, "v4")
		fmt.Fprint(v, "Rainbow")
	}
	return nil
}

func rainbow(g *gocui.Gui, viewName string) {
	colors := []gocui.Attribute{
		gocui.ColorRed,
		gocui.ColorGreen,
		gocui.ColorYellow,
		gocui.ColorBlue,
		gocui.ColorMagenta,
		gocui.ColorCyan,
		gocui.ColorWhite,
	}

	go func() {
		i := 0
		for {
			time.Sleep(time.Second)
			color := colors[i%len(colors)]
			i++
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View(viewName)
				if err != nil {
					return err
				}
				v.FrameColor = color
				return nil
			})
		}
	}()
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(layout)

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
