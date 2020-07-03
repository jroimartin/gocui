// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/awesome-gocui/gocui"
)

const NumGoroutines = 20

var (
	done = make(chan struct{})
	wg   sync.WaitGroup

	mu  sync.Mutex // protects ctr
	ctr = 0
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	for i := 0; i < NumGoroutines; i++ {
		wg.Add(1)
		go counter(g)
	}

	if err := g.MainLoop(); err != nil && !gocui.IsQuit(err) {
		log.Panicln(err)
	}

	wg.Wait()
}

func layout(g *gocui.Gui) error {
	if v, err := g.SetView("ctr", 2, 2, 22, 2+NumGoroutines+1, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		v.Clear()
		if _, err := g.SetCurrentView("ctr"); err != nil {
			return err
		}
	}
	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	close(done)
	return gocui.ErrQuit
}

func counter(g *gocui.Gui) {
	defer wg.Done()

	for {
		select {
		case <-done:
			return
		case <-time.After(500 * time.Millisecond):
			mu.Lock()
			n := ctr
			ctr++
			mu.Unlock()

			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("ctr")
				if err != nil {
					return err
				}
				// use ctr to make it more chaotic
				// "pseudo-randomly" print in one of two columns (x = 0, and x = 10)
				x := (ctr / NumGoroutines) & 1
				if x != 0 {
					x = 10
				}
				y := ctr % NumGoroutines
				v.SetWritePos(x, y)
				fmt.Fprintln(v, n)
				return nil
			})
		}
	}
}
