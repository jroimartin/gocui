// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package gocui allows to create console user interfaces.

Example:

	func layout(g *gocui.Gui) error {
		maxX, maxY := g.Size()
		if v, err := g.SetView("center", maxX/2-10, maxY/2, maxX/2+10, maxY/2+2); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			fmt.Fprintln(v, "This is an example")
		}
		return nil
	}
	func quit(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}
	func main() {
		var err error
		g := gocui.NewGui()
		if err := g.Init(); err != nil {
			log.Panicln(err)
		}
		defer g.Close()
		g.SetLayout(layout)
		if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
			log.Panicln(err)
		}
		err = g.MainLoop()
		if err != nil && err != gocui.ErrQuit {
			log.Panicln(err)
		}
	}
*/
package gocui
