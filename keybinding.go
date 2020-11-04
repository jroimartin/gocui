// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"strings"
)

// Keybidings are used to link a given key-press event with a handler.
type keybinding struct {
	viewName string
	key      Key
	ch       rune
	mod      Modifier
	handler  func(*Gui, *View) error
}

// Parse takes the input string and extracts the keybinding.
// Returns a Key / rune, a Modifier and an error.
func Parse(input string) (interface{}, Modifier, error) {
	if len(input) == 1 {
		_, r, err := getKey(rune(input[0]))
		if err != nil {
			return nil, ModNone, err
		}
		return r, ModNone, nil
	}

	var modifier Modifier
	cleaned := make([]string, 0)

	tokens := strings.Split(input, "+")
	for _, t := range tokens {
		normalized := strings.Title(strings.ToLower(t))
		if t == "Alt" {
			modifier = ModAlt
			continue
		}
		cleaned = append(cleaned, normalized)
	}

	key, exist := translate[strings.Join(cleaned, "")]
	if !exist {
		return nil, ModNone, ErrNoSuchKeybind
	}

	return key, modifier, nil
}

// ParseAll takes an array of strings and returns a map of all keybindings.
func ParseAll(input []string) (map[interface{}]Modifier, error) {
	ret := make(map[interface{}]Modifier)
	for _, i := range input {
		k, m, err := Parse(i)
		if err != nil {
			return ret, err
		}
		ret[k] = m
	}
	return ret, nil
}

// MustParse takes the input string and returns a Key / rune and a Modifier.
// It will panic if any error occured.
func MustParse(input string) (interface{}, Modifier) {
	k, m, err := Parse(input)
	if err != nil {
		panic(err)
	}
	return k, m
}

// MustParseAll takes an array of strings and returns a map of all keybindings.
// It will panic if any error occured.
func MustParseAll(input []string) map[interface{}]Modifier {
	result, err := ParseAll(input)
	if err != nil {
		panic(err)
	}
	return result
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
	// if the user is typing in a field, ignore char keys
	if v == nil || (v.Editable && kb.ch != 0 && !v.KeybindOnEdit) {
		return false
	}
	return kb.viewName == v.name
}

// translations for strings to keys
var translate = map[string]Key{
	"F1":         KeyF1,
	"F2":         KeyF2,
	"F3":         KeyF3,
	"F4":         KeyF4,
	"F5":         KeyF5,
	"F6":         KeyF6,
	"F7":         KeyF7,
	"F8":         KeyF8,
	"F9":         KeyF9,
	"F10":        KeyF10,
	"F11":        KeyF11,
	"F12":        KeyF12,
	"Insert":     KeyInsert,
	"Delete":     KeyDelete,
	"Home":       KeyHome,
	"End":        KeyEnd,
	"Pgup":       KeyPgup,
	"Pgdn":       KeyPgdn,
	"ArrowUp":    KeyArrowUp,
	"ArrowDown":  KeyArrowDown,
	"ArrowLeft":  KeyArrowLeft,
	"ArrowRight": KeyArrowRight,
	// "CtrlTilde":      KeyCtrlTilde,
	"Ctrl2": KeyCtrl2,
	// "CtrlSpace":      KeyCtrlSpace,
	"CtrlA":          KeyCtrlA,
	"CtrlB":          KeyCtrlB,
	"CtrlC":          KeyCtrlC,
	"CtrlD":          KeyCtrlD,
	"CtrlE":          KeyCtrlE,
	"CtrlF":          KeyCtrlF,
	"CtrlG":          KeyCtrlG,
	"Backspace":      KeyBackspace,
	"CtrlH":          KeyCtrlH,
	"Tab":            KeyTab,
	"CtrlI":          KeyCtrlI,
	"CtrlJ":          KeyCtrlJ,
	"CtrlK":          KeyCtrlK,
	"CtrlL":          KeyCtrlL,
	"Enter":          KeyEnter,
	"CtrlM":          KeyCtrlM,
	"CtrlN":          KeyCtrlN,
	"CtrlO":          KeyCtrlO,
	"CtrlP":          KeyCtrlP,
	"CtrlQ":          KeyCtrlQ,
	"CtrlR":          KeyCtrlR,
	"CtrlS":          KeyCtrlS,
	"CtrlT":          KeyCtrlT,
	"CtrlU":          KeyCtrlU,
	"CtrlV":          KeyCtrlV,
	"CtrlW":          KeyCtrlW,
	"CtrlX":          KeyCtrlX,
	"CtrlY":          KeyCtrlY,
	"CtrlZ":          KeyCtrlZ,
	"Esc":            KeyEsc,
	"CtrlLsqBracket": KeyCtrlLsqBracket,
	"Ctrl3":          KeyCtrl3,
	"Ctrl4":          KeyCtrl4,
	"CtrlBackslash":  KeyCtrlBackslash,
	"Ctrl5":          KeyCtrl5,
	"CtrlRsqBracket": KeyCtrlRsqBracket,
	"Ctrl6":          KeyCtrl6,
	"Ctrl7":          KeyCtrl7,
	"CtrlSlash":      KeyCtrlSlash,
	"CtrlUnderscore": KeyCtrlUnderscore,
	"Space":          KeySpace,
	"Backspace2":     KeyBackspace2,
	// "Ctrl8":          KeyCtrl8,
	"Mouseleft":      MouseLeft,
	"Mousemiddle":    MouseMiddle,
	"Mouseright":     MouseRight,
	"Mouserelease":   MouseRelease,
	"MousewheelUp":   MouseWheelUp,
	"MousewheelDown": MouseWheelDown,
}
