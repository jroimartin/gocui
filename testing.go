// Copyright 2021 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"errors"
	"fmt"
	"log"
	"strings"

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

func (t *TestingScreen) StartTestingScreen() func() {
	go func() {
		if err := t.gui.MainLoop(); (err != nil && !errors.Is(err, ErrQuit)) {
			log.Panic(err)
		}
	}()

	// Return a func that will stop the main loop
	return func() { 
		t.gui.stop <- struct{}{}
	}
}

func (t *TestingScreen) SendStringAsKeys(input string) {
	t.injectString(input)
}

func (t *TestingScreen) SendKey(key Key) {
	t.screen.InjectKey(tcell.Key(key), rune(key), tcell.ModNone)
}

func (s *TestingScreen) GetViewContent(viewName string) (string, error) {
	view, err := s.gui.View(viewName)
	if err != nil {
		return "", fmt.Errorf("Failed to retreive view: %w", err)
	}

	x0, y0, x1, y1 := view.Dimensions()

	// Account for the border
	x0++
	y0++
	x1--
	y1--

	// Walk each line in the view
	var result strings.Builder
	Xcurrent := x0
	Ycurrent := y0
	for y0 < y1 || y0 == y1 {
		// Did we reach the end of the line?
		if Xcurrent > x1 {
			Xcurrent = x0
			Ycurrent++
			result.WriteString("\n")
		}

		// Did we reach the bottom of the view?
		if Ycurrent > y1 {
			break
		}

		// Get the content (without formatting) at that position
		content, err := s.gui.Rune(Xcurrent, Ycurrent)
		if err != nil {
			return "", fmt.Errorf("Failed reading rune from simulation screen: %w", err)
		}
		result.WriteRune(content)
		Xcurrent++
	}

	return result.String(), nil
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
