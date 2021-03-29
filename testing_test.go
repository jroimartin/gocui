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

	. "github.com/onsi/gomega"
)

func TestTestingScreenReturnsCorrectContent(t *testing.T) {
	// Setup gomega for async "eventually" assertions
	// see: http://onsi.github.io/gomega/#making-asynchronous-assertions
	assert := NewGomegaWithT(t)

	// Track what happened in the view, we'll assert on these
	didCallCTRLC := false
	expectedViewContent := "Hello world!"
	viewName := "testView1"

	// Create a view
	g := setupViews(viewName, expectedViewContent)

	// Create a key binding which sets "didCallCTRLC" when triggered
	if err := g.SetKeybinding("", KeyCtrlC, ModNone, func(g *Gui, v *View) error { didCallCTRLC = true; return nil }); err != nil {
		log.Panicln(err)
	}

	// Create a test screen and start gocui
	testingScreen := g.GetTestingScreen()
	cleanup := testingScreen.StartTestingScreen()
	defer cleanup()

	// Send a key to gocui
	testingScreen.SendKey(KeyCtrlC)

	// Use gomega asserts "eventually" to handle the async drawing
	// of the view and handling of the events
	//
	// Check the key binding was called
	assert.
		Eventually(func() bool { return didCallCTRLC }).
		Should(Equal(true), "Expect the simulator to invoke the key handler for CTRLC")

	// Check the content was drawn onto the screen
	assert.
		Eventually(func() string {
			// Get the content of the "hello" view
			actualContent, err := testingScreen.GetViewContent(viewName)
			if err != nil {
				t.Error(err)
			}
			return strings.TrimSpace(actualContent)
		}).Should(Equal(expectedViewContent))
}

func setupViews(viewName, expectedViewContent string) *Gui {
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

			fmt.Fprintln(v, expectedViewContent)
		}

		return nil
	})

	return g
}
