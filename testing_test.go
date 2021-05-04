// Copyright 2021 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"
)

func TestTestingScreenReturnsCorrectContent(t *testing.T) {
	// Track what happened in the view, we'll assert on these
	didCallCTRLC := false
	expectedViewContent := "Hello world!"
	viewName := "testView1"

	// Create a view specifying the "OutputSimulator" mode
	g, err := NewGui(OutputSimulator, true)
	if err != nil {
		log.Panicln(err)
	}
	g.SetManagerFunc(func(g *Gui) error {
		maxX, maxY := g.Size()
		if v, err := g.SetView(viewName, maxX/2-7, maxY/2, maxX/2+7, maxY/2+2, 0); err != nil {
			if !errors.Is(err, ErrUnknownView) {
				return err
			}

			if _, err := g.SetCurrentView(viewName); err != nil {
				return err
			}

			// Have the view draw "Hello world!"
			fmt.Fprintln(v, expectedViewContent)
		}

		return nil
	})

	// Create a key binding which sets "didCallCTRLC" when triggered
	exampleBindingToTest := func(g *Gui, v *View) error {
		didCallCTRLC = true
		return nil
	}
	if err := g.SetKeybinding("", KeyCtrlC, ModNone, exampleBindingToTest); err != nil {
		log.Panicln(err)
	}

	// Create a test screen and start gocui
	testingScreen := g.GetTestingScreen()
	cleanup := testingScreen.StartGui()
	defer cleanup()

	// Send a key to gocui
	testingScreen.SendKey(KeyCtrlC)

	// Wait for key to be processed
	<-time.After(time.Millisecond * 50)

	// Test that the keybinding fired and set "didCallCTRLC" to true
	if !didCallCTRLC {
		t.Error("Expect the simulator to invoke the key handler for CTRLC")
	}

	// Get the content from the testing screen
	actualContent, err := testingScreen.GetViewContent(viewName)
	if err != nil {
		t.Error(err)
	}

	// Test that it contains the "Hello World!" we thought the view should draw
	if strings.TrimSpace(actualContent) != expectedViewContent {
		t.Error(fmt.Printf("Expected view content to be: %q got: %q", expectedViewContent, actualContent))
	}
}
