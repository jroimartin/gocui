// Copyright 2020 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

var simulationScreen tcell.SimulationScreen

type TestingScreen struct {
	screen tcell.SimulationScreen
	gui    *Gui
}

func isSimScreen() {
	if simulationScreen == nil {
		panic("Cannot use testing methods with a real screen use ")
	}
}

func (g *Gui) GetTestingScreen() TestingScreen {
	isSimScreen()

	return TestingScreen{
		screen: simulationScreen,
		gui:    g,
	}
}

func (t *TestingScreen) SendString(input string) {
	t.injectString(input)
}

func (s *TestingScreen) GetViewContent(viewName string) (string, error) {
	view, err := s.gui.View(viewName)
	if err != nil {
		return "", fmt.Errorf("Failed to retreive view: %w", err)
	}

	// Todo: Should we return the buffer here or the content of the sim screen?
	return view.Buffer(), nil
}

// Used from Micro MIT Licensed: https://github.com/zyedidia/micro/blob/c0907bb58e35ee05202a78226a2f53909af228ca/cmd/micro/micro_test.go#L183
func (t *TestingScreen) injectString(str string) {
	// the tcell simulation screen event channel can only handle
	// 10 events at once, so we need to divide up the key events
	// into chunks of 10 and handle the 10 events before sending
	// another chunk of events
	iters := len(str) / 10
	extra := len(str) % 10

	for i := 0; i < iters; i++ {
		s := i * 10
		e := i*10 + 10
		simulationScreen.InjectKeyBytes([]byte(str[s:e]))
		for i := 0; i < 10; i++ {
			// Trigger GoCUI to update with new input
			// Todo: Is this necessary?
			t.gui.Update(func(*Gui) error {
				return nil
			})
		}
	}

	simulationScreen.InjectKeyBytes([]byte(str[len(str)-extra:]))
	for i := 0; i < extra; i++ {
		// Trigger GoCUI to update with new input
		// Todo: Is this necessary?
		t.gui.Update(func(*Gui) error {
			return nil
		})
	}
}
