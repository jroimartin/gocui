// Copyright 2020 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This is copy from original github.com/gdamore/tcell/termbox package
// for easier adoption of tcell.
// There are some changes made, to make it work with termbox keys (Ctrl modifier)

package gocui

import (
	"github.com/gdamore/tcell/v2"
)

var screen tcell.Screen
var outMode OutputMode

// tcellInit initializes the screen for use.
func tcellInit() error {
	outMode = OutputNormal
	if s, e := tcell.NewScreen(); e != nil {
		return e
	} else if e = s.Init(); e != nil {
		return e
	} else {
		screen = s
		return nil
	}
}

// tcellClose cleans up the terminal, restoring terminal modes, etc.
func tcellClose() {
	screen.Fini()
}

// tcellFlush updates the screen.
func tcellFlush() error {
	screen.Show()
	return nil
}

// tcellSetCursor displays the terminal cursor at the given location.
func tcellSetCursor(x, y int) {
	screen.ShowCursor(x, y)
}

// tcellHideCursor hides the terminal cursor.
func tcellHideCursor() {
	tcellSetCursor(-1, -1)
}

// tcellSize returns the screen size as width, height in character cells.
func tcellSize() (int, int) {
	return screen.Size()
}

// tcellClear clears the screen with the given attributes.
func tcellClear(fg, bg Attribute) {
	st := mkStyle(fg, bg)
	w, h := screen.Size()
	for row := 0; row < h; row++ {
		for col := 0; col < w; col++ {
			screen.SetContent(col, row, ' ', nil, st)
		}
	}
}

// InputMode is not used.
type InputMode int

// Unused input modes; here for compatibility.
const (
	InputEsc InputMode = 1 << iota
	InputAlt
	InputMouse
	InputCurrent InputMode = 0
)

// tcellSetInputMode does not do anything in this version.
func tcellSetInputMode(mode InputMode) InputMode {
	if mode&InputMouse != 0 {
		screen.EnableMouse()
		return InputEsc | InputMouse
	}
	// We don't do anything else right now
	return InputEsc
}

// OutputMode represents an output mode, which determines how colors
// are used.  See the termbox documentation for an explanation.
type OutputMode int

// OutputMode values.
const (
	OutputCurrent OutputMode = iota
	OutputNormal
	Output256
	Output216
	OutputGrayscale
	OutputTrue
)

// tcellSetOutputMode is used to set the color palette used.
func tcellSetOutputMode(mode OutputMode) OutputMode {
	if screen.Colors() < 256 {
		mode = OutputNormal
	}
	switch mode {
	case OutputCurrent:
		return outMode
	case OutputNormal, Output256, Output216, OutputGrayscale, OutputTrue:
		outMode = mode
		return mode
	default:
		return outMode
	}
}

// tcellSync forces a resync of the screen.
func tcellSync() error {
	screen.Sync()
	return nil
}

// tcellSetCell sets the character cell at a given location to the given
// content (rune) and attributes.
func tcellSetCell(x, y int, ch rune, fg, bg Attribute) {
	st := mkStyle(fg, bg)
	screen.SetContent(x, y, ch, nil, st)
}

// EventType represents the type of event.
type EventType uint8

// Event represents an event like a key press, mouse action, or window resize.
type Event struct {
	Type   EventType
	Mod    Modifier
	Key    Key
	Ch     rune
	Width  int
	Height int
	Err    error
	MouseX int
	MouseY int
	N      int
}

// Event types.
const (
	EventNone EventType = iota
	EventKey
	EventResize
	EventMouse
	EventInterrupt
	EventError
	EventRaw
)

var (
	lastMouseKey tcell.ButtonMask = tcell.ButtonNone
	lastMouseMod tcell.ModMask    = tcell.ModNone
)

func makeEvent(tev tcell.Event) Event {
	switch tev := tev.(type) {
	case *tcell.EventInterrupt:
		return Event{Type: EventInterrupt}
	case *tcell.EventResize:
		w, h := tev.Size()
		return Event{Type: EventResize, Width: w, Height: h}
	case *tcell.EventKey:
		k := tev.Key()
		ch := rune(0)
		if k == tcell.KeyRune {
			k = 0 // if rune remove key (so it can match, for example spacebar)
			ch = tev.Rune()
		}
		mod := tev.Modifiers()
		// remove control modifier and setup special handling of ctrl+spacebar, etc.
		if mod == tcell.ModCtrl && k == 0 && ch == ' ' {
			mod = 0
			ch = rune(0)
			k = tcell.KeyCtrlSpace
		} else if mod == tcell.ModCtrl || mod == tcell.ModShift {
			// remove Ctrl or Shift if specified
			// - shift - will be translated to the final code of rune
			// - ctrl  - is translated in the key
			mod = 0
		}
		return Event{
			Type: EventKey,
			Key:  Key(k),
			Ch:   ch,
			Mod:  Modifier(mod),
		}
	case *tcell.EventMouse:
		x, y := tev.Position()
		button := tev.Buttons()
		mouseKey := MouseRelease
		mouseMod := ModNone
		// process mouse wheel
		if button&tcell.WheelUp != 0 {
			mouseKey = MouseWheelUp
		}
		if button&tcell.WheelDown != 0 {
			mouseKey = MouseWheelDown
		}
		if button&tcell.WheelLeft != 0 {
			mouseKey = MouseWheelLeft
		}
		if button&tcell.WheelRight != 0 {
			mouseKey = MouseWheelRight
		}

		// process button events (not wheel events)
		button &= tcell.ButtonMask(0xff)
		if button != tcell.ButtonNone && lastMouseKey == tcell.ButtonNone {
			lastMouseKey = button
			lastMouseMod = tev.Modifiers()
		}

		switch tev.Buttons() {
		case tcell.ButtonNone:
			if lastMouseKey != tcell.ButtonNone {
				switch lastMouseKey {
				case tcell.ButtonPrimary:
					mouseKey = MouseLeft
				case tcell.ButtonSecondary:
					mouseKey = MouseRight
				case tcell.ButtonMiddle:
					mouseKey = MouseMiddle
				}
				mouseMod = Modifier(lastMouseMod)
				lastMouseMod = tcell.ModNone
				lastMouseKey = tcell.ButtonNone
			}
		}

		return Event{
			Type:   EventMouse,
			MouseX: x,
			MouseY: y,
			Key:    mouseKey,
			Ch:     0,
			Mod:    mouseMod,
		}
	default:
		return Event{Type: EventNone}
	}
}
