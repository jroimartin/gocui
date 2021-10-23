// Copyright 2021 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
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
	testingScreen.WaitSync()

	// Test that the keybinding fired and set "didCallCTRLC" to true
	if !didCallCTRLC {
		t.Error("Expect the simulator to invoke the key handler for CTRLC")
	}

	// check view content
	assertView(t, testingScreen, viewName, expectedViewContent)
}

func TestTestingScreenMultipleKeys(t *testing.T) {
	// Track what happened in the view, we'll assert on these
	didCallCTRLC := false
	expectedViewContent := "Hello world!"
	expectedViewContent1 := "Hello World!"
	expectedViewContent2 := "HELLO WORLD!"
	expectedViewContent3 := "Hello lord!!"
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

	if err := g.SetKeybinding("", KeyF1, ModNone, func(g *Gui, v *View) error {
		v.Clear()
		fmt.Fprintln(v, expectedViewContent1)
		return nil
	}); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", KeyF2, ModNone, func(g *Gui, v *View) error {
		v.Clear()
		<-time.After(time.Millisecond * 100)
		fmt.Fprintln(v, expectedViewContent2)
		return nil
	}); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", KeyF3, ModNone, func(g *Gui, v *View) error {
		v.Clear()
		fmt.Fprintln(v, expectedViewContent3)
		return nil
	}); err != nil {
		log.Panicln(err)
	}

	// Create a test screen and start gocui
	testingScreen := g.GetTestingScreen()
	cleanup := testingScreen.StartGui()
	defer cleanup()

	// check view content
	assertView(t, testingScreen, viewName, expectedViewContent)

	// Send a key to gocui
	testingScreen.SendKeySync(KeyCtrlC)

	// Test that the keybinding fired and set "didCallCTRLC" to true
	if !didCallCTRLC {
		t.Error("Expect the simulator to invoke the key handler for CTRLC")
	}

	// Send a key to gocui
	testingScreen.SendKeySync(KeyF1)

	// check view content
	assertView(t, testingScreen, viewName, expectedViewContent1)

	// Send a key to gocui
	testingScreen.SendKeySync(KeyF2)

	// check view content
	assertView(t, testingScreen, viewName, expectedViewContent2)

	// Send a key to gocui
	testingScreen.SendKeySync(KeyF3)

	// check view content
	assertView(t, testingScreen, viewName, expectedViewContent3)
}

func TestTestingScreenParallelKeys(t *testing.T) {
	// Track what happened in the view, we'll assert on these
	didCallCTRLC := false
	didCallF1 := false
	didCallF2 := false
	didCallF3 := false
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

	// Create a key bindings
	if err := g.SetKeybinding("", KeyCtrlC, ModNone, func(g *Gui, v *View) error {
		didCallCTRLC = true
		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", KeyF1, ModNone, func(g *Gui, v *View) error {
		didCallF1 = true
		return nil
	}); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", KeyF2, ModNone, func(g *Gui, v *View) error {
		didCallF2 = true
		return nil
	}); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", KeyF3, ModNone, func(g *Gui, v *View) error {
		didCallF3 = true
		return nil
	}); err != nil {
		log.Panicln(err)
	}

	// Create a test screen and start gocui
	testingScreen := g.GetTestingScreen()
	cleanup := testingScreen.StartGui()
	defer cleanup()

	// check view content
	assertView(t, testingScreen, viewName, expectedViewContent)

	// Send a key to gocui
	testingScreen.SendKeySync(KeyCtrlC)
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		testingScreen.SendKeySync(KeyF1)
		wg.Done()
	}()
	go func() {
		testingScreen.SendKeySync(KeyF2)
		wg.Done()
	}()
	go func() {
		testingScreen.SendKeySync(KeyF3)
		wg.Done()
	}()

	wg.Wait()

	// Test that the keybinding fired
	if !didCallCTRLC {
		t.Error("Expect the simulator to invoke the key handler for CTRLC")
	}
	if !didCallF1 || !didCallF2 || !didCallF3 {
		t.Error("Expect the simulator to invoke the key handler for F1, F2 and F3")
	}
}

// assertView checks if view contains provided content.
func assertView(t *testing.T, ts TestingScreen, viewName, content string) {
	t.Helper()
	// Get the content from the testing screen
	if actualContent, err := ts.GetViewContent(viewName); err != nil {
		t.Error(err)
	} else {
		// Test that it contains the "Hello World!" we thought the view should draw
		if strings.TrimSpace(actualContent) != content {
			t.Error(fmt.Printf("Expected view content to be: %q got: %q", content, actualContent))
		}
	}
}
