// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import "github.com/nsf/termbox-go"

// Keybidings are used to link a given key-press event with a handler.
type keybinding struct {
	viewName string
	key      Key
	ch       rune
	mod      Modifier
	handler  func(*Gui, *View) error
}

// newKeybinding returns a new Keybinding object.
func newKeybinding(viewname string, key Key, ch rune, mod Modifier, handler func(*Gui, *View) error) (kb *keybinding) {
	kb = &keybinding{
		viewName: viewname,
		key:      key,
		ch:       ch,
		mod:      mod,
		handler:  handler,
	}
	return kb
}

// matchKeypress returns if the keybinding matches the keypress.
func (kb *keybinding) matchKeypress(key Key, ch rune, mod Modifier) bool {
	return kb.key == key && kb.ch == ch && kb.mod == mod
}

// matchView returns if the keybinding matches the current view.
func (kb *keybinding) matchView(v *View) bool {
	if kb.viewName == "" {
		return true
	}
	return v != nil && kb.viewName == v.name
}

// Key represents special keys or keys combinations.
type Key termbox.Key

// Special keys.
const (
	KeyF1         Key = Key(termbox.KeyF1)
	KeyF2             = Key(termbox.KeyF2)
	KeyF3             = Key(termbox.KeyF3)
	KeyF4             = Key(termbox.KeyF4)
	KeyF5             = Key(termbox.KeyF5)
	KeyF6             = Key(termbox.KeyF6)
	KeyF7             = Key(termbox.KeyF7)
	KeyF8             = Key(termbox.KeyF8)
	KeyF9             = Key(termbox.KeyF9)
	KeyF10            = Key(termbox.KeyF10)
	KeyF11            = Key(termbox.KeyF11)
	KeyF12            = Key(termbox.KeyF12)
	KeyInsert         = Key(termbox.KeyInsert)
	KeyDelete         = Key(termbox.KeyDelete)
	KeyHome           = Key(termbox.KeyHome)
	KeyEnd            = Key(termbox.KeyEnd)
	KeyPgup           = Key(termbox.KeyPgup)
	KeyPgdn           = Key(termbox.KeyPgdn)
	KeyArrowUp        = Key(termbox.KeyArrowUp)
	KeyArrowDown      = Key(termbox.KeyArrowDown)
	KeyArrowLeft      = Key(termbox.KeyArrowLeft)
	KeyArrowRight     = Key(termbox.KeyArrowRight)

	MouseLeft      = Key(termbox.MouseLeft)
	MouseMiddle    = Key(termbox.MouseMiddle)
	MouseRight     = Key(termbox.MouseRight)
	MouseRelease   = Key(termbox.MouseRelease)
	MouseWheelUp   = Key(termbox.MouseWheelUp)
	MouseWheelDown = Key(termbox.MouseWheelDown)
)

// Keys combinations.
const (
	KeyCtrlTilde      Key = Key(termbox.KeyCtrlTilde)
	KeyCtrl2              = Key(termbox.KeyCtrl2)
	KeyCtrlSpace          = Key(termbox.KeyCtrlSpace)
	KeyCtrlA              = Key(termbox.KeyCtrlA)
	KeyCtrlB              = Key(termbox.KeyCtrlB)
	KeyCtrlC              = Key(termbox.KeyCtrlC)
	KeyCtrlD              = Key(termbox.KeyCtrlD)
	KeyCtrlE              = Key(termbox.KeyCtrlE)
	KeyCtrlF              = Key(termbox.KeyCtrlF)
	KeyCtrlG              = Key(termbox.KeyCtrlG)
	KeyBackspace          = Key(termbox.KeyBackspace)
	KeyCtrlH              = Key(termbox.KeyCtrlH)
	KeyTab                = Key(termbox.KeyTab)
	KeyCtrlI              = Key(termbox.KeyCtrlI)
	KeyCtrlJ              = Key(termbox.KeyCtrlJ)
	KeyCtrlK              = Key(termbox.KeyCtrlK)
	KeyCtrlL              = Key(termbox.KeyCtrlL)
	KeyEnter              = Key(termbox.KeyEnter)
	KeyCtrlM              = Key(termbox.KeyCtrlM)
	KeyCtrlN              = Key(termbox.KeyCtrlN)
	KeyCtrlO              = Key(termbox.KeyCtrlO)
	KeyCtrlP              = Key(termbox.KeyCtrlP)
	KeyCtrlQ              = Key(termbox.KeyCtrlQ)
	KeyCtrlR              = Key(termbox.KeyCtrlR)
	KeyCtrlS              = Key(termbox.KeyCtrlS)
	KeyCtrlT              = Key(termbox.KeyCtrlT)
	KeyCtrlU              = Key(termbox.KeyCtrlU)
	KeyCtrlV              = Key(termbox.KeyCtrlV)
	KeyCtrlW              = Key(termbox.KeyCtrlW)
	KeyCtrlX              = Key(termbox.KeyCtrlX)
	KeyCtrlY              = Key(termbox.KeyCtrlY)
	KeyCtrlZ              = Key(termbox.KeyCtrlZ)
	KeyEsc                = Key(termbox.KeyEsc)
	KeyCtrlLsqBracket     = Key(termbox.KeyCtrlLsqBracket)
	KeyCtrl3              = Key(termbox.KeyCtrl3)
	KeyCtrl4              = Key(termbox.KeyCtrl4)
	KeyCtrlBackslash      = Key(termbox.KeyCtrlBackslash)
	KeyCtrl5              = Key(termbox.KeyCtrl5)
	KeyCtrlRsqBracket     = Key(termbox.KeyCtrlRsqBracket)
	KeyCtrl6              = Key(termbox.KeyCtrl6)
	KeyCtrl7              = Key(termbox.KeyCtrl7)
	KeyCtrlSlash          = Key(termbox.KeyCtrlSlash)
	KeyCtrlUnderscore     = Key(termbox.KeyCtrlUnderscore)
	KeySpace              = Key(termbox.KeySpace)
	KeyBackspace2         = Key(termbox.KeyBackspace2)
	KeyCtrl8              = Key(termbox.KeyCtrl8)
)

// Modifier allows to define special keys combinations. They can be used
// in combination with Keys or Runes when a new keybinding is defined.
type Modifier termbox.Modifier

// Modifiers.
const (
	ModNone Modifier = Modifier(0)
	ModAlt           = Modifier(termbox.ModAlt)
)

// DescribeKey generates a human-readable description of a key combo.
func DescribeKey(key interface{}, mod Modifier) string {

	var k string

	switch mod {
	case ModNone:
	case ModAlt:
		k = "Alt+"
	default:
		k = "<unknown modifier>+"
	}
	
	switch t := key.(type) {
	case Key:
		s,ok := keysymbols[t]
		if !ok {
			k += "<unknown key>"
		} else {
			k += s
		}
	case rune:
		k += string(t)
	default:
		k += "<unknown key type>"
	}

	return k

}

var keysymbols = map[Key]string{
	// Entries commented out duplicate codes with a more prominent combo.
	KeyF1: "F1",
	KeyF2: "F2",
	KeyF3: "F3",
	KeyF4: "F4",
	KeyF5: "F5",
	KeyF6: "F6",
	KeyF7: "F7",
	KeyF8: "F8",
	KeyF9: "F9",
	KeyF10: "F10",
	KeyF11: "F11",
	KeyF12: "F12",
	KeyInsert: "Insert",
	KeyDelete: "Delete",
	KeyHome: "Home",
	KeyEnd: "End",
	KeyPgup: "PgUp",
	KeyPgdn: "PgDn",
	KeyArrowUp: "Up",
	KeyArrowDown: "Down",
	KeyArrowLeft: "Left",
	KeyArrowRight: "Right",
//	KeyCtrlTilde: "Ctrl+~",
//	KeyCtrl2: "Ctrl+2",
	KeyCtrlSpace: "Ctrl+Space",
	KeyCtrlA: "Ctrl+a",
	KeyCtrlB: "Ctrl+b",
	KeyCtrlC: "Ctrl+c",
	KeyCtrlD: "Ctrl+d",
	KeyCtrlE: "Ctrl+e",
	KeyCtrlF: "Ctrl+f",
	KeyCtrlG: "Ctrl+g",
	KeyBackspace: "Backspace",
//	KeyCtrlH: "Ctrl+h",
	KeyTab: "Ctrl+Tab",
//	KeyCtrlI: "Ctrl+i",
	KeyCtrlJ: "Ctrl+j",
	KeyCtrlK: "Ctrl+k",
	KeyCtrlL: "Ctrl+l",
	KeyEnter: "Ctrl+Enter",
//	KeyCtrlM: "Ctrl+m",
	KeyCtrlN: "Ctrl+n",
	KeyCtrlO: "Ctrl+o",
	KeyCtrlP: "Ctrl+p",
	KeyCtrlQ: "Ctrl+q",
	KeyCtrlR: "Ctrl+r",
	KeyCtrlS: "Ctrl+s",
	KeyCtrlT: "Ctrl+t",
	KeyCtrlU: "Ctrl+u",
	KeyCtrlV: "Ctrl+v",
	KeyCtrlW: "Ctrl+w",
	KeyCtrlX: "Ctrl+x",
	KeyCtrlY: "Ctrl+y",
	KeyCtrlZ: "Ctrl+z",
	KeyEsc: "Esc",
//	KeyCtrlLsqBracket: "Ctrl+[",
//	KeyCtrl3: "Ctrl+3",
//	KeyCtrl4: "Ctrl+4",
	KeyCtrlBackslash: "Ctrl+\\",
//	KeyCtrl5: "Ctrl+5",
	KeyCtrlRsqBracket: "Ctrl+]",
	KeyCtrl6: "Ctrl+6",
//	KeyCtrl7: "Ctrl+7",
	KeyCtrlSlash: "Ctrl+/",
//	KeyCtrlUnderscore: "Ctrl+_",
	KeySpace: "Ctrl+Space",
	KeyBackspace2: "Ctrl+Backspace",
//	KeyCtrl8: "Ctrl+8",
}
