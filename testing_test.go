// Copyright 2020 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"errors"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestTestingScreenReturnsCorrectContent(t *testing.T) {
	didCallCTRLC := false
	viewContent := "Hello world!"
	viewName := "testView1"

	g, err := NewGui(OutputSimulator, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(func (g *Gui) error {
		maxX, maxY := g.Size()
		if v, err := g.SetView(viewName, maxX/2-7, maxY/2, maxX/2+7, maxY/2+2, 0); err != nil {
			if !errors.Is(err, ErrUnknownView) {
				return err
			}
	
			if _, err := g.SetCurrentView(viewName); err != nil {
				return err
			}
	
			fmt.Fprintln(v, viewContent)
		}

		return nil
	})

	if err := g.SetKeybinding("", KeyCtrlC, ModNone, func(g *Gui, v *View) error { didCallCTRLC = true; return nil }); err != nil {
		log.Panicln(err)
	}

	testCompleted := false
	go func() {
		if err := g.MainLoop(); (err != nil && !errors.Is(err, ErrQuit)) || testCompleted {
			log.Panicln(err)
		}
	}()

	testingScreen := g.GetTestingScreen()
	testingScreen.SendKey(KeyCtrlC)

	<-time.After(time.Second * 5)

	if !didCallCTRLC {
		t.Error("Simulator didn't send CTRLC command correctly")
	}

	// Get the content of the "hello" view
	content, err := testingScreen.GetViewContent(viewName)
	if err != nil {
		t.Error(err)
	}

	if content != viewContent {
		t.Error("View content doesn't match expected")
	}
}
